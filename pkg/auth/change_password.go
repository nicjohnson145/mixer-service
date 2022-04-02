package auth

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type ChangePasswordResponse struct {
	Success bool `json:"success"`
}

func changePassword(db *sql.DB) FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		payload := new(ChangePasswordRequest)
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		hashedPw, err := hashPassword(payload.NewPassword)
		if err != nil {
			return common.NewInternalServerErrorResp("hashing password", err)
		}

		err = updatePassword(UserModel{Username: claims.Username, Password: hashedPw}, db)
		if err != nil {
			return common.NewInternalServerErrorResp("updating DB", err)
		}

		return c.JSON(ChangePasswordResponse{Success: true})
	}
}
