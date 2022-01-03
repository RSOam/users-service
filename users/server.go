package users

import (
	"context"
	"net/http"

	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHttpServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/users").Handler(ht.NewServer(
		endpoints.CreateUser,
		decodeCreateUserRequest,
		encodeResponse,
	))
	r.Methods("POST").Path("/users/login").Handler(ht.NewServer(
		endpoints.UserLogin,
		decodeUserLoginRequest,
		encodeResponse,
	))
	r.Methods("PUT").Path("/users/{id}").Handler(ht.NewServer(
		endpoints.UpdateUser,
		decodeUpdateUserRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/users/{id}").Handler(ht.NewServer(
		endpoints.GetUser,
		decodeGetUserRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/users").Handler(ht.NewServer(
		endpoints.GetUsers,
		decodeGetUsersRequest,
		encodeResponse,
	))
	r.Methods("DELETE").Path("/users/{id}").Handler(ht.NewServer(
		endpoints.DeleteUser,
		decodeDeleteUserRequest,
		encodeResponse,
	))
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
