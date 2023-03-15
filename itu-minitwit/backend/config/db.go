package config

import (
	"fmt"
	"minitwit-backend/init/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect_prod_db() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s client_encoding=%s",
		"cicdont-do-user-13570987-0.b.db.ondigitalocean.com",
		25060,
		"doadmin",
		"AVNS_FeRFl5bSz6UNMVF6Llx",
		"minitwit",
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
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s client_encoding=%s",
		"cicdont-do-user-13570987-0.b.db.ondigitalocean.com",
		25060,
		"doadmin",
		"AVNS_FeRFl5bSz6UNMVF6Llx",
		"minitwit_test",
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
