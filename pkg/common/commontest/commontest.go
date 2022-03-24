package commontest

import (
	"database/sql"
	"io"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/nicjohnson145/mixer-service/pkg/jwt"
	"github.com/stretchr/testify/require"
)

func newDB(t *testing.T, name string) (*sql.DB, func()) {
	db, err := db.NewDB(name)
	require.NoError(t, err)

	cleanup := func() {
		err := os.Remove(name)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	return db, cleanup
}

func SetupDbAndRouter(t *testing.T, name string, routeFunc func(*fiber.App, *sql.DB)) (*fiber.App, func()) {
	db, cleanup := newDB(t, name)
	app := common.NewApp()
	routeFunc(app, db)
	return app, cleanup
}

func SetJsonHeader(r *http.Request) {
	r.Header.Set("Content-type", "application/json")
}

func RequireOkStatus(t *testing.T, resp *http.Response) {
	t.Helper()
	require.True(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)
}

func RequireNotOkStatus(t *testing.T, resp *http.Response) {
	t.Helper()
	require.False(t, resp.StatusCode >= 200 && resp.StatusCode <= 299)
}

type Req[T any] struct {
	Method string
	Path string
	Body *T
	Auth *AuthOpts
}


func T_call_ok[T any](t *testing.T, app *fiber.App, req *http.Request) (int, T) {
	t.Helper()
	resp, err := app.Test(req)
	require.NoError(t, err)
	RequireOkStatus(t, resp)

	defer resp.Body.Close()
	var rp T
	err = json.NewDecoder(resp.Body).Decode(&rp)
	require.NoError(t, err)

	return resp.StatusCode, rp
}

func T_call_fail(t *testing.T, app *fiber.App, req *http.Request) (int, common.OutboundErrResponse) {
	t.Helper()
	resp, err := app.Test(req)
	require.NoError(t, err)
	RequireNotOkStatus(t, resp)

	defer resp.Body.Close()
	var rp common.OutboundErrResponse
	err = json.NewDecoder(resp.Body).Decode(&rp)
	require.NoError(t, err)

	return resp.StatusCode, rp
}

func T_req[T any](t *testing.T, r Req[T]) *http.Request {
	t.Helper()
	var body io.Reader
	if r.Body != nil {
		bodyBytes, err := json.Marshal(r.Body)
		require.NoError(t, err)
		body = strings.NewReader(string(bodyBytes))
	}

	req, err := http.NewRequest(
		r.Method,
		r.Path,
		body,
	)
	require.NoError(t, err)

	if r.Auth != nil {
		AuthenticatedRequest(t, req, *r.Auth)
	}
	if r.Body != nil {
		SetJsonHeader(req)
	}

	return req
}

const (
	DefaultUsername = "foobar"
)

type AuthOpts struct {
	Username *string
}

func AuthenticatedRequest(t *testing.T, r *http.Request, opts AuthOpts) {
	if opts.Username == nil {
		opts.Username = to.StringPtr(DefaultUsername)
	}

	token, err := jwt.GenerateAccessToken(jwt.TokenInputs{
		Username: *opts.Username,
	})
	require.NoError(t, err)

	r.Header.Set(jwt.AuthenticationHeader, token)
}
