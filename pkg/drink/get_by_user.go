package drink

import (
	"database/sql"
	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type GetDrinksByUserResponse struct {
	Error   string  `json:"error"`
	Success bool    `json:"success"`
	Drinks  []Drink `json:"drinks"`
}

func getDrinksByUser(db *sql.DB) auth.FiberClaimsHandler {
	return func(c *fiber.Ctx, claims auth.Claims) error {
		username := c.Params("username")

		var modelList []Model
		var err error
		if username == claims.Username {
			modelList, err = getAllDrinksByUser(username, db)
		} else {
			modelList, err = getAllPublicDrinksByUser(username, db)
		}

		if err != nil {
			return err
		}

		drinks := make([]Drink, 0, len(modelList))
		for _, m := range modelList {
			d, err := fromDb(m)
			if err != nil {
				return err
			}
			drinks = append(drinks, d)
		}

		return c.JSON(GetDrinksByUserResponse{
			Success: true,
			Drinks: drinks,
		})
	}

}
