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
	Success bool   `json:"success"`
}

func updateDrink(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return common.NewBadRequestResponse(err)
		}
		model, err := getByID(id, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return common.NewGenericNotFoundResponse("getting drink from DB")
			} else {
				return common.NewInternalServerErrorResp("getting drink from DB", err)
			}
		}
		if model.Username != claims.Username {
			return common.NewGenericNotFoundResponse("username mismatch")
		}

		var payload UpdateDrinkRequest
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		err = validate.Struct(payload)
		if err != nil {
			return common.NewBadRequestResponse(err)
		}

		drink := Drink{}
		setDrinkDataAttributes(&drink, payload)
		drink.ID = id
		drink.Username = claims.Username

		newModel, err := toDb(drink)
		if err != nil {
			return common.NewInternalServerErrorResp("converting to DB model", err)
		}
		err = updateModel(newModel, db)
		if err != nil {
			return common.NewInternalServerErrorResp("updating in DB", err)
		}

		return c.JSON(UpdateDrinkResponse{
			Success: true,
		})
	}
}
