package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}

func refresh() FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		newToken, err := jwt.GenerateAccessToken(jwt.TokenInputs{Username: claims.Username})
		if err != nil {
			return common.NewInternalServerErrorResp("generating access token", err)
		}

		return c.JSON(RefreshTokenResponse{
			AccessToken: newToken,
		})
	}
}
