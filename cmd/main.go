package main

import (
	"spotsync/internal/config"
	"spotsync/internal/server"
)

func main() {
	// load env
	cfg := config.LoadEnv()

	// connect database
	db := config.ConnectDatabase(cfg)

	// start server
	server.Start(db, cfg)
}
