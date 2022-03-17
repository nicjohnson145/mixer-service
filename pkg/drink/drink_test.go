package drink

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/gorilla/mux"
	"github.com/nicjohnson145/mixer-service/pkg/auth/authtest"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupDbAndRouter(t *testing.T) (*mux.Router, func()) {
	return commontest.SetupDbAndRouter(t, "drink.db", defineRoutes)
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

func t_updateDrink(t *testing.T, router *mux.Router, id int64, r UpdateDrinkRequest, o authtest.AuthOpts) (int, UpdateDrinkResponse) {
	bodyBytes, err := json.Marshal(r)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodPut,
		common.DrinksV1+fmt.Sprintf("/%v", id),
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

func t_getDrinksByUser(t *testing.T, router *mux.Router, user string, o authtest.AuthOpts) (int, GetDrinksByUserResponse) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(
		http.MethodGet,
		common.DrinksV1+"/by-user/"+user,
		nil,
	)
	require.NoError(t, err)
	authtest.AuthenticatedRequest(t, req, o)
	router.ServeHTTP(rr, req)

	defer rr.Result().Body.Close()
	var resp GetDrinksByUserResponse
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
		Publicity: DrinkPublicityPrivate,
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
		Publicity: DrinkPublicityPrivate,
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
			ID:        1,
			Username:  "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Fetch it as someone else should fail
	status, _ = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusNotFound, status)

	// Update it
	updateReq := UpdateDrinkRequest{
		drinkData: updatedDrinkData,
	}
	// Updating as someone else should not work
	status, _ = t_updateDrink(t, router, createResp.ID, updateReq, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusNotFound, status)
	// Should still be the same
	status, getResp = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Update it as the original author should work
	status, _ = t_updateDrink(t, router, createResp.ID, updateReq, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	status, getResp = t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
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
			ID:        1,
			Username:  "user1",
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

func TestPublicDrinksFetchableByAnyone(t *testing.T) {
	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	drinkData := drinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPublic,
	}

	body := CreateDrinkRequest{drinkData: drinkData}

	// Creating a drink
	status, createResp := t_createDrink(t, router, body, authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)

	// Fetch it as someone else, it should succeed since it's public
	status, getResp := t_getDrink(t, router, createResp.ID, authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp := GetDrinkResponse{
		Success: true,
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}

func TestGetDrinksByUser(t *testing.T) {
	router, cleanup := setupDbAndRouter(t)
	defer cleanup()

	first := drinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPublic,
	}
	second := drinkData{
		Name:           "Bee's Knees",
		PrimaryAlcohol: "Gin",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"0.75 oz honey syrup",
			"0.75 oz lemon",
			"2 oz gin",
		},
		Publicity: DrinkPublicityPublic,
	}
	third := drinkData{
		Name:           "Secret Drink",
		PrimaryAlcohol: "Scotch",
		PreferredGlass: "Rocks",
		Ingredients: []string{
			"2 oz scotch",
		},
		Publicity: DrinkPublicityPrivate,
	}

	drinks := []drinkData{first, second, third}
	for _, d := range drinks {
		status, _ := t_createDrink(t, router, CreateDrinkRequest{drinkData: d}, authtest.AuthOpts{Username: to.StringPtr("user1")})
		require.Equal(t, http.StatusOK, status)
	}

	// Fetching as user1 should result in all drinks
	status, getResp := t_getDrinksByUser(t, router, "user1", authtest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp := GetDrinksByUserResponse{
		Success: true,
		Drinks: []Drink{
			{
				ID:        1,
				Username:  "user1",
				drinkData: first,
			},
			{
				ID:        2,
				Username:  "user1",
				drinkData: second,
			},
			{
				ID:        3,
				Username:  "user1",
				drinkData: third,
			},
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Fetching as user2 should only return the public drinks
	status, getResp = t_getDrinksByUser(t, router, "user1", authtest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusOK, status)
	expectedGetResp = GetDrinksByUserResponse{
		Success: true,
		Drinks: []Drink{
			{
				ID:        1,
				Username:  "user1",
				drinkData: first,
			},
			{
				ID:        2,
				Username:  "user1",
				drinkData: second,
			},
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}
