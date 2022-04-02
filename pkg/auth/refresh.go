package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToekn string `json:"refresh_token"`
}

func refresh() FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		newToken, err := jwt.GenerateAccessToken(jwt.TokenInputs{Username: claims.Username})
		if err != nil {
			return common.NewInternalServerErrorResp("generating access token", err)
		}

		newRefresh, err := jwt.GenerateRefreshToken(jwt.TokenInputs{Username: claims.Username})
		if err != nil {
			return common.NewInternalServerErrorResp("generating refresh token", err)
		}

		return c.JSON(RefreshTokenResponse{
			AccessToken:  newToken,
			RefreshToekn: newRefresh,
		})
	}
}
