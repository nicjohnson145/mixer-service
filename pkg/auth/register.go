package auth

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type RegisterNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterNewUserResponse struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

func registerNewUser(db *sql.DB) FiberClaimsHandler {
	return func(c *fiber.Ctx, claims Claims) error {
		var payload RegisterNewUserRequest
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(RegisterNewUserResponse{
				Error: err.Error(),
				Success: false,
			})
		}

		existingUser, err := getUserByName(payload.Username, db)
		if err != nil {
			return err
		}
		
		if existingUser != nil {
			return c.Status(fiber.StatusBadRequest).JSON(RegisterNewUserResponse{
				Error: fmt.Sprintf("user %v already exists", payload.Username),
				Success: false,
			})
		}

		hashedPw, err := hashPassword(payload.Password)
		if err != nil {
			return err
		}
		
		newUser := UserModel{
			Username: payload.Username,
			Password: hashedPw,
		}
		err = createUser(newUser, db)
		if err != nil {
			return err
		}

		return c.JSON(RegisterNewUserResponse{
			Success: true,
		})
	}

}
