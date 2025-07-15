package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"sykell-backend/internal/service"
	"sykell-backend/pkg/logger"
)

func GetUsers(c *fiber.Ctx) error {
	users, err := service.FetchUsers()
	if err != nil {
		logger.Sugar().Errorf("FetchUsers error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch users",
		})
	}
	return c.JSON(fiber.Map{
		"data": users,
		"count": len(users),
	})
}

func GetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := service.GetUserByID(id)
	if err != nil {
		logger.Sugar().Errorf("GetUserByID error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": user,
	})
}

func CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	user, err := service.CreateUser(req)
	if err != nil {
		logger.Sugar().Errorf("CreateUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": user,
		"message": "User created successfully",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var req service.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}

	user, err := service.UpdateUser(id, req)
	if err != nil {
		logger.Sugar().Errorf("UpdateUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not update user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": user,
		"message": "User updated successfully",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = service.DeleteUser(id)
	if err != nil {
		logger.Sugar().Errorf("DeleteUser error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var creds service.Credentials
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid login payload",
		})
	}

	token, err := service.Authenticate(creds)
	if err != nil {
		logger.Sugar().Errorf("Authentication error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication failed",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"message": "Login successful",
	})
}