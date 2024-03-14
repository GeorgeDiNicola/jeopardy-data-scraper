package main

import (
	"encoding/csv"
	"fmt"
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
			strconv.Itoa(game.TotalDdCorrect),
			strconv.Itoa(game.TotalDdIncorrect),
			strconv.Itoa(game.TotalDdWinnings),
			strconv.Itoa(game.FinalScore),
			strconv.Itoa(game.TotalTripleStumpers),
		}

		if err := writer.Write(record); err != nil {
			panic("Cannot write record")
		}
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
			game.TotalDdCorrect,
			game.TotalDdIncorrect,
			game.TotalDdWinnings,
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
