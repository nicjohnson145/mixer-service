package user

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/nicjohnson145/mixer-service/pkg/settings"

	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/stretchr/testify/require"
	"sort"
)

func setupDbAndRouter(t *testing.T) (*fiber.App, func()) {
	t.Helper()

	db, cleanup := commontest.NewDB(t, "user.db")
	app := common.NewApp()
	defineRoutes(app, db)

	err := auth.Init(app, db)
	require.NoError(t, err)

	err = settings.Init(app, db)
	require.NoError(t, err)

	return app, cleanup
}

func TestGetAllPublicUsers(t *testing.T) {
	app, cleanup := setupDbAndRouter(t)
	defer cleanup()

	// Register some users
	allUsers := []string{
		"explicitly_public_user",
		"no_settings_user",
		"explicitly_private_user",
	}
	for _, name := range allUsers {
		status, _ := auth.T_RegisterUser(t, app, auth.RegisterNewUserRequest{Username: name, Password: "foo"})
		require.Equal(t, 200, status)
	}

	// Setup the user settings for the explicit users
	status, _ := settings.T_WriteSettings(
		t,
		app,
		settings.WriteSettingsRequest{Settings: settings.UserSettings{PublicProfile: true}},
		commontest.AuthOpts{Username: commontest.Ptr("explicitly_public_user")},
	)
	require.Equal(t, 200, status)

	status, _ = settings.T_WriteSettings(
		t,
		app,
		settings.WriteSettingsRequest{Settings: settings.UserSettings{PublicProfile: false}},
		commontest.AuthOpts{Username: commontest.Ptr("explicitly_private_user")},
	)
	require.Equal(t, 200, status)

	status, resp := T_GetPublicUsers(t, app, commontest.AuthOpts{Username: commontest.Ptr("test_user")})
	require.Equal(t, 200, status)
	sort.Strings(resp.Users)
	require.Equal(t, []string{"explicitly_public_user", "no_settings_user"}, resp.Users)
}
