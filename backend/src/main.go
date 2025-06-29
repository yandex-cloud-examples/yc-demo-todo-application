package main

import (
    "os"
    "log"
    "fmt"
    "todo/migration"
    "todo/database"
    "todo/router"
    "todo/config"
    "todo/handler"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    _ "github.com/lib/pq"
)

func main() {
    config := config.NewDbConfig()
    db := database.ConnectRW(&config.RW)
    args := os.Args[1:]
    action := "run"
    for i := 0; i < len(args); i++ {
        action = args[i]
        switch action {
            case "migrate":
                migration.Migrate(db)
            case "load":
                i++
                if i >= len(args) {
                    log.Fatal("Filename required after load\n")
                }
                migration.LoadData(db, args[i])
            case "run":
                break
            default:
                log.Fatal(fmt.Sprintf("Unknown command `%s'\n", action))
        }
    }

    if action != "run" {
        return
    }

    if config.RW.Host != config.RO.Host {
        database.ConnectRO(&config.RO, db)
    }

    h := handler.Handle{DB: db}
    app := fiber.New()
    app.Use(logger.New())
    app.Use(cors.New())

    router.SetupRoutes(app, h)

    app.Use(func(c *fiber.Ctx) error {
        return c.SendStatus(404)
    })
    app.Listen(":8080")
}
