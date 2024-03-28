package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/georgedinicola/jeopardy-data-scraper/model"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	CreateDatabaseIfDoesNotExist() error
	CreateGormDbConnection() (*gorm.DB, error)
	CreateJeopardyGameBoxScoreTable() error
	GetMostRecentEpisodeNumber() (string, error)
	SaveJeopardyGameBoxScore(scores []model.JeopardyGameBoxScore) error
}

type DatabaseConnx struct {
	gorm *gorm.DB
}

func CreateNewGormDbConnection() (DatabaseConnx, error) {
	dbHost, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	dbName, dbPort, dbTimezone := os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE")

	gormDB, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUsername, dbPassword, dbName, dbTimezone)), &gorm.Config{})
	if err != nil {
		return DatabaseConnx{}, err
	}

	return DatabaseConnx{gormDB}, nil
}

// no usage of gorm. Nneeded for DB creation since gorm only does ORM
func CreateDatabaseIfDoesNotExist() error {
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

func (d *DatabaseConnx) CreateJeopardyGameBoxScoreTable() error {
	err := d.gorm.AutoMigrate(&model.JeopardyGameBoxScore{})
	if err != nil {
		log.Printf("failed to migrate: %v", err)
	}

	return nil
}

func (d *DatabaseConnx) GetMostRecentEpisodeNumber() (string, error) {
	var mostRecentBoxScore model.JeopardyGameBoxScore

	result := d.gorm.Order("episode_number DESC").First(&mostRecentBoxScore)
	if result.Error != nil {
		return "", result.Error
	}

	return mostRecentBoxScore.EpisodeNumber, nil
}

func (d *DatabaseConnx) SaveJeopardyGameBoxScore(scores []model.JeopardyGameBoxScore) error {
	result := d.gorm.Create(&scores)
	if result.Error != nil {
		return result.Error
	} else {
		return nil
	}
}
