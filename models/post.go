package models

import (
	"context"
	"databases"
	"log"
	"time"
	"utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title     string
	Author    primitive.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string
}

// CRUD here
func AddPost(inputPost interface{}) (interface{}, error) {

	result, err := databases.DB.Collection("post").InsertOne(context.TODO(), inputPost)

	utils.ErrorChecking(err) // Maybe we can remove it later to boost the speed of server side

	return result.InsertedID, err // Check the error in the router so we can decide how to handle the respnse json there
}

func AddPosts(inputPosts []interface{}) (interface{}, error) {

	// Should I combine this one and the one above?
	result, err := databases.DB.Collection("post").InsertMany(context.TODO(), inputPosts)

	utils.ErrorChecking(err)

	return result.InsertedIDs, err

}

func DeletePosByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	_, err = databases.DB.Collection("post").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return err
}

func UpdatePostById(id string, upadateDetail bson.M) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	result, err := databases.DB.Collection("post").UpdateOne(context.TODO(), bson.M{"_id": oid}, upadateDetail)

	return result.UpsertedID, err
}

func UpdatePosts(filterDetail bson.M, updateDetail bson.M) {
	result, err := databases.DB.Collection("post").UpdateMany(context.TODO(), filterDetail, updateDetail)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}
}

func FindPostByID(id string) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	var post Post
	err = databases.DB.Collection("post").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&post)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func FindPosts(filterDetail bson.M) (interface{}, error) {
	var posts []*Post
	result, err := databases.DB.Collection("post").Find(context.TODO(), filterDetail)
	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Post
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &elem)
	}

	return posts, nil

}

//MGDB_APIKEY=mongodb+srv://mike820808:0933a887632@cluster0-a0dck.mongodb.net/test?retryWrites=true&w=majority
func FindPostsByUserID(id string) (interface{}, error) {

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	posts, err := FindPosts(bson.M{"author": oid})

	if err != nil {
		return nil, err
	}

	return posts, nil
}
