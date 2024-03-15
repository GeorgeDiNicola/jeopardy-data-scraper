package main

import (
	"log"

	"github.com/xuri/excelize/v2"
)

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
			score.StartingFjScore, score.FjWager, score.FjFinalScore,
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
