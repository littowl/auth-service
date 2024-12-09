package main

import (
	"auth-service/internal/domain"
	db "auth-service/internal/repository"
	"auth-service/transport/transport_http"
	"auth-service/transport/transport_http/handlers"
	"auth-service/utils"
)

func main() {
	logger := utils.CreateLogger()
	defer logger.Sync()

	config := utils.NewConfig()
	pool := db.DbStart(config, logger)
	if pool == nil {
		logger.Fatal("Failed to connect to database")
		return
	}
	database := db.NewDB(pool, logger)
	if database == nil {
		logger.Fatal("Failed to start database")
		return
	}

	domainService := domain.NewDomainService(database)
	httpHandlers := handlers.NewBaseHandler(logger, domainService)
	httpServer := transport_http.NewHttpServer(logger, config.HTTPAddr)

	httpServer.StartHTTPServer(httpHandlers)
	logger.Error("HTTP server down")
}
