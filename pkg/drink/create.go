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
	DrinkData
}

type CreateDrinkResponse struct {
	ID int64 `json:"id,omitempty"`
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

		id, err := createDrinkInternal(db, c, claims, payload.DrinkData)
		if err != nil {
			return err
		}

		return c.JSON(CreateDrinkResponse{
			ID: id,
		})
	}
}

type DrinkAlreadyExistsErrorResponse struct {
	DrinkId int64
	Msg     string
	Status  int
}

func (e DrinkAlreadyExistsErrorResponse) Error() string {
	return e.Msg
}

func createDrinkInternal(db *sql.DB, c *fiber.Ctx, claims jwt.Claims, drinkData DrinkData) (int64, error) {
	existingDrink, err := getByNameAndUsername(drinkData.Name, claims.Username, db)
	if err != nil && !errors.Is(err, common.ErrNotFound) {
		return 0, common.NewInternalServerErrorResp("checking existance in DB", err)
	}

	if c.Query("overwrite") != "true" {
		if existingDrink != nil {
			return 0, common.NewBadRequestResponse(DrinkAlreadyExistsErrorResponse{
				DrinkId: existingDrink.ID,
				Msg:     fmt.Sprintf("user %v already has a drink named %v", claims.Username, drinkData.Name),
				Status:  fiber.StatusBadRequest,
			})
		}
	}

	drink := Drink{}
	setDrinkDataAttributes(&drink, drinkData)
	drink.Username = claims.Username

	model, err := toDb(drink)
	if err != nil {
		return 0, common.NewInternalServerErrorResp("converting to DB model", err)
	}

	var id int64
	if existingDrink != nil {
		model.ID = existingDrink.ID
		err = updateModel(model, db)
		id = existingDrink.ID
	} else {
		id, err = create(model, db)
	}
	if err != nil {
		return 0, common.NewInternalServerErrorResp("inserting into DB", err)
	}

	return id, nil
}
