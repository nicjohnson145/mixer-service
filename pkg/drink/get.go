package drink

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type GetDrinkResponse struct {
	Drink *Drink `json:"drink"`
}

func getDrink(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {

		drink, err := getDrinkInternal(db, c, claims)
		if err != nil {
			return err
		}

		if drink.Username != claims.Username && drink.Publicity != DrinkPublicityPublic {
			return common.NewGenericNotFoundResponse("non-public access")
		}

		return c.JSON(GetDrinkResponse{
			Drink: drink,
		})
	}
}

func getDrinkInternal(db *sql.DB, c *fiber.Ctx, claims jwt.Claims) (*Drink, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return nil, common.NewBadRequestResponse(err)
	}
	model, err := getByID(id, db)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.NewGenericNotFoundResponse("not found in DB")
		} else {
			return nil, common.NewInternalServerErrorResp("getting drink from DB", err)
		}
	}
	drink, err := fromDb(*model)
	if err != nil {
		return nil, common.NewInternalServerErrorResp("converting from DB model", err)
	}

	return &drink, nil
}
