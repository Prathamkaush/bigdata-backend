package main

import (
	"log"

	"bigdata-api/internal/api/routes"
	"bigdata-api/internal/config"
	"bigdata-api/internal/database"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Init databases
	database.ConnectClickHouse(&cfg)
	database.ConnectPostgres(&cfg)
	database.ConnectRedis(&cfg)

	// Start API
	app := routes.InitRoutes(&cfg)

	log.Println("🚀 Server running on port", cfg.ServerPort)
	app.Listen(":" + cfg.ServerPort)
}
