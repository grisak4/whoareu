package routes

import (
	"whoareu/middlewares/authmiddleware"
	"whoareu/middlewares/cors"
	"whoareu/services/auth/login"
	"whoareu/services/auth/regin"
	"whoareu/services/chat/createchat"
	"whoareu/services/chat/getchats"
	"whoareu/services/chat/getmessages"
	"whoareu/services/chat/joinchat"
	"whoareu/services/chat/ws/startchatting"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const (
	api_url = "/api/v1"
)

func InitRoutes(r *gin.Engine, db *gorm.DB, mdb *mongo.Database) {
	cors.InitCors(r)

	r.POST(api_url+"/regin", func(ctx *gin.Context) {
		regin.Regin(ctx, db)
	})
	r.POST(api_url+"/login", func(ctx *gin.Context) {
		login.Login(ctx, db)
	})

	userRoutes := r.Group("/user")
	userRoutes.Use(authmiddleware.AuthMiddleware([]string{"user"}))
	{
		r.GET(api_url+"/chats/getmessages/:chat_id", func(ctx *gin.Context) { getmessages.GetMessages(ctx, mdb) })
		r.GET(api_url+"/chats/getchats/:user_id", func(ctx *gin.Context) { getchats.GetAllChats(ctx, mdb) })
		r.POST(api_url+"/chats/create-chat", func(ctx *gin.Context) { createchat.CreateChat(ctx, mdb) })
		r.POST(api_url+"/chats/joinchat/:chat_id/:user_id", func(ctx *gin.Context) { joinchat.JoinChat(ctx, mdb) })

		// websockets
		r.GET(api_url+"/ws/startchat/:chat_id/:user_id", func(ctx *gin.Context) { startchatting.ConnectChat(ctx, mdb) })
	}
}
