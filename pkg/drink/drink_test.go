package drink

import (
	"database/sql"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/nicjohnson145/mixer-service/pkg/auth/authtest"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
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

func setupDbAndRouter(t *testing.T) (*mux.Router, func()) {
		db, cleanup := newDB(t)
		router := mux.NewRouter()
		defineRoutes(router, db)
		return router, cleanup
}

func postCreateDrink(t *testing.T, router *mux.Router, r CreateDrinkRequest, o authtest.AuthOpts) (int, CreateDrinkResponse) {
	bodyBytes, err := json.Marshal(r)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPost,
		common.DrinksV1+"/create",
		strings.NewReader(string(bodyBytes)),
	)
	require.NoError(t, err)
	authtest.AuthenticatedRequest(t, req, o)

	router.ServeHTTP(rr, req)

	defer rr.Result().Body.Close()
	var resp CreateDrinkResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&resp)
	require.NoError(t, err)

	return rr.Result().StatusCode, resp
}

func getDrinkByID(t *testing.T, router *mux.Router, id int64, o authtest.AuthOpts) (int, GetDrinkResponse) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		common.DrinksV1+fmt.Sprintf("/%v", id),
		nil,
	)
	require.NoError(t, err)
	authtest.AuthenticatedRequest(t, req, o)
	router.ServeHTTP(rr, req)

	defer rr.Result().Body.Close()
	var resp GetDrinkResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&resp)
	require.NoError(t, err)

	return rr.Result().StatusCode, resp
}

func TestCreateGet(t *testing.T) {
	t.Run("create_get_happy", func(t *testing.T) {
		router, cleanup := setupDbAndRouter(t)
		defer cleanup()

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

		status, resp := postCreateDrink(t, router, body, authtest.AuthOpts{})
		require.Equal(t, http.StatusOK, status)

		status, getResp := getDrinkByID(t, router, resp.ID, authtest.AuthOpts{})
		require.Equal(t, http.StatusOK, status)

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
		router, cleanup := setupDbAndRouter(t)
		defer cleanup()

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

		status, resp := postCreateDrink(t, router, body, authtest.AuthOpts{Username: to.StringPtr("user1")})
		require.Equal(t, http.StatusOK, status)

		status, _ = getDrinkByID(t, router, resp.ID, authtest.AuthOpts{Username: to.StringPtr("user2")})
		require.Equal(t, http.StatusNotFound, status)
	})
}
