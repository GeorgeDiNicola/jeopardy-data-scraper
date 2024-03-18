package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func writeToPostgresDB(scores []JeopardyGameBoxScore) {
	dbHost, dbUsername, dbPassword := os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")
	dbName, dbPort, dbTimezone := os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_TIMEZONE")

	// ensure the DB exists
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable TimeZone=%s",
		dbHost, dbPort, dbUsername, dbPassword, dbTimezone)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, _ = db.Exec("CREATE DATABASE jeopardata")

	// Connect to the DB and output the results
	gormDB, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUsername, dbPassword, dbName, dbTimezone)), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = gormDB.AutoMigrate(&JeopardyGameBoxScore{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	// Create records in the database
	result := gormDB.Create(&scores)
	if result.Error != nil {
		log.Printf("failed to insert records: %v", result.Error)
	} else {
		fmt.Println("Records inserted successfully")
	}
}

func writeBoxScoreHistoryToExcel(scores []JeopardyGameBoxScore) {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	headers := []string{
		"Episode Number", "Episode Title", "Date", "Last Name", "First Name", "Home City", "Home State", "Game Winner",
		"R1 Att", "R1 Buz", "R1 Buz Percentage", "R1 Correct", "R1 Incorrect", "R1 Correct Percentage", "R1 Daily Double", "R1 Eor",
		"R2 Att", "R2 Buz", "R2 Buz Percentage", "R2 Correct", "R2 Incorrect", "R2 Correct Percentage", "R2 Daily Double 1", "R2 Daily Double 2", "R2 Eor",
		"Starting FJ Score", "FJ Wager", "Final Score",
		"Att Total", "Buz Total", "Buz Percentage Total", "Correct Total", "Incorrect Total", "Correct Percentage Total",
		"Daily Double Correct Total", "Daily Double Incorrect Total", "Daily Double Winnings Total",
		"Final Score Total", "Total Triple Strumpers", "Coryat Score",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, score := range scores {
		gameWinnerStr := "No"
		if score.IsWinner {
			gameWinnerStr = "Yes"
		}
		values := []interface{}{
			score.EpisodeNumber, score.EpisodeTitle, score.EpisodeDate, score.ContestantLastName, score.ContestantFirstName, score.HomeCity, score.HomeState, gameWinnerStr,
			score.JeopardyAttempts, score.JeopardyBuzzes, score.JeopardyBuzzPercentage, score.JeopardyCorrectAnswers, score.JeopardyIncorrectAnswers, score.JeopardyCorrectAnswerPercentage, score.JeopardyDailyDoublesFound, score.JeopardyScore,
			score.DoubleJeopardyAttempts, score.DoubleJeopardyBuzzes, score.DoubleJeopardyBuzzPercentage, score.DoubleJeopardyCorrectAnswers, score.DoubleJeopardyIncorrectAnswers, score.DoubleJeopardyCorrectAnswerPercentage, score.DoubleJeopardyDailyDouble1Found, score.DoubleJeopardyDailyDouble2Found, score.DoubleJeopardyScore,
			score.FinalJeopardyStartingScore, score.FinalJeopardyWager, score.FinalJeopardyScore,
			score.TotalAttempts, score.TotalBuzzes, score.TotalBuzzPercentage, score.TotalCorrectAnswers, score.TotalIncorrectAnswers, score.TotalCorrectAnswerPercentage,
			score.TotalDailyDoublesCorrect, score.TotalDailyDoublesIncorrect, score.TotalDailyDoubleWinnings,
			score.TotalScore, score.TotalTripleStumpers, score.CoryatScore,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	if err := f.SaveAs("all_jeopardy_box_scores.xlsx"); err != nil {
		log.Printf("Failed to save the Excel file: %v", err)
	}
}
