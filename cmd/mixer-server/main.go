package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/nicjohnson145/mixer-service/pkg/health"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Use(logger.New())

	db := db.NewDBOrDie(common.DefaultedEnvVar("DB_PATH", "mixer.db"))

	if err := auth.Init(app, db); err != nil {
		log.Fatal(err)
	}

	//if err := drink.Init(r, db); err != nil {
	//    log.Fatal(err)
	//}

	if err := health.Init(app, db); err != nil {
		log.Fatal(err)
	}

	port := common.DefaultedEnvVar("PORT", "30000")

	log.Info("Listening on port ", port)
	app.Listen(":" + port)
}
