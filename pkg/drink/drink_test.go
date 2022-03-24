package drink

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/gofiber/fiber/v2"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/common/commontest"
	"github.com/stretchr/testify/require"
)

func setupDbAndApp(t *testing.T) (*fiber.App, func()) {
	return commontest.SetupDbAndRouter(t, "drink.db", defineRoutes)
}

func t_createDrink_ok(t *testing.T, app *fiber.App, r CreateDrinkRequest, o commontest.AuthOpts) (int, CreateDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CreateDrinkRequest]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + "/create",
		Body:   &r,
		Auth:   &o,
	})
	return commontest.T_call_ok[CreateDrinkResponse](t, app, req)
}

func t_createDrink_fail(t *testing.T, app *fiber.App, r CreateDrinkRequest, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CreateDrinkRequest]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + "/create",
		Body:   &r,
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func t_getDrink_ok(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, GetDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Auth:   &o,
	})
	return commontest.T_call_ok[GetDrinkResponse](t, app, req)
}

func t_getDrink_fail(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func t_updateDrink_ok(t *testing.T, app *fiber.App, id int64, r UpdateDrinkRequest, o commontest.AuthOpts) (int, UpdateDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[UpdateDrinkRequest]{
		Method: http.MethodPut,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Body:   &r,
		Auth:   &o,
	})
	return commontest.T_call_ok[UpdateDrinkResponse](t, app, req)
}

func t_updateDrink_fail(t *testing.T, app *fiber.App, id int64, r UpdateDrinkRequest, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[UpdateDrinkRequest]{
		Method: http.MethodPut,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Body:   &r,
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func t_deleteDrink_ok(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, DeleteDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodDelete,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Auth:   &o,
	})
	return commontest.T_call_ok[DeleteDrinkResponse](t, app, req)
}

func t_deleteDrink_fail(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodDelete,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id),
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func t_getDrinksByUser_ok(t *testing.T, app *fiber.App, user string, o commontest.AuthOpts) (int, GetDrinksByUserResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.DrinksV1 + "/by-user/" + user,
		Auth:   &o,
	})
	return commontest.T_call_ok[GetDrinksByUserResponse](t, app, req)
}

func t_getDrinksByUser_fail(t *testing.T, app *fiber.App, user string, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[any]{
		Method: http.MethodGet,
		Path:   common.DrinksV1 + "/by-user/" + user,
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func TestFullCRUDLoop(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
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
	_, createResp := t_createDrink_ok(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})

	// Fetch it as the orignal author
	_, getResp := t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp := GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Fetch it as someone else should fail
	_, _ = t_getDrink_fail(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})

	// Update it
	updateReq := UpdateDrinkRequest{
		drinkData: updatedDrinkData,
	}
	// Updating as someone else should not work
	_, _ = t_updateDrink_fail(t, app, createResp.ID, updateReq, commontest.AuthOpts{Username: to.StringPtr("user2")})
	// Should still be the same
	_, getResp = t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Update it as the original author should work
	_, _ = t_updateDrink_ok(t, app, createResp.ID, updateReq, commontest.AuthOpts{Username: to.StringPtr("user1")})
	_, getResp = t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Deleting it as someone else should not be possible
	_, _ = t_deleteDrink_fail(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	_, getResp = t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// But deleting it as the orignal author should work
	_, _ = t_deleteDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	_, _ = t_getDrink_fail(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
}

func TestPublicDrinksFetchableByAnyone(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
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
	_, createResp := t_createDrink_ok(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})

	// Fetch it as someone else, it should succeed since it's public
	_, getResp := t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	expectedGetResp := GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			drinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}

func TestGetDrinksByUser(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
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
		_, _ = t_createDrink_ok(t, app, CreateDrinkRequest{drinkData: d}, commontest.AuthOpts{Username: to.StringPtr("user1")})
	}

	// Fetching as user1 should result in all drinks
	_, getResp := t_getDrinksByUser_ok(t, app, "user1", commontest.AuthOpts{Username: to.StringPtr("user1")})
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
	_, getResp = t_getDrinksByUser_ok(t, app, "user1", commontest.AuthOpts{Username: to.StringPtr("user2")})
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
