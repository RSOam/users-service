package users

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type database struct {
	db     *mongo.Database
	logger log.Logger
	consul consulapi.Client
}

func NewDatabase(db *mongo.Database, logger log.Logger, consul consulapi.Client) UserDB {
	return &database{
		db:     db,
		logger: log.With(logger, "database", "mongoDB"),
		consul: consul,
	}
}

func (dat *database) CreateUser(ctx context.Context, username string, password string) error {
	user := User{
		Username:     username,
		Password:     password,
		Ratings:      []Rating{},
		Comments:     []Comment{},
		Reservations: []Reservation{},
		Created:      time.Now().Format(time.RFC3339),
		Modified:     time.Now().Format(time.RFC3339),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := dat.db.Collection("Users").InsertOne(ctx, user)
	if err != nil {
		dat.logger.Log("Error inserting user into DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetUser(ctx context.Context, id string) (User, error) {
	tempUser := User{}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		dat.logger.Log("Error getting user from DB: ", err.Error())
		return tempUser, err
	}
	val, _ := getConsulValue(dat.consul, dat.logger, "commratService")
	val2, _ := getConsulValue(dat.consul, dat.logger, "reservationsService")
	ratings, err := getUserRatings(val, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting user from DB: ", err.Error())
		return tempUser, err
	}
	comments, err := getUserComments(val, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting user from DB: ", err.Error())
		return tempUser, err
	}
	reservations, err := getUserReservations(val2, dat.logger, id)
	if err != nil {
		dat.logger.Log("Error getting user from DB: ", err.Error())
		return tempUser, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = dat.db.Collection("Users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempUser)
	if err != nil {
		dat.logger.Log("Error getting user from DB: ", err.Error())
		return tempUser, err
	}
	tempUser.Ratings = ratings
	tempUser.Comments = comments
	tempUser.Reservations = reservations
	return tempUser, nil
}
func (dat *database) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error deleting user from DB: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := dat.db.Collection("Users").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		dat.logger.Log("Error deleting user from DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetUsers(ctx context.Context) ([]User, error) {
	tempUser := User{}
	tempUsers := []User{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := dat.db.Collection("Users").Find(ctx, bson.D{})
	if err != nil {
		dat.logger.Log("Error getting users from DB: ", err.Error())
		return tempUsers, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempUser)
		if err != nil {
			dat.logger.Log("Error getting users from DB: ", err.Error())
			return tempUsers, err
		}
		tempUsers = append(tempUsers, tempUser)
	}
	return tempUsers, nil
}
func (dat *database) UpdateUser(ctx context.Context, id string, username string, password string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error updating user: ", err.Error())
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"username": username,
			"passowrd": password,
			"modified": time.Now().Format(time.RFC3339),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = dat.db.Collection("Users").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		dat.logger.Log("Error updating user: ", err.Error())
		return err
	}

	return nil
}
func (dat *database) UserLogin(ctx context.Context, username string, password string) (string, error) {
	tempUser := User{
		Username: username,
		Password: password,
	}
	dbUser := User{}
	token := ""
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": tempUser.Username}
	err := dat.db.Collection("Users").FindOne(ctx, filter).Decode(&dbUser)
	if err != nil {
		dat.logger.Log("Error getting user: ", err.Error())
		return "", err
	}

	if verifyPassword(tempUser.Password, dbUser.Password) {
		val, _ := getConsulValue(dat.consul, dat.logger, "jwtSecret")
		token, err = ToJWT(dbUser, val)
		if err != nil {
			dat.logger.Log("Error creating JWT token: ", err.Error())
			return "", err
		}
	} else {
		dat.logger.Log("Problem verifying user: ", errors.New("credentials missmatch"))
		return "", err
	}
	return token, nil
}
func getUserRatings(commratAddr string, logger log.Logger, userID string) ([]Rating, error) {
	requestBody, err := json.Marshal(GetUserRatingsRequest{})
	tempResponse := GetUserRatingsResponse{}
	tempRatings := []Rating{}
	if err != nil {
		return tempRatings, err
	}
	client := &http.Client{}
	commratUri := commratAddr + "/ratings/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?user="+userID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempRatings, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempRatings, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempRatings = tempResponse.Ratings
	if err != nil {
		return tempRatings, err
	}
	client.CloseIdleConnections()
	return tempRatings, nil
}
func getUserComments(commratAddr string, logger log.Logger, userID string) ([]Comment, error) {
	requestBody, err := json.Marshal(GetUserCommentsRequest{})
	tempResponse := GetUserCommentsResponse{}
	tempComments := []Comment{}
	if err != nil {
		return tempComments, err
	}
	client := &http.Client{}
	commratUri := commratAddr + "/comments/"
	req, err := http.NewRequest(http.MethodGet, commratUri+"?user="+userID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempComments, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempComments, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempComments = tempResponse.Comments
	if err != nil {
		return tempComments, err
	}
	client.CloseIdleConnections()
	return tempComments, nil
}
func getUserReservations(reservationsAddr string, logger log.Logger, userID string) ([]Reservation, error) {
	requestBody, err := json.Marshal(GetUserReservationsRequest{})
	tempResponse := GetUserReservationsResponse{}
	tempReservations := []Reservation{}
	if err != nil {
		return tempReservations, err
	}
	client := &http.Client{}
	reservationsUri := reservationsAddr + "/reservations/"
	req, err := http.NewRequest(http.MethodGet, reservationsUri+"?user="+userID, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempReservations, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempReservations, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tempResponse)

	tempReservations = tempResponse.Reservations
	if err != nil {
		return tempReservations, err
	}
	client.CloseIdleConnections()
	return tempReservations, nil
}
func ToJWT(u User, secret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = u.ID
	claims["user_username"] = u.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

//FromJWT returns a User id from a JWT token
func FromJWT(token string, secret string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claims["user_id"].(string), nil
}
func getConsulValue(consul consulapi.Client, logger log.Logger, key string) (string, error) {
	kv := consul.KV()
	keyPair, _, err := kv.Get(key, nil)
	if err != nil {
		logger.Log("msg", "Failed getting consul key")
		return "", err
	}
	return string(keyPair.Value), nil
}
