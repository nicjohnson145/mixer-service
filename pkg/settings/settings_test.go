package settings

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"

	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/stretchr/testify/require"
)

func setupDbAndApp(t *testing.T) (*fiber.App, func()) {
	return commontest.SetupDbAndRouter(t, "settings.db", defineRoutes)
}

func t_getSettings(t *testing.T, app *fiber.App, o commontest.AuthOpts) (int, GetSettingsResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.SettingsV1,
		Auth:   &o,
	})
	return commontest.T_call_ok[GetSettingsResponse](t, app, req)
}

func t_writeSettings(t *testing.T, app *fiber.App, b WriteSettingsRequest, o commontest.AuthOpts) (int, WriteSettingsResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[WriteSettingsRequest]{
		Method: http.MethodPut,
		Path:   common.SettingsV1,
		Body:   &b,
		Auth:   &o,
	})
	return commontest.T_call_ok[WriteSettingsResponse](t, app, req)
}

func TestFullCRUDLoop(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	auth := commontest.AuthOpts{Username: commontest.Ptr("user1")}

	// Getting the settings from an empty DB should just yield the defaults
	status, resp := t_getSettings(t, app, auth)
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
	status, _ = t_writeSettings(t, app, WriteSettingsRequest{Settings: settings}, auth)
	require.Equal(t, 200, status)

	// Now get them
	status, resp = t_getSettings(t, app, auth)
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
	status, _ = t_writeSettings(t, app, WriteSettingsRequest{Settings: newSettings}, auth)
	require.Equal(t, 200, status)

	// Now get them
	status, resp = t_getSettings(t, app, auth)
	require.Equal(t, 200, status)
	require.Equal(
		t,
		UserSettings{
			PublicProfile: true,
		},
		resp.Settings,
	)
}
