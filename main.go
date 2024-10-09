package main

import (
	"file_store/lib"
	"file_store/model/mysql"
	"file_store/router"
	"log"
)

func main() {
	serverConfig := lib.LoadServerConfig()
	mysql.InitDB(serverConfig)
	defer mysql.DB.Close()

	r := router.SetupRoute()
	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")

	if err := r.Run(":9090"); err != nil {
		log.Fatal("failed to start server ...")
	}
}
