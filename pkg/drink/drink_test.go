package drink

import (
	"database/sql"
	"encoding/json"
	"github.com/nicjohnson145/mixer-service/pkg/auth/authtest"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/Azure/go-autorest/autorest/to"
	// log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func newDB(t *testing.T) (*sql.DB, func()) {
	const name = "drink.db"
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

func TestCreateGet(t *testing.T) {
	t.Run("create_get_happy", func(t *testing.T) {
		db, cleanup := newDB(t)
		defer cleanup()

		router := mux.NewRouter()
		defineRoutes(router, db)

		body := CreateDrinkRequest{
			Name:           "Daquari",
			PrimaryAlcohol: "Rum",
			PreferredGlass: "Coupe",
			Ingredients: []string{
				"2.5 oz white rum",
				"0.5 oz simple syrup",
				"1 oz lime",
			},
		}

		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodPost,
			common.DrinksV1+"/create",
			strings.NewReader(string(bodyBytes)),
		)
		require.NoError(t, err)
		authtest.AuthenticatedRequest(t, req, authtest.AuthOpts{})

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)

		defer rr.Result().Body.Close()
		var resp CreateDrinkResponse
		err = json.NewDecoder(rr.Result().Body).Decode(&resp)
		require.NoError(t, err)

		// Now retrieve the drink
		rr = httptest.NewRecorder()
		req, err = http.NewRequest(
			http.MethodGet,
			common.DrinksV1+fmt.Sprintf("/%v", resp.ID),
			nil,
		)
		require.NoError(t, err)
		authtest.AuthenticatedRequest(t, req, authtest.AuthOpts{})
		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)

		defer rr.Result().Body.Close()
		var getResp GetDrinkResponse
		err = json.NewDecoder(rr.Result().Body).Decode(&getResp)
		require.NoError(t, err)

		expectedDrink := &Drink{
			Name:           "Daquari",
			Username:       authtest.DefaultUsername,
			PrimaryAlcohol: "Rum",
			PreferredGlass: "Coupe",
			Ingredients: []string{
				"2.5 oz white rum",
				"0.5 oz simple syrup",
				"1 oz lime",
			},
		}
		require.Equal(t, expectedDrink, getResp.Drink)
	})

	t.Run("fetch_other_users_drink", func(t *testing.T) {
		db, cleanup := newDB(t)
		defer cleanup()

		router := mux.NewRouter()
		defineRoutes(router, db)

		body := CreateDrinkRequest{
			Name:           "Daquari",
			PrimaryAlcohol: "Rum",
			PreferredGlass: "Coupe",
			Ingredients: []string{
				"2.5 oz white rum",
				"0.5 oz simple syrup",
				"1 oz lime",
			},
		}

		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodPost,
			common.DrinksV1+"/create",
			strings.NewReader(string(bodyBytes)),
		)
		require.NoError(t, err)
		authtest.AuthenticatedRequest(t, req, authtest.AuthOpts{Username: to.StringPtr("user1")})

		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Result().StatusCode)

		defer rr.Result().Body.Close()
		var resp CreateDrinkResponse
		err = json.NewDecoder(rr.Result().Body).Decode(&resp)
		require.NoError(t, err)

		// Now retrieve the drink
		rr = httptest.NewRecorder()
		req, err = http.NewRequest(
			http.MethodGet,
			common.DrinksV1+fmt.Sprintf("/%v", resp.ID),
			nil,
		)
		require.NoError(t, err)
		authtest.AuthenticatedRequest(t, req, authtest.AuthOpts{Username: to.StringPtr("user2")})
		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNotFound, rr.Result().StatusCode)
	})
}
