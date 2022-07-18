package drink

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type CopyDrinkResponse struct {
	ID int64 `json:"id,omitempty"`
}

func copyDrink(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		drink, err := getDrinkInternal(db, c, claims)
		if err != nil {
			return err
		}

		if c.Query("newName") != "" {
			drink.DrinkData.Name = c.Query("newName")
		}

		id, err := createDrinkInternal(db, c, claims, drink.DrinkData)
		if err != nil {
			return err
		}

		return c.JSON(CopyDrinkResponse{
			ID: id,
		})
	}
}
