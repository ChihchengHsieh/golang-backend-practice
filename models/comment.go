package models

import (
	"context"
	"databases"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id",bson:"_id,omitempty"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    primitive.ObjectID
	Post      primitive.ObjectID
}

// CRUD here
func AddComment(inputComment interface{}) (interface{}, error) {

	result, err := databases.DB.Collection("comment").InsertOne(context.TODO(), inputComment)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil

}

func AddComments(inputComments []interface{}) {

	// Should I combine this one and the one above?
	result, err := databases.DB.Collection("comment").InsertMany(context.TODO(), inputComments)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}

}

func DeleteCommentByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
	}

	_, err = databases.DB.Collection("comment").DeleteOne(context.TODO(), bson.M{"_id": oid})

	if err != nil {
		return err
	}

	return nil
}

func UpdateCommentById(id string, upadateDetail bson.M) (interface{}, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}
	result, err := databases.DB.Collection("comment").UpdateOne(context.TODO(), bson.M{"_id": oid}, upadateDetail)

	if err != nil {
		return nil, err
	}

	return result.UpsertedID, nil
}

func UpdateComments(filterDetail bson.M, updateDetail bson.M) {
	result, err := databases.DB.Collection("comment").UpdateMany(context.TODO(), filterDetail, updateDetail)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(result)
	}
}

func FindCommentByID(id string) Comment {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Fatal(err)
	}
	var comment Comment
	err = databases.DB.Collection("comment").FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&comment)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(comment)
	}

	return comment
}

func FindComments(filterDetail bson.M) ([]*Comment, error) {
	var comments []*Comment
	result, err := databases.DB.Collection("comment").Find(context.TODO(), filterDetail)
	defer result.Close(context.TODO())

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Comment
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &elem)
	}

	return comments, nil

}

// FindCommentByUserID - Find all the0 comments for a certain user
func FindCommentByUserID(id string) ([]*Comment, error) {

	var comments []*Comment

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	result, err := databases.DB.Collection("comment").Find(context.TODO(), bson.M{"author": oid})

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Comment
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &elem)

	}

	return comments, nil

}

// FindCommentByPostID - Find all the comments for a certain post
func FindCommentByPostID(id string) ([]*Comment, error) {
	var comments []*Comment

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	result, err := databases.DB.Collection("comment").Find(context.TODO(), bson.M{"post": oid})

	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var elem Comment
		err := result.Decode(&elem)
		if err != nil {
			return nil, err
		}

		comments = append(comments, &elem)

	}

	return comments, nil

}
