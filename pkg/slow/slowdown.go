package slow

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func SlowDown(app *fiber.App) {
	val, ok := os.LookupEnv("SLOWDOWN_AMOUNT")
	if !ok {
		val = "2s"
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Unable to parse SLOWDOWN_AMOUNT as duration: %v", err)
	}

	log.Infof("Slowdown active, all requests will have a %v delay added to them", duration)

	app.Use(func(c *fiber.Ctx) error {
		log.Debug("Sleeping due to slowdown mode")
		time.Sleep(duration)
		return c.Next()
	})
}
