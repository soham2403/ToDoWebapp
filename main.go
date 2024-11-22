package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type ToDo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var ToDos []ToDo

func validateTodoId(id string) (int, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	if idInt <= 0 || idInt > len(ToDos) {
		return 0, fmt.Errorf("%s out of bound", id)
	}
	return idInt - 1, nil
}

func getTodos(c *fiber.Ctx) error {
	return c.Status(200).JSON(ToDos)
}

func createTodo(c *fiber.Ctx) error {
	todo := &ToDo{}

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"err": "body field is required."})
	}

	todo.Id = len(ToDos) + 1
	ToDos = append(ToDos, *todo)
	return c.Status(201).JSON(todo)
}

func toggleTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	index, err := validateTodoId(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"err": err.Error()})
	}

	ToDos[index].Completed = !ToDos[index].Completed
	return c.Status(201).JSON(ToDos[index])
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	index, err := validateTodoId(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"err": err.Error()})
	}

	deletedTodo := ToDos[index]
	//use to change but didnt change index here.
	ToDos = append(ToDos[:index], ToDos[index+1:]...)

	//used to change the index.
	for i := range ToDos {
		ToDos[i].Id = i + 1
	}

	return c.Status(200).JSON(deletedTodo)
}

func main() {
	fmt.Println("hello world")
	app := fiber.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while loading .env file")
	}

	PORT := os.Getenv("PORT")

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", toggleTodo)
	app.Delete("/api/todos/delete/:id", deleteTodo)

	log.Fatal(app.Listen(":" + PORT))
}
