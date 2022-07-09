package settings

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"

	"github.com/stretchr/testify/require"
)

func setupDbAndApp(t *testing.T) (*fiber.App, func()) {
	return commontest.SetupDbAndRouter(t, "settings.db", defineRoutes)
}

func TestFullCRUDLoop(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	auth := commontest.AuthOpts{Username: commontest.Ptr("user1")}

	// Getting the settings from an empty DB should just yield the defaults
	status, resp := T_GetSettings(t, app, auth)
	require.Equal(t, 200, status)
	require.Equal(
		t,
		UserSettings{
			PublicProfile: true,
		},
		resp.Settings,
	)

	// Write some settings back to make sure they're persisted
	settings := UserSettings{
		PublicProfile: false,
	}
	status, _ = T_WriteSettings(t, app, WriteSettingsRequest{Settings: settings}, auth)
	require.Equal(t, 200, status)

	// Now get them
	status, resp = T_GetSettings(t, app, auth)
	require.Equal(t, 200, status)
	require.Equal(
		t,
		UserSettings{
			PublicProfile: false,
		},
		resp.Settings,
	)

	// Update them again, to show you can overwrite
	newSettings := UserSettings{
		PublicProfile: true,
	}
	status, _ = T_WriteSettings(t, app, WriteSettingsRequest{Settings: newSettings}, auth)
	require.Equal(t, 200, status)

	// Now get them
	status, resp = T_GetSettings(t, app, auth)
	require.Equal(t, 200, status)
	require.Equal(
		t,
		UserSettings{
			PublicProfile: true,
		},
		resp.Settings,
	)
}
