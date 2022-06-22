package settings

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
)

type GetSettingsResponse struct {
	Settings UserSettings `json:"settings"`
}

func getUserSettings(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims jwt.Claims) error {
		kvs, err := getByUsername(claims.Username, db)
		if err != nil {
			return common.NewInternalServerErrorResp(fmt.Sprintf("getting settings for user %v", claims.Username), err)
		}

		settings := UserSettings{}

		if val, ok := kvs[publicProfile]; ok {
			settings.PublicProfile = val == "true"
		} else {
			settings.PublicProfile = defaultPublicProfile
		}

		return c.JSON(GetSettingsResponse{
			Settings: settings,
		})
	}
}
