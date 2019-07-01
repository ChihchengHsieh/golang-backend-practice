package databases

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var AnotherNumner int

// For getting the client
func InitDB() {
	fmt.Println(os.Getenv("MGDB_APIKEY"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MGDB_APIKEY"))
	AnotherNumner = 5

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	DB = client.Database("api")

}
