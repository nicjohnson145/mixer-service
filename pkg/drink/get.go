package drink

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"strconv"
)

type GetDrinkResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
	Drink   *Drink `json:"drink"`
}

func getDrink(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims auth.Claims) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return err
		}
		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(GetDrinkResponse{
					Success: false,
				})
			} else {
				return err
			}
		}

		drink, err := fromDb(*model)
		if err != nil {
			return err
		}

		if drink.Username != claims.Username && drink.Publicity != DrinkPublicityPublic {
			return c.Status(fiber.StatusNotFound).JSON(GetDrinkResponse{
				Success: false,
			})
		}

		return c.JSON(GetDrinkResponse{
			Success: true,
			Drink:   &drink,
		})
	}
}
