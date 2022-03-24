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
	ID      int64  `json:"id,omitempty"`
}

func createDrink(db *sql.DB) auth.FiberClaimsHandler {

	return func(c *fiber.Ctx, claims jwt.Claims) error {
		var payload CreateDrinkRequest
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}
		if err := validate.Struct(payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		existingDrink, err := getByNameAndUsername(payload.Name, claims.Username, db)
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			return common.NewInternalServerErrorResp("checking existance in DB", err)
		}

		if existingDrink != nil {
			return common.ErrorResponse{
				Msg: fmt.Sprintf("user %v already has a drink named %v", claims.Username, payload.Name),
				Err: fmt.Errorf("exising name/user combination"),
				Status: fiber.StatusBadRequest,
			}
		}

		drink := Drink{}
		setDrinkDataAttributes(&drink, payload)
		drink.Username = claims.Username

		model, err := toDb(drink)
		if err != nil {
			return common.NewInternalServerErrorResp("converting to DB model", err)
		}

		id, err := create(model, db)
		if err != nil {
			return common.NewInternalServerErrorResp("inserting into DB", err)
		}

		return c.JSON(CreateDrinkResponse{
			ID:      id,
		})
	}
}
