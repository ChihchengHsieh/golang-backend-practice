package apis

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { // what's this function for?
		return true
	},
}

// InitWebSocketApi - init the whole websocket api here
func InitWebSocketApi(router *gin.Engine) {

	router.GET("/ws", func(c *gin.Context) {
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		defer ws.Close()
		if err != nil {
			log.Printf("The websocket is not working due to the error : %+v \n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error occur: %+v\b", err)
				break
			}

			if string(message) == "ping" {
				message = []byte("pong")
			}

			err = ws.WriteMessage(mt, message)

			if err != nil {
				break
			}
		}

	})
}
