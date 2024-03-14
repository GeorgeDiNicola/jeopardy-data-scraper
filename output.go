package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func writeGameTotalsOutToCsv(boxScoreTotals []JeopardyGameBoxScoreTotal) {
	// write out to CSV
	file, err := os.Create("jeopardy_scores.csv")
	if err != nil {
		panic("Cannot create file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing the header
	header := []string{"EpisodeNumber", "Date", "LastName", "FirstName", "City", "State", "GameWinner", "TotalAtt", "TotalBuz", "TotalBuzPercentage", "TotalCorrect", "TotalIncorrect", "CorrectPercentage", "TotalDdCorrect", "TotalDdIncorrect", "TotalDdWinnings", "FinalScore", "TotalTripleStumpers"}
	if err := writer.Write(header); err != nil {
		panic("Cannot write header")
	}

	// Writing the data
	for _, game := range boxScoreTotals {
		gameWinnerStr := "false"
		if game.GameWinner {
			gameWinnerStr = "true"
		}
		record := []string{
			game.EpisodeNumber,
			game.Date,
			game.LastName,
			game.FirstName,
			game.City,
			game.State,
			gameWinnerStr,
			strconv.Itoa(game.TotalAtt),
			strconv.Itoa(game.TotalBuz),
			strconv.Itoa(game.TotalBuzPercentage),
			strconv.Itoa(game.TotalCorrect),
			strconv.Itoa(game.TotalIncorrect),
			strconv.Itoa(game.CorrectPercentage),
			strconv.Itoa(game.TotalDailyDoubleCorrect),
			strconv.Itoa(game.TotalDailyDoubleIncorrect),
			strconv.Itoa(game.TotalDailyDoubleWinnings),
			strconv.Itoa(game.FinalScore),
			strconv.Itoa(game.TotalTripleStumpers),
		}

		if err := writer.Write(record); err != nil {
			panic("Cannot write record")
		}
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
		"Final Score Total", "Coryat Score",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, score := range scores {
		gameWinnerStr := "No"
		if score.GameWinner {
			gameWinnerStr = "Yes"
		}
		values := []interface{}{
			score.EpisodeNumber, score.EpisodeTitle, score.Date, score.LastName, score.FirstName, score.HomeCity, score.HomeState, gameWinnerStr,
			score.R1Att, score.R1Buz, score.R1BuzPercentage, score.R1Correct, score.R1Incorrect, score.R1CorrectPercentage, score.R1DailyDouble, score.R1Eor,
			score.R2Att, score.R2Buz, score.R2BuzPercentage, score.R2Correct, score.R2Incorrect, score.R2CorrectPercentage, score.R2DailyDouble1, score.R2DailyDouble2, score.R2Eor,
			score.StartingFjScore, score.FjWager, score.FinalScore,
			score.AttTotal, score.BuzTotal, score.BuzPercentageTotal, score.CorrectTotal, score.IncorrectTotal, score.CorrectPercentageTotal,
			score.DailyDoubleCorrectTotal, score.DailyDoubleIncorrectTotal, score.DailyDoubleWinningsTotal,
			score.FinalScoreTotal, score.CoryatScore,
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

func writeGameTotalsOutToExcel(boxScoreTotals []JeopardyGameBoxScoreTotal) {
	f := excelize.NewFile()
	headers := []string{"EpisodeNumber", "Date", "LastName", "FirstName", "City", "State", "GameWinner", "TotalAtt", "TotalBuz", "TotalBuzPercentage", "TotalCorrect", "TotalIncorrect", "CorrectPercentage", "TotalDdCorrect", "TotalDdIncorrect", "TotalDdWinnings", "FinalScore", "TotalTripleStumpers"}
	for i, title := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, title)
	}

	for i, game := range boxScoreTotals {
		gameWinnerStr := "No"
		if game.GameWinner {
			gameWinnerStr = "Yes"
		}
		values := []interface{}{
			game.EpisodeNumber,
			game.Date,
			game.LastName,
			game.FirstName,
			game.City,
			game.State,
			gameWinnerStr,
			game.TotalAtt,
			game.TotalBuz,
			game.TotalBuzPercentage,
			game.TotalCorrect,
			game.TotalIncorrect,
			game.CorrectPercentage,
			game.TotalDailyDoubleCorrect,
			game.TotalDailyDoubleIncorrect,
			game.TotalDailyDoubleWinnings,
			game.FinalScore,
			game.TotalTripleStumpers,
		}
		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2) // Rows start from 1, but we have a title row
			f.SetCellValue("Sheet1", cell, value)
		}
	}

	if err := f.SaveAs("JeopardyGameBoxScores.xlsx"); err != nil {
		fmt.Println(err)
	}

}
