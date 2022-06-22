package settings

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type WriteSettingsResponse struct {
	Success bool `json:"success"`
}

type WriteSettingsRequest struct {
	Settings UserSettings `json:"settings"`
}

func writeUserSettings(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		var payload WriteSettingsRequest
		if err := c.BodyParser(&payload); err != nil {
			return common.NewBadRequestResponse(err)
		}

		err := writeSettingsForUser(db, claims.Username, payload.Settings)
		if err != nil {
			return common.NewInternalServerErrorResp("upserting settings to DB", err)
		}

		return c.JSON(WriteSettingsResponse{
			Success: true,
		})
	}
}
