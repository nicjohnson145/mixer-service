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

func t_createDrinkOverwrite_ok(t *testing.T, app *fiber.App, r CreateDrinkRequest, o commontest.AuthOpts) (int, CreateDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CreateDrinkRequest]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + "/create?overwrite=true",
		Body:   &r,
		Auth:   &o,
	})
	return commontest.T_call_ok[CreateDrinkResponse](t, app, req)
}

func t_createDrink_DrinkAlreadyExists(t *testing.T, app *fiber.App, r CreateDrinkRequest, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
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

func t_copyDrink_ok(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, CopyDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CopyDrinkResponse]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id) + "/copy",
		Auth:   &o,
	})
	return commontest.T_call_ok[CopyDrinkResponse](t, app, req)
}

func t_copyDrinkRename_ok(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, CopyDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CopyDrinkResponse]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id) + "/copy?newName=DrinkCopy",
		Auth:   &o,
	})
	return commontest.T_call_ok[CopyDrinkResponse](t, app, req)
}

func t_copyDrinkOverwrite_ok(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, CopyDrinkResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CopyDrinkResponse]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id) + "/copy?overwrite=true",
		Auth:   &o,
	})
	return commontest.T_call_ok[CopyDrinkResponse](t, app, req)
}

func t_copyDrink_DrinkAlreadyExists(t *testing.T, app *fiber.App, id int64, o commontest.AuthOpts) (int, common.OutboundErrResponse) {
	t.Helper()
	req := commontest.T_req(t, commontest.Req[CopyDrinkResponse]{
		Method: http.MethodPost,
		Path:   common.DrinksV1 + fmt.Sprintf("/%v", id) + "/copy",
		Auth:   &o,
	})
	return commontest.T_call_fail(t, app, req)
}

func TestFullCRUDLoop(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	origDrinkData := DrinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPrivate,
		UnderDevelopment: true,
	}
	updatedDrinkData := DrinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"0.75 oz lime",
		},
		Publicity: DrinkPublicityPrivate,
		UnderDevelopment: false,
	}

	body := CreateDrinkRequest{DrinkData: origDrinkData}

	// Creating a drink
	_, createResp := t_createDrink_ok(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})

	// Fetch it as the orignal author
	_, getResp := t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp := GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			DrinkData: origDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Fetch it as someone else should fail
	_, _ = t_getDrink_fail(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})

	// Update it
	updateReq := UpdateDrinkRequest{
		DrinkData: updatedDrinkData,
	}
	// Updating as someone else should not work
	_, _ = t_updateDrink_fail(t, app, createResp.ID, updateReq, commontest.AuthOpts{Username: to.StringPtr("user2")})
	// Should still be the same
	_, getResp = t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			DrinkData: origDrinkData,
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
			DrinkData: updatedDrinkData,
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
			DrinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// But deleting it as the orignal author should work
	_, _ = t_deleteDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	_, _ = t_getDrink_fail(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
}

func TestCreateWithOverwriteAndGet(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	origDrinkData := DrinkData{
		Name:           "Daquiri",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz simple syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPrivate,
	}
	overwriteDrinkData := DrinkData{
		Name:           "Daquiri",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz strawberry syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPrivate,
	}
	body := CreateDrinkRequest{DrinkData: origDrinkData}
	overwriteBody := CreateDrinkRequest{DrinkData: overwriteDrinkData}

	// Creating a drink
	_, createResp := t_createDrink_ok(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})
	// Create it again with overwrite
	_, overwriteResp := t_createDrinkOverwrite_ok(t, app, overwriteBody, commontest.AuthOpts{Username: to.StringPtr("user1")})

	require.Equal(t, createResp.ID, overwriteResp.ID)

	// Fetch it and it should be eqwual to the overwritten value
	_, getResp := t_getDrink_ok(t, app, overwriteResp.ID, commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp := GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			DrinkData: overwriteDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}

func TestPublicDrinksFetchableAndCopyableByAnyone(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	drinkData := DrinkData{
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

	updatedDrinkData := DrinkData{
		Name:           "Daquari",
		PrimaryAlcohol: "Rum",
		PreferredGlass: "Coupe",
		Ingredients: []string{
			"2.5 oz white rum",
			"0.5 oz strawberry syrup",
			"1 oz lime",
		},
		Publicity: DrinkPublicityPublic,
	}

	body := CreateDrinkRequest{DrinkData: drinkData}

	// Creating a drink
	_, createResp := t_createDrink_ok(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})
	// Creating exact same drink again should fail
	status, errResp := t_createDrink_DrinkAlreadyExists(t, app, body, commontest.AuthOpts{Username: to.StringPtr("user1")})
	require.Equal(t, http.StatusConflict, status)
	require.Equal(t, "existing drink named Daquari", errResp.Error)

	// Fetch it as someone else, it should succeed since it's public
	_, getResp := t_getDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	expectedGetResp := GetDrinkResponse{
		Drink: &Drink{
			ID:        1,
			Username:  "user1",
			DrinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Copy it as someone else
	_, copyResp := t_copyDrink_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	// Copy again without overwrite or rename should fail
	status, errResp = t_copyDrink_DrinkAlreadyExists(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	require.Equal(t, http.StatusConflict, status)
	require.Equal(t, "existing drink named Daquari", errResp.Error)

	// Get it and it should be the same as the original but with a new owner and id
	_, getResp = t_getDrink_ok(t, app, copyResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        2,
			Username:  "user2",
			DrinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// update it (so we can copy it again with overwrite)
	// Update it
	updateReq := UpdateDrinkRequest{
		DrinkData: updatedDrinkData,
	}
	_, _ = t_updateDrink_ok(t, app, 2, updateReq, commontest.AuthOpts{Username: to.StringPtr("user2")})
	// Get it and it should be updated
	_, getResp = t_getDrink_ok(t, app, copyResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        2,
			Username:  "user2",
			DrinkData: updatedDrinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)
	// Copy it again with overwrite
	_, copyResp = t_copyDrinkOverwrite_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	// Get it and it should be back to the original
	_, getResp = t_getDrink_ok(t, app, copyResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        2,
			Username:  "user2",
			DrinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)

	// Copy it again with new name
	_, copyResp = t_copyDrinkRename_ok(t, app, createResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})

	// Get it and it should be the same as the original but with a new name and id
	_, getResp = t_getDrink_ok(t, app, copyResp.ID, commontest.AuthOpts{Username: to.StringPtr("user2")})
	drinkData.Name = "DrinkCopy"
	expectedGetResp = GetDrinkResponse{
		Drink: &Drink{
			ID:        3,
			Username:  "user2",
			DrinkData: drinkData,
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}

func TestGetDrinksByUser(t *testing.T) {
	app, cleanup := setupDbAndApp(t)
	defer cleanup()

	first := DrinkData{
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
	second := DrinkData{
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
	third := DrinkData{
		Name:           "Secret Drink",
		PrimaryAlcohol: "Scotch",
		PreferredGlass: "Rocks",
		Ingredients: []string{
			"2 oz scotch",
		},
		Publicity: DrinkPublicityPrivate,
	}

	drinks := []DrinkData{first, second, third}
	for _, d := range drinks {
		_, _ = t_createDrink_ok(t, app, CreateDrinkRequest{DrinkData: d}, commontest.AuthOpts{Username: to.StringPtr("user1")})
	}

	// Fetching as user1 should result in all drinks
	_, getResp := t_getDrinksByUser_ok(t, app, "user1", commontest.AuthOpts{Username: to.StringPtr("user1")})
	expectedGetResp := GetDrinksByUserResponse{
		Success: true,
		Drinks: []Drink{
			{
				ID:        1,
				Username:  "user1",
				DrinkData: first,
			},
			{
				ID:        2,
				Username:  "user1",
				DrinkData: second,
			},
			{
				ID:        3,
				Username:  "user1",
				DrinkData: third,
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
				DrinkData: first,
			},
			{
				ID:        2,
				Username:  "user1",
				DrinkData: second,
			},
		},
	}
	require.Equal(t, expectedGetResp, getResp)
}
