package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username     string `json:"username,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func login(db *sql.DB) common.FiberHandler {
	return func(c *fiber.Ctx) error {
		var payload LoginRequest
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		existingUser, err := getUserByName(payload.Username, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return common.NewGenericUnauthorizedResponse(fmt.Sprintf("no user %v", payload.Username))
			} else {
				return common.NewInternalServerErrorResp(fmt.Sprintf("fetching %v from DB", payload.Username), err)
			}
		}

		if !comparePasswords(existingUser.Password, payload.Password) {
			return common.NewGenericUnauthorizedResponse(fmt.Sprintf("password mismatch for %v", payload.Username))
		}

		accessStr, err := jwt.GenerateAccessToken(jwt.TokenInputs{Username: payload.Username})
		if err != nil {
			return common.NewInternalServerErrorResp("generating access token", err)
		}

		refreshStr, err := jwt.GenerateRefreshToken(jwt.TokenInputs{Username: payload.Username})
		if err != nil {
			return common.NewInternalServerErrorResp("generating refresh token", err)
		}

		return c.JSON(LoginResponse{
			Username:     payload.Username,
			AccessToken:  accessStr,
			RefreshToken: refreshStr,
		})
	}
}
