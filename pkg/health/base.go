package health

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
)

func Init(app *fiber.App, db *sql.DB) error {
	defineRoutes(app, db)
	return nil
}

func defineRoutes(app *fiber.App, db *sql.DB) {
	app.Get(common.HealthV1, healthCheck)
}

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("ok")
}
