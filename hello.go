package main

import (
	"apis"
	"databases"

	_ "github.com/joho/godotenv/autoload"
)

// var exisitingNum int = 45 // can this one be used in another file?

func main() {

	// Setting the environment variables here
	databases.InitDB() // the client can be extracted to here // init the database to create the client

	// newUser := models.User{"Chihcheng", "Hsieh", "mike820808@gmail.com"}

	// models.AddUser(newUser)

	// models.DeleteUser("5d15825b6bfd7042f777f7c9")

	// users := models.FindUsers(bson.M{"email": "mike820808@gmail.com"})

	// fmt.Println(users)

	// if len(users) > 2 {
	// 	fmt.Println("Found User")
	// }

	// fmt.Println("This length is: ", len(users))

	// usersJSON, err := json.Marshal(users)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(usersJSON))

	//
	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080

	router := apis.InitRouter()

	// Setup the auth middleware for home page

	// router.GET("/", func(c *gin.Context) {

	// 	user := models.User{
	// 		FirstName: "Chihcheng",
	// 		LastName:  "Hsieh",
	// 		Password:  "password",
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"user": user,
	// 	})
	// })

	router.Run()

}
