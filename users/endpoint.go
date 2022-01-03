package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateUser endpoint.Endpoint
	GetUser    endpoint.Endpoint
	GetUsers   endpoint.Endpoint
	UpdateUser endpoint.Endpoint
	DeleteUser endpoint.Endpoint
	UserLogin  endpoint.Endpoint
}

func MakeEndpoints(s UsersService) Endpoints {
	return Endpoints{
		CreateUser: makeCreateUserEndpoint(s),
		GetUser:    makeGetUserEndpoint(s),
		GetUsers:   makeGetUsersEndpoint(s),
		UpdateUser: makeUpdateUserEndpoint(s),
		DeleteUser: makeDeleteUserEndpoint(s),
		UserLogin:  makeUserLoginEndpoint(s),
	}
}

func makeCreateUserEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		status, err := s.CreateUser(ctx, req.Username, req.Pasword)
		return CreateUserResponse{Status: status}, err
	}
}
func makeGetUserEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		user, err := s.GetUser(ctx, req.Id)
		return GetUserResponse{
			Username:     user.Username,
			Created:      user.Created,
			Modified:     user.Modified,
			Ratings:      user.Ratings,
			Comments:     user.Comments,
			Reservations: user.Reservations,
		}, err
	}
}
func makeGetUsersEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetUsers(ctx)
		return GetUsersResponse{
			Users: users,
		}, err
	}
}

func makeDeleteUserEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteUserRequest)
		status, err := s.DeleteUser(ctx, req.Id)
		return DeleteUserResponse{
			Status: status,
		}, err
	}
}
func makeUpdateUserEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserRequest)
		status, err := s.UpdateUser(ctx, req.Id, req.Username, req.Password)
		return CreateUserResponse{Status: status}, err
	}
}
func makeUserLoginEndpoint(s UsersService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UserLoginRequest)
		token, status, err := s.UserLogin(ctx, req.Username, req.Pasword)
		return UserLoginResponse{
			Token:  token,
			Status: status,
		}, err
	}
}
