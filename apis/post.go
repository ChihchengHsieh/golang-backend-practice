package apis

import (
	"models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

func PostApiInit(router *gin.Engine) {

	postRouter := router.Group("/post")
	{

		postRouter.GET("/", func(c *gin.Context) {
			posts, err := models.FindPosts(bson.M{})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot retrieve posts",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"posts": posts,
			})
		})

		postRouter.POST("/", func(c *gin.Context) {

			/*
				Post Fields:
					type Post struct {
						ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
						Title     string
						Author    primitive.ObjectID
						CreatedAt time.Time
						UpdatedAt time.Time
						Content   string
					}
			*/

			// Retrieve the data we need

			title, authorId, content := c.PostForm("title"), c.PostForm("author"), c.PostForm("content")

			oid, err := primitive.ObjectIDFromHex(authorId)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Can not get the author properly",
				})

				return
			}

			postID, err := models.AddPost(models.Post{
				Author:    oid,
				Title:     title,
				Content:   content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg":   "fail to add this post",
					"error": err,
				})

				return
			}

			c.JSON(http.StatusOK, gin.H{
				"postID": postID,
			})

		})

		// GET all post

		// Root for changing the post
		postRouter.PUT("/:id", func(c *gin.Context) {

			// Only the title and the content can be changed
			postID, title, content := c.Param("id"), c.PostForm("title"), c.PostForm("content")

			resultID, err := models.UpdatePostById(postID, bson.M{"$set": bson.M{"title": title, "content": content, "updatedAt": time.Now()}})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot update the Post",
				})

				return
			}

			c.JSON(http.StatusOK, gin.H{
				"updateID": resultID,
				"PostID":   postID,
			})

		})

		postRouter.DELETE("/:id", func(c *gin.Context) {
			postID := c.Param("id")
			err := models.DeletePosByID(postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot Delete the post",
				})
			}
		})

		// Get post for a specific user
		postRouter.GET("/user/:id", func(c *gin.Context) {
			posts, err := models.FindPostsByUserID(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get the posts for the user",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"posts": posts,
			})
		})

		// Get post through id
		postRouter.GET("/id/:id", func(c *gin.Context) {
			postID := c.Param("id")
			post, err := models.FindPostByID(postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get the post",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"post": post,
			})

		})

	}

}
