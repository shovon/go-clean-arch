package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DatabaseUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UniquenessViolationError struct {
	emailAddress string
}

var _ error = UniquenessViolationError{}

func (u UniquenessViolationError) Error() string {
	return fmt.Sprintf("User with email %s already exists", u.emailAddress)
}

type UserNotFoundError struct{}

var _ error = UserNotFoundError{}

func (u UserNotFoundError) Error() string {
	return "User not found"
}

type UsersModel interface {
	GetUsers() ([]DatabaseUser, error)
	GetUser(id string) (DatabaseUser, error)
	CreateUser(user User) (DatabaseUser, error)
	DeleteUser(id string) error
}

func NewUsersModel(db *sql.DB) UsersModel {
	return &usersModelImpl{db: db}
}

type usersModelImpl struct {
	db *sql.DB
}

func (u *usersModelImpl) GetUsers() ([]DatabaseUser, error) {
	rows, err := u.db.Query("SELECT id, name, email_address FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []DatabaseUser{}

	for rows.Next() {
		var u DatabaseUser
		err := rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (u *usersModelImpl) GetUser(id string) (DatabaseUser, error) {
	var user DatabaseUser

	err := u.db.QueryRow("SELECT id, name, email_address FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return DatabaseUser{}, UserNotFoundError{}
		} else {
			return DatabaseUser{}, err
		}
	}

	return user, nil
}

func (u *usersModelImpl) CreateUser(user User) (DatabaseUser, error) {
	var id int64
	err := u.db.
		QueryRow(
			"INSERT INTO users (name, email_address) VALUES ($1, $2) RETURNING id",
			user.Name,
			user.Email,
		).
		Scan(&id)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok {
			if pgErr.Code.Name() == "unique_violation" {
				return DatabaseUser{}, UniquenessViolationError{emailAddress: user.Email}
			}
		}
		return DatabaseUser{}, err
	}
	var dbUser DatabaseUser
	u.db.QueryRow("SELECT id, name, email_address FROM users WHERE id = $1", id).Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email)
	return dbUser, nil
}

func (u *usersModelImpl) DeleteUser(id string) error {
	_, err := u.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
