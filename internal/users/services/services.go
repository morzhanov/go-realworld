package services

import (
	"errors"
	"github.com/morzhanov/go-realworld/internal/common/helper"

	"github.com/jmoiron/sqlx"
	urpc "github.com/morzhanov/go-realworld/api/grpc/users"
	. "github.com/morzhanov/go-realworld/internal/users/models"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	db *sqlx.DB
}

type UsersService interface {
	GetUserData(userId string) (user *urpc.UserMessage, err error)
	GetUserDataByUsername(username string) (user *urpc.UserMessage, err error)
	ValidateUserPassword(data *urpc.ValidateUserPasswordRequest) error
	CreateUser(data *urpc.CreateUserRequest) (res *urpc.UserMessage, err error)
	DeleteUser(userId string) error
}

func (s *usersService) GetUserData(userId string) (user *urpc.UserMessage, err error) {
	q := `SELECT id, username FROM users
		WHERE id = $1`
	row := s.db.QueryRow(q, userId)

	user = &urpc.UserMessage{}
	err = row.Scan(&user.Id, &user.Username)
	if err != nil {
		if helper.CheckNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *usersService) GetUserDataByUsername(username string) (user *urpc.UserMessage, err error) {
	q := `SELECT id, username FROM users
		WHERE username = $1`
	row := s.db.QueryRow(q, username)

	user = &urpc.UserMessage{}
	err = row.Scan(&user.Id, &user.Username)
	if err != nil {
		if helper.CheckNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (s *usersService) ValidateUserPassword(data *urpc.ValidateUserPasswordRequest) error {
	q := `SELECT id, username, password FROM users
		WHERE username = $1`
	row := s.db.QueryRow(q, data.Username)

	user := &User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New("user not found")
		}
		return err
	}

	if s.checkPasswordHash(data.Password, user.Password) {
		return nil
	}
	return errors.New("wrong password")
}

func (s *usersService) CreateUser(data *urpc.CreateUserRequest) (res *urpc.UserMessage, err error) {
	hashedPassword, err := s.hashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, username`
	row := s.db.QueryRow(q, data.Username, hashedPassword)

	res = &urpc.UserMessage{}
	err = row.Scan(&res.Id, &res.Username)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *usersService) DeleteUser(userId string) error {
	q := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Query(q, userId)
	return err
}

func (s *usersService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *usersService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUsersService(db *sqlx.DB) UsersService {
	return &usersService{db}
}
