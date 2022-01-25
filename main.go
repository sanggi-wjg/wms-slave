package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wms_slave/route"
	"wms_slave/server"
	"wms_slave/server/database"
	"wms_slave/server/logger"
)

func init() {
	server.Init()
	database.Init()
	logger.Init(server.Config.Debug)
	if server.Config.Debug {
		gin.ForceConsoleColor()
	}
	gin.SetMode(server.Config.RunMode)
}

func main() {
	router := route.Init()

	app := &http.Server{
		Addr:         server.Config.HttpPort,
		Handler:      router,
		WriteTimeout: server.Config.WriteTimeout,
		ReadTimeout:  server.Config.ReadTimeout,
	}
	if err := app.ListenAndServe(); err != nil {
		panic(err)
	}
	//err := manners.ListenAndServe(server.Config.HttpPort, router)
	//if err != nil {
	//	return
	//}
	//manners.Close()
}
