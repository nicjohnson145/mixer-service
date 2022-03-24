package drink

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type CreateDrinkRequest struct {
	drinkData
}

type CreateDrinkResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	ID      int64  `json:"id,omitempty"`
}

func createDrink(db *sql.DB) auth.FiberClaimsHandler {

	return func(c *fiber.Ctx, claims jwt.Claims) error {
		var payload CreateDrinkRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(CreateDrinkResponse{
				Success: false,
			})
		}
		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(CreateDrinkResponse{
				Error:   err.Error(),
				Success: false,
			})
		}

		existingDrink, err := getByNameAndUsername(payload.Name, claims.Username, db)
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			return err
		}

		if existingDrink != nil {
			return c.Status(fiber.StatusBadRequest).JSON(CreateDrinkResponse{
				Error:   fmt.Sprintf("user %v already has a drink named %v", claims.Username, payload.Name),
				Success: false,
			})
		}

		drink := Drink{}
		setDrinkDataAttributes(&drink, payload)
		drink.Username = claims.Username

		model, err := toDb(drink)
		if err != nil {
			return err
		}

		id, err := create(model, db)
		if err != nil {
			return err
		}

		return c.JSON(CreateDrinkResponse{
			Success: true,
			ID:      id,
		})
	}
}
