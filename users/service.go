package users

import (
	"context"
)

type UsersService interface {
	CreateUser(ctx context.Context, username string, password string) (string, error)
	GetUser(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id string, username string, password string) (string, error)
	DeleteUser(ctx context.Context, id string) (string, error)
	UserLogin(ctx context.Context, username string, password string) (string, string, error)
}
