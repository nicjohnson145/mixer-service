package drink

import (
	"database/sql"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/nicjohnson145/mixer-service/pkg/auth/authtest"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
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

func t_createDrink(t *testing.T, router *mux.Router, r CreateDrinkRequest, o authtest.AuthOpts) (int, CreateDrinkResponse) {
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

func t_getDrink(t *testing.T, router *mux.Router, id int64, o authtest.AuthOpts) (int, GetDrinkResponse) {
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

func t_updateDrink(t *testing.T, router *mux.Router, r UpdateDrinkRequest, o authtest.AuthOpts) (int, UpdateDrinkResponse) {
	bodyBytes, err := json.Marshal(r)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut,
		common.DrinksV1+fmt.Sprintf("/%v", r.ID),
		strings.NewReader(string(bodyBytes)),
	)
	require.NoError(t, err)
	authtest.AuthenticatedRequest(t, req, o)
	router.ServeHTTP(rr, req)

	defer rr.Result().Body.Close()
	var resp UpdateDrinkResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&resp)
	require.NoError(t, err)

	return rr.Result().StatusCode, resp
}

func t_deleteDrink(t *testing.T, router *mux.Router, id int64, o authtest.AuthOpts) (int, DeleteDrinkResponse) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodDelete,
		common.DrinksV1+fmt.Sprintf("/%v", id),
		nil,
	)
	require.NoError(t, err)
	authtest.AuthenticatedRequest(t, req, o)
	router.ServeHTTP(rr, req)

	defer rr.Result().Body.Close()
	var resp DeleteDrinkResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&resp)
	require.NoError(t, err)

	return rr.Result().StatusCode, resp
}

func TestFullCRUDLoop(t *testing.T) {
	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	origDrinkData := drinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"1 oz lime",
		},
	}
	updatedDrinkData := drinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"0.75 oz lime",
		},
	}

	body := CreateDrinkRequest{drinkData: origDrinkData}

	// Creating a drink
	status, createResp := t_createDrink(t, router, body, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)

	// Fetch it as the orignal author
	status, getResp := t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp := GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID: 1,
			Username: "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Fetch it as someone else should fail
	status, _ = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusNotFound, status)

	// Update it
	updateReq := UpdateDrinkRequest{
		ID: createResp.ID,
		drinkData: updatedDrinkData,
	}
	// Updating as someone else should not work
	status, _ = t_updateDrink(t, router, updateReq, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusNotFound, status)
	// Should still be the same
	status, getResp = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID: 1,
			Username: "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Update it as the original author should work
	status, _ = t_updateDrink(t, router, updateReq, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	status, getResp = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID: 1,
			Username: "user1",
			drinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Deleting it as someone else should not be possible
	status, _ = t_deleteDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusNotFound, status)
	status, getResp = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID: 1,
			Username: "user1",
			drinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// But deleting it as the orignal author should work
	status, _ = t_deleteDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	status, _ = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusNotFound, status)
}
