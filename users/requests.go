package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	CreateUserRequest struct {
		Username string `json:"username"`
		Pasword  string `json:"password"`
	}
	CreateUserResponse struct {
		Status string `json:"status"`
	}
	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		Username     string        `json:"username"`
		Created      string        `json:"created"`
		Modified     string        `json:"modified"`
		Ratings      []Rating      `json:"ratings"`
		Comments     []Comment     `json:"comments"`
		Reservations []Reservation `json:"reservations"`
	}
	GetUsersRequest struct {
	}
	GetUsersResponse struct {
		Users []User `json:"users"`
	}
	UpdateUserRequest struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
	}
	UpdateUserResponse struct {
		Status string `json:"status"`
	}
	DeleteUserRequest struct {
		Id string `json:"id"`
	}
	DeleteUserResponse struct {
		Status string `json:"status"`
	}
	UserLoginRequest struct {
		Username string `json:"username"`
		Pasword  string `json:"password"`
	}
	UserLoginResponse struct {
		Token  string `token:"status"`
		Status string `json:"status"`
	}
	//OTHER
	GetUserRatingsRequest struct {
	}
	GetUserRatingsResponse struct {
		Ratings []Rating `json:"ratings"`
	}
	GetUserCommentsRequest struct {
	}
	GetUserCommentsResponse struct {
		Comments []Comment `json:"comments"`
	}
	GetUserReservationsRequest struct {
	}
	GetUserReservationsResponse struct {
		Reservations []Reservation `json:"reservations"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetUserRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetUserRequest{}
	return req, nil
}
func decodeDeleteUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := DeleteUserRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeUserLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := UserLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
