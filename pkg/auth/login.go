package auth

import (
	"github.com/gofiber/fiber/v2"
	"database/sql"
	"errors"
	"github.com/nicjohnson145/mixer-service/pkg/common"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Error        string `json:"error,omitempty"`
	Success      bool   `json:"success"`
	Username     string `json:"username,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func login(db *sql.DB) common.FiberHandler {
	return func(c *fiber.Ctx) error {
		var payload LoginRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(LoginResponse{
				Error: err.Error(),
				Success: false,
			})
		}

		existingUser, err := getUserByName(payload.Username, db)
		if err != nil {
			if errors.Is(err, common.ErrNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
					Error: err.Error(),
					Success: false,
					Username: payload.Username,
				})
			} else {
				return err
			}
		}

		if !comparePasswords(existingUser.Password, payload.Password) {
			return c.Status(fiber.StatusUnauthorized).JSON(LoginResponse{
				Error: "unauthorized",
				Success: false,
				Username: payload.Username,
			})
		}

		accessStr, err := GenerateAccessToken(TokenInputs{Username: payload.Username})
		if err != nil {
			return err
		}

		refreshStr, err := generateRefreshToken(TokenInputs{Username: payload.Username})
		if err != nil {
			return err
		}

		return c.JSON(LoginResponse{
			Success: true,
			Username: payload.Username,
			AccessToken: accessStr,
			RefreshToken: refreshStr,
		})
	}
}
