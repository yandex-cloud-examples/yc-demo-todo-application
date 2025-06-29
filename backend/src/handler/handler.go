package handler
import (
    "errors"
    "todo/model"
    "gorm.io/gorm"
    "gorm.io/plugin/dbresolver"
    "github.com/gofiber/fiber/v2"
)

type Handle struct {
    DB *gorm.DB
}

func ProcessDbError(c *fiber.Ctx, err error) error {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return c.Status(404).JSON(fiber.Map{"message": "Not found"})
    }
    return c.Status(500).JSON(fiber.Map{"message": "Unexpected error"})
}

func (h Handle) ListTodos(c *fiber.Ctx) error {
    var todos []model.Todo
    db := h.DB.Clauses(dbresolver.Read)
    if r := db.Find(&todos); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    return c.Status(200).JSON(&todos)
}

func (h Handle) ListTodos20(c *fiber.Ctx) error {
    var todos []model.Todo
    // Choose the connection for reading
    tx := h.DB.Clauses(dbresolver.Read).Session(&gorm.Session{})
    // Select 20 ids
    if r := tx.Limit(20).Select("id").Find(&todos); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    // Loop through the list of ids and select a record for each id
    for idx, _ := range todos {
        if r := tx.First(&todos[idx]); r.Error != nil {
            return ProcessDbError(c, r.Error)
        }
    }
    return c.Status(200).JSON(&todos)
}

func (h Handle) CreateTodo(c *fiber.Ctx) error {
    var todo model.Todo

    if err := c.BodyParser(&todo); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Bad request"})
    }
    if r := h.DB.Create(&todo); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    return c.Status(201).JSON(&todo)
}

func (h Handle) ReadTodo(c *fiber.Ctx) error {
    var todo model.Todo
    if r := h.DB.Clauses(dbresolver.Read).First(&todo, c.Params("id")); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    return c.Status(200).JSON(&todo)
}

func (h Handle) UpdateTodo(c *fiber.Ctx) error {
    var todo model.Todo
    var updateTodo model.Todo

    if err := c.BodyParser(&updateTodo); err != nil {
        return c.Status(400).JSON(fiber.Map{"message": "Bad request"})
    }
    tx := h.DB.Clauses(dbresolver.Read).Session(&gorm.Session{})
    if r := tx.First(&todo, c.Params("id")); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    todo.Title = updateTodo.Title
    todo.Description = updateTodo.Description
    todo.Completed = updateTodo.Completed
    if r := tx.Save(&todo); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    return c.Status(200).JSON(&todo)
}

func (h Handle) DeleteTodo(c *fiber.Ctx) error {
    var todo model.Todo
    if r := h.DB.Delete(&todo, "id = ?", c.Params("id")); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    return c.SendStatus(204)
}

func (h Handle) HealthCheck(c *fiber.Ctx) error {
    var result bool

    if r := h.DB.Clauses(dbresolver.Write).Raw("SELECT pg_is_in_recovery()").Scan(&result); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    
    if result {
        return c.Status(500).JSON(fiber.Map{"message": "Master is away"})
    } 

    if r := h.DB.Clauses(dbresolver.Read).Raw("SELECT true").Scan(&result); r.Error != nil {
        return ProcessDbError(c, r.Error)
    }
    
    return c.Status(200).JSON(fiber.Map{"message": "OK"})
}


