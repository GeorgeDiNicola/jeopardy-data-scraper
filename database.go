package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createDatabaseIfDoesNotExist() error {
	dbHost, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	dbName, dbPort, dbTimezone := os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s",
		dbHost, dbPort, dbUsername, dbPassword, dbTimezone)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	_, _ = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))

	return nil
}

func createJeopardyGameBoxScoreTable() error {
	dbHost, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	dbName, dbPort, dbTimezone := os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE")

	gormDB, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUsername, dbPassword, dbName, dbTimezone)), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = gormDB.AutoMigrate(&JeopardyGameBoxScore{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return nil
}

func getMostRecentEpisodeDate(gormDB *gorm.DB) (string, error) {
	var mostRecentBoxScore JeopardyGameBoxScore

	result := gormDB.Order("episode_date DESC").First(&mostRecentBoxScore)
	if result.Error != nil {
		panic("failed to connect to database")
	}

	return mostRecentBoxScore.EpisodeDate, nil
}

func getMostRecentEpisodeNumber(gormDB *gorm.DB) (string, error) {
	var mostRecentBoxScore JeopardyGameBoxScore

	result := gormDB.Order("episode_number DESC").First(&mostRecentBoxScore)
	if result.Error != nil {
		panic("failed to connect to database")
	}

	return mostRecentBoxScore.EpisodeNumber, nil
}

func saveJeopardyGameBoxScore(scores []JeopardyGameBoxScore) {
	dbHost, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	dbName, dbPort, dbTimezone := os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE")

	gormDB, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUsername, dbPassword, dbName, dbTimezone)), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	result := gormDB.Create(&scores)
	if result.Error != nil {
		log.Printf("failed to insert records: %v", result.Error)
	} else {
		fmt.Println("Records inserted successfully")
	}
}
