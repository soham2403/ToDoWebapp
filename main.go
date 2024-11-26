package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type ToDo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

// DataBase function
func createTable(db *sql.DB, tableName string) {
	if tableName == "" {
		log.Fatal("Table name can't be empty.")
	}
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INT AUTO_INCREMENT PRIMARY KEY,
		completed BOOL NOT NULL,
		body VARCHAR(50) NOT NULL
	);`, tableName)

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error while creating table", err)
	}
	fmt.Println("Table", tableName, "created successfully!!")
}

func getTodos(c *fiber.Ctx, db *sql.DB) error {
	query := `SELECT id, completed, body FROM todos`
	rows, err := db.Query(query)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"err": "Error while fetching todos"})
	}
	defer rows.Close()
	var ToDos []ToDo
	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.Id, &todo.Completed, &todo.Body); err != nil {
			return c.Status(500).JSON(fiber.Map{"err": "error scanning data"})
		}
		ToDos = append(ToDos, todo)
	}
	if len(ToDos) == 0 {
		return c.Status(200).JSON([]ToDo{})
	}
	return c.Status(200).JSON(ToDos)
}

func createTodo(c *fiber.Ctx, db *sql.DB) error {
	todo := &ToDo{}

	if err := c.BodyParser(todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"err": "body field is required."})
	}

	query := `INSERT INTO todos (completed, body) VALUES (?,?)`
	result, err := db.Exec(query, todo.Completed, todo.Body)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create todo"})
	}

	id, _ := result.LastInsertId()
	todo.Id = int(id)
	return c.Status(201).JSON(todo)
}

func toggleTodo(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")

	// Check if the ID exists
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM todos WHERE id = ?)`
	err := db.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil || !exists {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	// If exists, toggle the 'completed' status
	query := `UPDATE todos SET completed = NOT completed WHERE id = ?`
	_, err = db.Exec(query, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to toggle todo"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo toggled successfully"})
}

func deleteTodo(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	query := `DELETE FROM todos WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Todo deleted successfully"})
}

func main() {
	app := fiber.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while loading .env file")
	}

	PORT := os.Getenv("PORT")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf(`%s:%s@tcp(127.0.0.1:3306)/ToDo`, DB_USERNAME, DB_PASSWORD)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed", err)
	}
	fmt.Println("Connected to Mysql database")
	createTable(db, "todos")
	app.Get("/api/todos", func(c *fiber.Ctx) error { return getTodos(c, db) })
	app.Post("/api/todos", func(c *fiber.Ctx) error { return createTodo(c, db) })
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error { return toggleTodo(c, db) })
	app.Delete("/api/todos/delete/:id", func(c *fiber.Ctx) error { return deleteTodo(c, db) })

	log.Fatal(app.Listen(":" + PORT))
}
