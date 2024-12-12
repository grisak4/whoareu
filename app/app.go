package app

import (
	"fmt"
	"whoareu/config"
	"whoareu/config/confget/servconf"
	mongoDB "whoareu/databases/mongodb"
	mainDB "whoareu/databases/postgresql"
	"whoareu/routes"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	config.InitConfig()

	mainDB.InitDatabase()
	defer mainDB.CloseDB()

	mongoDB.InitMongoDB()
	defer mongoDB.CloseMongoDB()

	routes.InitRoutes(r, mainDB.GetDB(), mongoDB.GetMongoDB())

	serv := servconf.GetServConf()
	r.Run(fmt.Sprintf("%s:%d",
		serv.Host, serv.Port))
}
