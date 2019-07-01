package apis

import (
	"models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

func CommentApiInit(router *gin.Engine) {
	commentRouter := router.Group("/comment")
	{
		// Create comment for a certain post
		commentRouter.POST("/", func(c *gin.Context) {

			/*
				Comment Fields:
					Content   string
					CreatedAt time.Time
					UpdatedAt time.Time
					Author    primitive.ObjectID
					Post      primitive.ObjectID
			*/
			content, authorID, postID := c.PostForm("content"), c.PostForm("author"), c.PostForm("post")
			if content == "" || authorID == "" || postID == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "You have to provide content, postID and authorID",
				})
			}

			authorOID, err := primitive.ObjectIDFromHex(authorID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})

				return
			}
			postOID, err := primitive.ObjectIDFromHex(postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			commentID, err := models.AddComment(models.Comment{
				Content:   content,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Author:    authorOID,
				Post:      postOID,
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot Create the comment",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"commentID": commentID,
			})
		})

		// Update Comment
		commentRouter.PUT("/:id", func(c *gin.Context) {
			content := c.PostForm("content")

			upsertID, err := models.UpdateCommentById(c.Param("id"), bson.M{
				"$set": bson.M{"content": content, "updatedAt": time.Now()},
			})

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot update the comment",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"postID":   c.Param("id"),
				"upsertID": upsertID,
			})

		})

		// Fetch Comment for a user
		commentRouter.GET("/user/:id", func(c *gin.Context) {

			comments, err := models.FindCommentByUserID(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot Get the comments",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"comments": comments,
			})

		})

		// Delet Comment through ID
		commentRouter.DELETE("/:id", func(c *gin.Context) {
			err := models.DeleteCommentByID(c.Param("id"))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot Delete the post",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"postID": c.Param("id"),
			})
		})

		// Fetch Comment for a post

		commentRouter.GET("/post/:id", func(c *gin.Context) {
			comments, err := models.FindCommentByPostID(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
					"msg":   "Cannot get the posts",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"comments": comments,
			})
		})
	}
}
