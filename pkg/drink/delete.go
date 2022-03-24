package drink

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	"strconv"
)

type DeleteDrinkResponse struct {
	Success bool `json:"success"`
}

func deleteDrink(db *sql.DB) auth.FiberClaimsHandler {

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

		err = deleteModel(id, db)
		if err != nil {
			return common.NewInternalServerErrorResp("deleting model from DB", err)
		}

		return c.JSON(DeleteDrinkResponse{
			Success: true,
		})
	}
}
