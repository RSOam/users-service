package users

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	consulapi "github.com/hashicorp/consul/api"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	db     UserDB
	logger log.Logger
	consul consulapi.Client
}

func NewService(db UserDB, logger log.Logger, consul consulapi.Client) UsersService {
	return &service{
		db:     db,
		logger: logger,
		consul: consul,
	}
}

func (s service) CreateUser(ctx context.Context, username string, password string) (string, error) {
	logger := log.With(s.logger, "method: ", "CreateUser")
	pw, err := passwordHash(password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	if err := s.db.CreateUser(ctx, username, pw); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("create User", nil)
	return "Ok", nil
}
func (s service) GetUser(ctx context.Context, id string) (User, error) {
	logger := log.With(s.logger, "method", "GetUser")
	user, err := s.db.GetUser(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return user, err
	}
	logger.Log("Get User", id)
	return user, nil
}
func (s service) GetUsers(ctx context.Context) ([]User, error) {
	logger := log.With(s.logger, "method", "GetUsers")
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		return users, err
	}
	logger.Log("Get Users")
	return users, nil
}

func (s service) DeleteUser(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteUser")
	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("Delete User", id)
	return "Ok", nil
}
func (s service) UpdateUser(ctx context.Context, id string, username string, password string) (string, error) {
	logger := log.With(s.logger, "method: ", "UpdateUser")
	pw, err := passwordHash(password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	if err := s.db.UpdateUser(ctx, id, username, pw); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("update User", id)
	return "Ok", nil
}
func (s service) UserLogin(ctx context.Context, username string, password string) (string, string, error) {
	logger := log.With(s.logger, "method: ", "UserLogin")
	token := ""
	token, err := s.db.UserLogin(ctx, username, password)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", "", err
	}
	logger.Log("update User")
	return token, "Ok", nil
}

func passwordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
func verifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return (err == nil)
}
