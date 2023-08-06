package gorm

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}
type SliceConfig []Config

func Init(cfg Config) *gorm.DB {
	var err error
	user := cfg.User
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	database := cfg.Database

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		database,
	)

	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	log.Println("Database connected")

	return db
}
