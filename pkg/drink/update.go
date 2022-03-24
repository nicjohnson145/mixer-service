package drink

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"strconv"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type UpdateDrinkRequest struct {
	drinkData
}

type UpdateDrinkResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func updateDrink(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return err
		}
		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(UpdateDrinkResponse{
					Success: false,
				})
			} else {
				return err
			}
		}
		if model.Username != claims.Username {
			return c.Status(fiber.StatusNotFound).JSON(UpdateDrinkResponse{
				Success: false,
			})
		}

		var payload UpdateDrinkRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(UpdateDrinkResponse{
				Error:   err.Error(),
				Success: false,
			})
		}

		err = validate.Struct(payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(UpdateDrinkResponse{
				Error:   err.Error(),
				Success: false,
			})
		}

		drink := Drink{}
		setDrinkDataAttributes(&drink, payload)
		drink.ID = id
		drink.Username = claims.Username

		newModel, err := toDb(drink)
		if err != nil {
			return err
		}
		err = updateModel(newModel, db)
		if err != nil {
			return err
		}

		return c.JSON(UpdateDrinkResponse{
			Success: true,
		})
	}
}
