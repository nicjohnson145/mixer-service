package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type RefreshTokenResponse struct {
	Error       string `json:"error,omitempty"`
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token,omitempty"`
}

func refresh() FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		newToken, err := jwt.GenerateAccessToken(jwt.TokenInputs{Username: claims.Username})
		if err != nil {
			return err
		}

		return c.JSON(RefreshTokenResponse{
			Success:     true,
			AccessToken: newToken,
		})
	}
}
