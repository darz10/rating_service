package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"ratingBookingService/config"
	_ "ratingBookingService/docs"
	"ratingBookingService/pkg/handlers"
	"ratingBookingService/pkg/repository"
)

// @title           Swagger rating services
// @version         1.0

// @host      localhost:8001
// @BasePath  /

// @securityDefinitions.apiKey  ApiKeyAAuth
// @in header
// @name Authorization

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("cannot init config:", err)
	}
	dbConnection, err := repository.ConnectDB(
		repository.Config{
			Host:     viper.GetString("DBHOST"),
			Port:     viper.GetString("DBPORT"),
			Username: viper.GetString("DBUSER"),
			Password: viper.GetString("DBPASSWORD"),
			DBName:   viper.GetString("DB"),
		})
	if err != nil {
		log.Fatal("Database connection failed", err)
	}
	repository := repository.CreateNewRatingRepository(dbConnection)
	handler := handlers.CreateNewRatingHandler(repository)
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		api.POST("/ratings", handlers.AuthJWTMiddleware(), handler.CreateRating)
		api.PATCH("/ratings/:id", handlers.AuthJWTMiddleware(), handler.UpdateRating)
		api.DELETE("/ratings/:id", handlers.AuthJWTMiddleware(), handler.DeleteRating)
		api.GET("/ratings", handler.GetRatingsAll)
		api.GET("/ratings/:id", handler.GetRatingByID)
	}

	serverAddress := viper.GetString("SERVERADDRESS")
	router.Run(serverAddress)
}
