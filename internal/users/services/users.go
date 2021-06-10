package services

import (
	"errors"

	"github.com/jmoiron/sqlx"
	. "github.com/morzhanov/go-realworld/internal/users/dto"
	. "github.com/morzhanov/go-realworld/internal/users/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	db *sqlx.DB
}

func (s *UsersService) GetUserData(userId string) (user *GetUserDto, err error) {
	q := `SELECT id, username FROM users
		WHERE id = $1`
	row := s.db.QueryRow(q, userId)

	user = &GetUserDto{}
	err = row.Scan(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UsersService) ValidateUserPassword(username string, password string) error {
	q := `SELECT id, username FROM users
		WHERE username = $1`
	row := s.db.QueryRow(q, username)

	user := &User{}
	if err := row.Scan(user.ID, user.Username, user.Password); err != nil {
		return err
	}

	if checkPasswordHash(password, user.Password) {
		return nil
	}
	return errors.New("Wrong password")
}

func (s *UsersService) CreateUser(data *CreateUserDto) (res *GetUserDto, err error) {
	hashedPassword, err := hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, useraname`
	row := s.db.QueryRow(q, data.Username, hashedPassword)

	res = &GetUserDto{}
	err = row.Scan(res.ID, res.Username)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *UsersService) DeleteUser(userId string) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Query(q, userId)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUsersService(db *sqlx.DB) *UsersService {
	return &UsersService{db}
}
