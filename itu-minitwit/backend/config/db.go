package config

import (
	"fmt"
	"log"
	"minitwit-backend/init/models"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect_prod_db() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("err = ", err)
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		fmt.Println("failed to convert db port to int")
	}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s client_encoding=%s",
		dbHost,
		dbPortInt,
		dbUser,
		dbPassword,
		dbName,
		"require",
		"Europe/Berlin",
		"UTF8")
	db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//Migrate schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Follower{})

	DB = db
}

func Connect_test_db() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("err = ", err)
		log.Fatal("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		fmt.Println("failed to convert db port to int")
	}
	dbName := os.Getenv("DB_TEST_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s client_encoding=%s",
		dbHost,
		dbPortInt,
		dbUser,
		dbPassword,
		dbName,
		"require",
		"Europe/Berlin",
		"UTF8")
	db, err := gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//Migrate schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Follower{})

	DB = db
}
