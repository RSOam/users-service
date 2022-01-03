package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Charger struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Location      Location           `json:"location"`
	AverageRating float64            `json:"averageRating"`
	Ratings       []Rating           `json:"ratings"`
	Comments      []Comment          `json:"comments"`
	Reservations  []Reservation      `json:"reservations"`
	Created       string             `json:"created"`
	Modified      string             `json:"modified"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Text      string             `json:"text"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}
type Rating struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Rating    int                `json:"rating"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type Reservation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	From      string             `json:"from"`
	To        string             `json:"to"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username"`
	Password     string             `json:"password,omitempty"`
	Created      string             `json:"created"`
	Modified     string             `json:"modified"`
	Ratings      []Rating           `json:"ratings"`
	Comments     []Comment          `json:"comments"`
	Reservations []Reservation      `json:"reservations"`
}

type UserDB interface {
	CreateUser(ctx context.Context, username string, password string) error
	GetUser(ctx context.Context, id string) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id string, username string, password string) error
	DeleteUser(ctx context.Context, id string) error
	UserLogin(ctx context.Context, username string, password string) (string, error)
}
