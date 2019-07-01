package models

import (
	"context"
	"databases"
	"log"
	"utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

var projectionForRemovingPassword = bson.D{
	{"password", 0},
}

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName string
	LastName  string
	Email     string
	Password  string
}

// CRUD here
func AddUser(inputUser *User) (interface{}, error) {

	result, err := databases.DB.Collection("user").InsertOne(context.TODO(), inputUser)

	return result.InsertedID, err
}

func AddUsers(inputUsers *[]interface{}) {

	// Should I combine this one and the one above?
	result, err := databases.DB.Collection("user").InsertMany(context.TODO(), *inputUsers)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}

}

func DeleteUserByID(id string) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Fatal(err)
	}

	result, err := databases.DB.Collection("user").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}
}

func UpdateUserById(id string, upadateDetail bson.M) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Fatal(err)
	}
	result, err := databases.DB.Collection("user").UpdateOne(context.TODO(), bson.M{"_id": oid}, upadateDetail)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}
}

func UpdateUsers(filterDetail bson.M, updateDetail bson.M) {
	result, err := databases.DB.Collection("user").UpdateMany(context.TODO(), filterDetail, updateDetail)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}
}

func FindUserByID(id string) (User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	utils.ErrorChecking(err)
	// projection := bson.D{
	// 	{"password", 0},
	// }

	var user User
	err = databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"_id": oid},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	return user, err
}

func FindUserByEmailForLogin(email string) (User, error) {
	var user User
	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return user, err
}

// FindUserByEmail - Finding the user through input Email
func FindUserByEmail(email string) (User, error) {
	var user User

	err := databases.DB.Collection("user").FindOne(context.TODO(), bson.M{"email": email},
		options.FindOne().SetProjection(projectionForRemovingPassword)).Decode(&user)

	return user, err
}

func FindUsers(filterDetail bson.M) []*User {
	var users []*User
	result, err := databases.DB.Collection("user").Find(context.TODO(), filterDetail,
		options.Find().SetProjection(projectionForRemovingPassword))

	if err != nil {
		log.Fatal(err)
	}
	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {
		var elem User
		err := result.Decode(&elem)
		utils.ErrorChecking(err)
		users = append(users, &elem)
	}

	return users

}
