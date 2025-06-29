package router

import (
    "todo/handler"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h handler.Handle) {
    app.Get("/health", h.HealthCheck)
    api := app.Group("/api")
    todos := api.Group("/todos")
    todos.Get("/", h.ListTodos)
    todos.Post("/", h.CreateTodo)
    todos.Get("/:id", h.ReadTodo)
    todos.Put("/:id", h.UpdateTodo)
    todos.Delete("/:id", h.DeleteTodo)
    todos20 := api.Group("/todos20")
    todos20.Get("/", h.ListTodos20)
}
