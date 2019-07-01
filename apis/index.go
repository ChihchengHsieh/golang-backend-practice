package apis

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	InitWebSocketApi(router)
	UserApiInit(router)
	PostApiInit(router)
	CommentApiInit(router)
	return router
}
