package routes

import (
	"whoareu/middlewares/cors"
	"whoareu/services/chat/getmessages"
	"whoareu/services/chat/ws/startchatting"
	"whoareu/services/testim"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const (
	api_url = "/api/v1"
)

func InitRoutes(r *gin.Engine, db *gorm.DB, mdb *mongo.Client) {
	cors.InitCors(r)

	r.GET(api_url+"/testim", func(ctx *gin.Context) {
		testim.GetTestim(ctx)
	})

	r.GET(api_url+"/chats/getmessages/:chat_id", func(ctx *gin.Context) {
		getmessages.GetMessages(ctx, mdb)
	})

	// websockets
	r.GET(api_url+"/ws/startchat/:chat_id/:user_id", func(ctx *gin.Context) {
		startchatting.ConnectChat(ctx, mdb)
	})
}
