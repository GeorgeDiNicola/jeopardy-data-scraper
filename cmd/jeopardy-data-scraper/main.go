package main

import (
	"georgedinicola/jeopardy-data-scraper/internal/db"
	"georgedinicola/jeopardy-data-scraper/internal/model"
	"georgedinicola/jeopardy-data-scraper/internal/scraper"
	"georgedinicola/jeopardy-data-scraper/internal/util"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <mode>", os.Args[0])
	}
	mode := os.Args[1]

	numberOfPages := 73 // TODO: add function that goes and gets max # of pages
	var jeopardyBoxScores []model.JeopardyGameBoxScore

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

	if mode == "FULL" {
		err = db.CreateJeopardyGameBoxScoreTable()
		if err != nil {
			log.Fatalf("failed to create the table: %v", err)
			return
		}

		jeopardyBoxScores = scraper.ScrapeGameDataFull(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			db.SaveJeopardyGameBoxScore(jeopardyBoxScores)
			util.WriteBoxScoreHistoryToExcel("jeopardata_box_scores_sample.xlsx", jeopardyBoxScores)
		} else {
			log.Println("no jeopardata records to extract")
		}

	} else if mode == "INCREMENTAL" {
		mostRecentEpisodeNum, err := db.GetMostRecentEpisodeNumber()
		if err != nil {
			log.Fatal("Error querying for the most recent episode date: ", err)
		}

		jeopardyBoxScores = scraper.ScrapeGameDataIncremental(mostRecentEpisodeNum, numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			db.SaveJeopardyGameBoxScore(jeopardyBoxScores)
			util.WriteBoxScoreHistoryToExcel("jeopardata_box_scores_sample.xlsx", jeopardyBoxScores)
		} else {
			log.Println("no new jeopardata records to extract")
		}
	} else {
		log.Fatalf("invalid mode: %s", mode)
		return
	}
}
