package main

import (
	"fmt"

	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	apiService "github.com/gvriofernando/test_synergizetech/app"
	"github.com/gvriofernando/test_synergizetech/config/gorm"
	"github.com/gvriofernando/test_synergizetech/config/header"
	"github.com/gvriofernando/test_synergizetech/config/redis"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Env initializations
	redisDatabase := os.Getenv("REDIS_DATABASE")
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	serverPort := os.Getenv("PORT")

	fmt.Print("Application starting up... \n")

	// Initiate Redis Client
	redisDB, err := strconv.Atoi(redisDatabase)
	if err != nil {
		log.Println("Failed convert redis database: ", err.Error())
		os.Exit(1)
	}
	rd, err := redis.NewClient(redis.Config{
		Address:  redisAddress,
		Password: redisPassword,
		Database: redisDB,
	})

	if err != nil {
		log.Println("initiate redis client failed", err.Error())
		os.Exit(1)
	}

	// Initiate Postgres database
	dbCon := gorm.Init(gorm.Config{
		User:     dbUser,
		Password: dbPassword,
		Host:     dbHost,
		Port:     dbPort,
		Database: dbName,
	})
	if err != nil {
		log.Println("Failed initiating Postgres database", "errorMsg", err.Error())
	}
	fmt.Println("Database initiation successful")

	// Initiate Gin Router
	httpS := gin.Default()
	ginConfig := apiService.Config{
		Rd: rd,
		Db: dbCon,
	}
	httpS.Use(header.CORSAllowHeaders())
	apiService.GinHttpRouter(ginConfig, httpS)

	httpS.Run(":" + serverPort)
}
