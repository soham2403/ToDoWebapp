package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ToDo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("hello world")
	app := fiber.New()
	ToDos := []ToDo{}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello world"})
	})
	// create a todo endpoint
	app.Post("/api/todos", func(c *fiber.Ctx) error {
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
	})
	// list of all tasks
	//toggle method
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"err": "invalid id"})
		}

		if idInt <= 0 || idInt > len(ToDos) {
			return c.Status(400).JSON(fiber.Map{"err": id + " out of bound"})
		}
		idInt--
		completedTask := ToDos[idInt]
		fmt.Println(completedTask)
		completedTask.Completed = !completedTask.Completed
		return c.Status(201).JSON(fiber.Map{"msg": id + " toggled"})
	})
	log.Fatal(app.Listen(":4000"))
}
