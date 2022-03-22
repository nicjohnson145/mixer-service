package drink

import (
	"database/sql"
	"errors"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type DeleteDrinkResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func deleteDrink(db *sql.DB) auth.FiberClaimsHandler {

	return func(c *fiber.Ctx, claims auth.Claims) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return err
		}
		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(DeleteDrinkResponse{
					Success: false,
				})
			} else {
				return err
			}
		}
		if model.Username != claims.Username {
			return c.Status(fiber.StatusNotFound).JSON(DeleteDrinkResponse{
				Success: false,
			})
		}

		err = deleteModel(id, db)
		if err != nil {
			return err
		}

		return c.JSON(DeleteDrinkResponse{
			Success: true,
		})
	}
}
