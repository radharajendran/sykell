package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"sykell-backend/internal/repository"
	"sykell-backend/pkg/logger"
)

var userRepo = repository.NewUserRepository()

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FetchUsers() ([]User, error) {
	repoUsers, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convert repository users to service users
	users := make([]User, len(repoUsers))
	for i, repoUser := range repoUsers {
		users[i] = User{
			ID:        repoUser.ID,
			Name:      repoUser.Name,
			Email:     repoUser.Email,
			CreatedAt: repoUser.CreatedAt,
			UpdatedAt: repoUser.UpdatedAt,
		}
	}

	return users, nil
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(req CreateUserRequest) (User, error) {
	if req.Name == "" {
		return User{}, errors.New("name is required")
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Sugar().Errorf("Failed to hash password: %v", err)
			return User{}, err
		}
		hashedPassword = string(hashed)
	}

	repoUser, err := userRepo.Create(req.Name, req.Email, hashedPassword)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        repoUser.ID,
		Name:      repoUser.Name,
		Email:     repoUser.Email,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

type Credentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Authenticate(creds Credentials) (string, error) {
	if creds.Name == "" {
		return "", errors.New("name is required")
	}

	// Get user from database
	user, err := userRepo.GetByName(creds.Name)
	if err != nil {
		logger.Sugar().Errorf("Failed to get user: %v", err)
		return "", errors.New("authentication failed")
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	// Verify password if provided
	if creds.Password != "" && user.Password != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
		if err != nil {
			return "", errors.New("invalid credentials")
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := getJWTSecret()
	signed, err := token.SignedString(secret)
	if err != nil {
		logger.Sugar().Errorf("JWT signing error: %v", err)
		return "", err
	}

	return signed, nil
}

func getJWTSecret() []byte {
	secret := "your-secret-key" // Default secret
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		secret = envSecret
	}
	return []byte(secret)
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetUserByID(id int) (*User, error) {
	repoUser, err := userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoUser == nil {
		return nil, nil
	}

	return &User{
		ID:        repoUser.ID,
		Name:      repoUser.Name,
		Email:     repoUser.Email,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

func UpdateUser(id int, req UpdateUserRequest) (*User, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	repoUser, err := userRepo.Update(id, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	if repoUser == nil {
		return nil, nil
	}

	return &User{
		ID:        repoUser.ID,
		Name:      repoUser.Name,
		Email:     repoUser.Email,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}, nil
}

func DeleteUser(id int) error {
	return userRepo.Delete(id)
}