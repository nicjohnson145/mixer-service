package user

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
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

func t_registerUser(t *testing.T, app *fiber.App, r auth.RegisterNewUserRequest) (int, auth.RegisterNewUserResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[auth.RegisterNewUserRequest]{
		Method: http.MethodPost,
		Path:   common.AuthV1 + "/register-user",
		Body:   &r,
	})
	return commontest.T_call_ok[auth.RegisterNewUserResponse](t, app, req)
}

func t_writeSettings(t *testing.T, app *fiber.App, b settings.WriteSettingsRequest, o commontest.AuthOpts) (int, settings.WriteSettingsResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[settings.WriteSettingsRequest]{
		Method: http.MethodPut,
		Path:   common.SettingsV1,
		Body:   &b,
		Auth:   &o,
	})
	return commontest.T_call_ok[settings.WriteSettingsResponse](t, app, req)
}

func t_getPublicUsers(t *testing.T, app *fiber.App, o commontest.AuthOpts) (int, GetPublicUsersResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.UserV1,
		Auth:   &o,
	})
	return commontest.T_call_ok[GetPublicUsersResponse](t, app, req)
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
		status, _ := t_registerUser(t, app, auth.RegisterNewUserRequest{Username: name, Password: "foo"})
		require.Equal(t, 200, status)
	}

	// Setup the user settings for the explicit users
	status, _ := t_writeSettings(
		t,
		app,
		settings.WriteSettingsRequest{Settings: settings.UserSettings{PublicProfile: true}},
		commontest.AuthOpts{Username: commontest.Ptr("explicitly_public_user")},
	)
	require.Equal(t, 200, status)

	status, _ = t_writeSettings(
		t,
		app,
		settings.WriteSettingsRequest{Settings: settings.UserSettings{PublicProfile: false}},
		commontest.AuthOpts{Username: commontest.Ptr("explicitly_private_user")},
	)
	require.Equal(t, 200, status)

	status, resp := t_getPublicUsers(t, app, commontest.AuthOpts{Username: commontest.Ptr("test_user")})
	require.Equal(t, 200, status)
	sort.Strings(resp.Users)
	require.Equal(t, []string{"explicitly_public_user", "no_settings_user"}, resp.Users)
}

