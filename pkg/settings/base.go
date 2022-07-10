package settings

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
)

type UserSettings struct {
	PublicProfile bool `json:"public_profile"`
}

const (
	PublicProfile = "public_profile"
	defaultPublicProfile = true
)

func Init(app *fiber.App, db *sql.DB) error {
	defineRoutes(app, db)
	return nil
}

func defineRoutes(app *fiber.App, db *sql.DB) {
	app.Get(common.SettingsV1, auth.RequiresValidAccessToken(getUserSettings(db)))
	app.Put(common.SettingsV1, auth.RequiresValidAccessToken(writeUserSettings(db)))
}
