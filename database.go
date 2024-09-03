package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(postgres.Open("user=dbuser password=123 dbname=mydb host=localhost port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	err = DB.AutoMigrate(&Item{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
