package repository

import (
	"database/sql"
	"time"

	"sykell-backend/pkg/database"
	"sykell-backend/pkg/logger"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) GetAll() ([]User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users ORDER BY id"
	
	rows, err := r.db.Query(query)
	if err != nil {
		logger.Sugar().Errorf("Failed to fetch users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var email sql.NullString
		
		err := rows.Scan(&user.ID, &user.Name, &email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logger.Sugar().Errorf("Failed to scan user: %v", err)
			return nil, err
		}
		
		if email.Valid {
			user.Email = email.String
		}
		
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Sugar().Errorf("Error iterating over users: %v", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Create(name, email, password string) (*User, error) {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?)"
	
	result, err := r.db.Exec(query, name, email, password)
	if err != nil {
		logger.Sugar().Errorf("Failed to create user: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Sugar().Errorf("Failed to get last insert ID: %v", err)
		return nil, err
	}

	return r.GetByID(int(id))
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?"
	
	var user User
	var email sql.NullString
	
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Sugar().Errorf("Failed to get user by ID: %v", err)
		return nil, err
	}

	if email.Valid {
		user.Email = email.String
	}

	return &user, nil
}

func (r *UserRepository) GetByName(name string) (*User, error) {
	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE name = ?"
	
	var user User
	var email sql.NullString
	
	err := r.db.QueryRow(query, name).Scan(&user.ID, &user.Name, &email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Sugar().Errorf("Failed to get user by name: %v", err)
		return nil, err
	}

	if email.Valid {
		user.Email = email.String
	}

	return &user, nil
}

func (r *UserRepository) Update(id int, name, email string) (*User, error) {
	query := "UPDATE users SET name = ?, email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	
	_, err := r.db.Exec(query, name, email, id)
	if err != nil {
		logger.Sugar().Errorf("Failed to update user: %v", err)
		return nil, err
	}

	return r.GetByID(id)
}

func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		logger.Sugar().Errorf("Failed to delete user: %v", err)
		return err
	}

	return nil
}
