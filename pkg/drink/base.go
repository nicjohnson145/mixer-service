package drink

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
)

const (
	DrinkPublicityPublic  = "public"
	DrinkPublicityPrivate = "private"
)

type Drink struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required"`
	DrinkData
}

var validate = validator.New()

func Init(app *fiber.App, db *sql.DB) error {
	defineRoutes(app, db)
	return nil
}

func defineRoutes(app *fiber.App, db *sql.DB) {
	app.Post(common.DrinksV1+"/create", auth.RequiresValidAccessToken(createDrink(db)))
	app.Get(common.DrinksV1+"/:id", auth.RequiresValidAccessToken(getDrink(db)))
	app.Delete(common.DrinksV1+"/:id", auth.RequiresValidAccessToken(deleteDrink(db)))
	app.Put(common.DrinksV1+"/:id", auth.RequiresValidAccessToken(updateDrink(db)))
	app.Get(common.DrinksV1+"/by-user/:username", auth.RequiresValidAccessToken(getDrinksByUser(db)))
	app.Post(common.DrinksV1+"/:id/copy", auth.RequiresValidAccessToken(copyDrink(db)))
	app.Static("/static", "../../pkg/public/webui")
}
