package user

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type GetPublicUsersResponse struct {
	Users []string `json:"users"`
}

func getAllPublicUsers(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		users, err := getPublicUsers(db)
		if err != nil {
			return common.NewInternalServerErrorResp("getting list of public users", err)
		}

		return c.JSON(GetPublicUsersResponse{
			Users: users,
		})
	}
}
