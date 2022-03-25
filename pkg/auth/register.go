package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Success bool `json:"success"`
}

func registerNewUser(db *sql.DB) FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		payload := new(RegisterNewUserRequest)
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		existingUser, err := getUserByName(payload.Username, db)
		if err != nil && !errors.Is(err, common.ErrNotFound) {
			return common.NewInternalServerErrorResp("checking DB for user", err)
		}

		if existingUser != nil {
			return common.ErrorResponse{
				Err:    fmt.Errorf("user already exists"),
				Msg:    fmt.Sprintf("user %v already exists", payload.Username),
				Status: fiber.StatusBadRequest,
			}
		}

		hashedPw, err := hashPassword(payload.Password)
		if err != nil {
			return common.NewInternalServerErrorResp("hashing password", err)
		}

		newUser := UserModel{
			Username: payload.Username,
			Password: hashedPw,
		}
		err = createUser(newUser, db)
		if err != nil {
			return common.NewInternalServerErrorResp("inserting into DB", err)
		}

		return c.JSON(RegisterNewUserResponse{
			Success: true,
		})
	}

}
