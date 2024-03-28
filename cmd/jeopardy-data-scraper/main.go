package main

import (
	"log"
	"os"

	"github.com/georgedinicola/jeopardy-data-scraper/internal/config"
	"github.com/georgedinicola/jeopardy-data-scraper/internal/db"
	"github.com/georgedinicola/jeopardy-data-scraper/internal/scraper"
	"github.com/georgedinicola/jeopardy-data-scraper/internal/util"
	"github.com/georgedinicola/jeopardy-data-scraper/model"
)

func main() {
	var jeopardyBoxScores []model.JeopardyGameBoxScore
	numberOfPages := 73 // TODO: add function that goes and gets max # of pages

	// control the mode from the container env
	appMode := os.Getenv("APP_MODE")
	if appMode == "" {
		appMode = "EXCEL"
	}

	// Handle excel mode
	if appMode == "EXCEL" {
		jeopardyBoxScores = scraper.ScrapeGameDataFull(numberOfPages)
		if len(jeopardyBoxScores) > 0 {
			err := util.WriteBoxScoreHistoryToExcel(config.OutputFileName, jeopardyBoxScores)
			if err != nil {
				log.Fatalf("failed to write the file: %v", err)
			}
		} else {
			log.Fatal("no jeopardata records found!")
		}

		os.Exit(0)
		return
	}

	// Handle DB modes
	err := db.CreateDatabaseIfDoesNotExist()
	if err != nil {
		log.Fatalf("failed to create the DB: %v", err)
		return
	}

	db, err := db.CreateNewGormDbConnection()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return
	}

	if appMode == "FULL" {
		err = db.CreateJeopardyGameBoxScoreTable()
		if err != nil {
			log.Fatalf("failed to create the table: %v", err)
			return
		}

		jeopardyBoxScores = scraper.ScrapeGameDataFull(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			db.SaveJeopardyGameBoxScore(jeopardyBoxScores)
		} else {
			log.Println("no jeopardata records to extract")
		}

	} else if appMode == "INCREMENTAL" {
		mostRecentEpisodeNum, err := db.GetMostRecentEpisodeNumber()
		if err != nil {
			log.Fatal("Error querying for the most recent episode date: ", err)
		}

		jeopardyBoxScores = scraper.ScrapeGameDataIncremental(mostRecentEpisodeNum, numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			db.SaveJeopardyGameBoxScore(jeopardyBoxScores)
		} else {
			log.Println("no new jeopardata records to extract")
		}
	} else {
		log.Fatalf("invalid mode: %s", appMode)
		return
	}
}
