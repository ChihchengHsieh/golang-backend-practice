package apis

import (
	"encoding/json"
	"fmt"
	"middlewares"
	"models"
	"net/http"
	"utils"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UserApiInit(router *gin.Engine) {

	userRouter := router.Group("/user")

	{

		// Register Route
		userRouter.POST("/register", func(c *gin.Context) {

			/*
				User Fields:
					Email:
					Password: (Hashed)
					FirstName:
					LastName:
			*/

			// Check if the email already exist

			if c.PostForm("email") == "" || c.PostForm("password") == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Must Provide Email and Password",
				})
				return
			}

			if exist := len(models.FindUsers(bson.M{"email": c.PostForm("email")})); exist > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "This Eamil is used",
				})
				return
			}

			registerUser := models.User{
				Email:     c.PostForm("email"),
				FirstName: c.PostForm("firstName"),
				LastName:  c.PostForm("lastName"),
				Password:  c.PostForm("password"),
			}

			insertID, err := models.AddUser(&registerUser)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err, "message": "Fail to register the user"})
				return
			}

			registerUser.ID = insertID.(primitive.ObjectID)

			c.JSON(http.StatusOK, gin.H{"user": utils.StructToMap(registerUser), "token": utils.GenerateAuthToken(c.PostForm("email"))})

		})

		// Login Route

		userRouter.POST("/login", func(c *gin.Context) {
			/*
				Require Field
					Email
					Password (Hashed)
			*/

			inputEmail := c.PostForm("email")
			inputPassword := c.PostForm("password") // Expect Hased Password

			if inputEmail == "" || inputPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Must Provide Email and Password",
				})
				return
			}

			// How can we check if the user exist or not
			user, err := models.FindUserByEmailForLogin(inputEmail)

			if err == nil || inputPassword != user.Password {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":         "The UserName or Password isn't correct",
					"err":           err,
					"inputPassword": inputPassword,
					"user.Password": user.Password,
					"isCorrect":     inputPassword == user.Password,
				})

			} else {

				c.JSON(http.StatusOK, gin.H{
					"token": utils.GenerateAuthToken(user.Email),
					"user":  utils.StructToMap(user),
				})
			}

		})

		// Get User through id
		userRouter.GET("/id/:id", func(c *gin.Context) {
			id := c.Param("id")
			user, err := models.FindUserByID(id)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User Not Found",
				})
				return
			}

			// userMap := utils.StructToMap(user)

			// RemovePasswordField(&userMap)

			c.JSON(http.StatusOK, gin.H{
				"user": utils.StructToMap(user),
			})
		})

		// Get User through Email

		userRouter.GET("/email/:email", func(c *gin.Context) {
			email := c.Param("email")
			user, err := models.FindUserByEmail(email)

			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User Not Found",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user": utils.StructToMap(user),
			})
		})

		// Get Users through an array of emails
		userRouter.POST("/emails", func(c *gin.Context) {
			/*
				Make Fask Post to test it
			*/
			// Get array(string) from PostForm
			emailsArrayStr := c.PostForm("emails")
			// Transfer it into array type
			var emailsArray []string

			json.Unmarshal([]byte(emailsArrayStr), &emailsArray)

			fmt.Printf("Input emails:\n %s\n", emailsArrayStr)

			fmt.Printf("The emailsArray:\n %+v\n", emailsArray)

			fmt.Printf("The array len is %d\n", len(emailsArray))

			// Find all users
			users := models.FindUsers(bson.M{"email": bson.M{"$in": emailsArray}})

			// // retrun users through json
			// usersJSON, err := json.Marshal(users)
			// utils.ErrorChecking(err)

			// var userMap interface{}

			// err = json.Unmarshal(usersJSON, &userMap)

			// utils.ErrorChecking(err)

			c.JSON(http.StatusOK, gin.H{"users": users})
		})

		// Get User through an array of IDs

		userRouter.POST("/ids", middlewares.LoginAuth(), func(c *gin.Context) {

			idsArrayStr := c.PostForm("ids")

			var idStrArray []string
			var idArray []primitive.ObjectID

			err := json.Unmarshal([]byte(idsArrayStr), &idStrArray)
			utils.ErrorChecking(err)

			for _, idStr := range idStrArray {
				oid, err := primitive.ObjectIDFromHex(idStr)
				utils.ErrorChecking(err)
				idArray = append(idArray, oid)
			}

			users := models.FindUsers(bson.M{"_id": bson.M{"$in": idArray}})

			// usersJSON, err := json.Marshal(users)
			// utils.ErrorChecking(err)

			// ["5d16ca2e6312db5e989a0339", "5d16ca43c5c78322dee1ed9c"]

			// var userMap interface{}

			// err = json.Unmarshal(usersJSON, &userMap)

			// utils.ErrorChecking(err)

			c.JSON(http.StatusOK, gin.H{"user": users})

		})

		// Get all Users

	}
}

func RemovePasswordField(inputMap *map[string]interface{}) {
	delete(*inputMap, "Password")
}

func RemovePasswordFieldForArry(inputMapArray *[]*map[string]interface{}) {
	for _, m := range *inputMapArray {
		delete(*m, "Password")
	}
}
