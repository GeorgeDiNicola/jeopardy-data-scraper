package util

import (
	"georgedinicola/jeopardy-data-scraper/internal/model"
	"log"
	"strings"

	"github.com/xuri/excelize/v2"
)

var UnitedStatesMap = map[string]string{
	"AL": "ALABAMA", "AK": "ALASKA", "AZ": "ARIZONA", "AR": "ARKANSAS", "CA": "CALIFORNIA",
	"CO": "COLORADO", "CT": "CONNECTICUT", "DE": "DELAWARE", "FL": "FLORIDA", "GA": "GEORGIA",
	"HI": "HAWAII", "ID": "IDAHO", "IL": "ILLINOIS", "IN": "INDIANA", "IA": "IOWA",
	"KS": "KANSAS", "KY": "KENTUCKY", "LA": "LOUISIANA", "ME": "MAINE", "MD": "MARYLAND",
	"MA": "MASSACHUSETTS", "MI": "MICHIGAN", "MN": "MINNESOTA", "MS": "MISSISSIPPI", "MO": "MISSOURI",
	"MT": "MONTANA", "NE": "NEBRASKA", "NV": "NEVADA", "NH": "NEW HAMPSHIRE", "NJ": "NEW JERSEY",
	"NM": "NEW MEXICO", "NY": "NEW YORK", "NC": "NORTH CAROLINA", "ND": "NORTH DAKOTA", "OH": "OHIO",
	"OK": "OKLAHOMA", "OR": "OREGON", "PA": "PENNSYLVANIA", "RI": "RHODE ISLAND", "SC": "SOUTH CAROLINA",
	"SD": "SOUTH DAKOTA", "TN": "TENNESSEE", "TX": "TEXAS", "UT": "UTAH", "VT": "VERMONT",
	"VA": "VIRGINIA", "WA": "WASHINGTON", "WV": "WEST VIRGINIA", "WI": "WISCONSIN", "WY": "WYOMING",
	"D.C.": "DC", "d.c": "DC", "D.C": "DC",
}

func GetStateFullName(input string) string {
	upInput := strings.ToUpper(input)
	if fullName, ok := UnitedStatesMap[upInput]; ok {
		return fullName
	} else {
		return upInput
	}
}

func WriteBoxScoreHistoryToExcel(filePath string, scores []model.JeopardyGameBoxScore) error {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	headers := []string{
		"Episode Number", "Episode Title", "Episode Date",
		"Contestant Last Name", "Contestant First Name", "Home City", "Home State", "Is Winner",
		"Round One Attempts", "Round One Buzzes", "Round One Buzz Percentage",
		"Round One Correct Answers", "Round One Incorrect Answers", "Round One Correct Answer Percentage",
		"Round One Daily Doubles", "Round One Score",
		"Round Two Attempts", "Round Two Buzzes", "Round Two Buzz Percentage",
		"Round Two Correct Answers", "Round Two Incorrect Answers", "Round Two Correct Answer Percentage",
		"Round Two Daily Double 1", "Round Two Daily Double 2", "Round Two Score",
		"Final Jeopardy Starting Score", "Final Jeopardy Wager", "Final Jeopardy Score",
		"Total Game Attempts", "Total Game Buzzes", "Total Game Buzz Percentage",
		"Total Game Correct Answers", "Total Game Incorrect Answers", "Total Game Correct Answer Percentage",
		"Total Game Daily Doubles Correct", "Total Game Daily Doubles Incorrect", "Total Game Daily Double Winnings",
		"Total Game Score", "Total Triple Stumpers",
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
			score.RoundOneAttempts, score.RoundOneBuzzes, score.RoundOneBuzzPercent, score.RoundOneCorrectAnswers, score.RoundOneIncorrectAnswers, score.RoundOneCorrectAnswerPercent, score.RoundOneDailyDoubles, score.RoundOneScore,
			score.RoundTwoAttempts, score.RoundTwoBuzzes, score.RoundTwoBuzzPercent, score.RoundTwoCorrectAnswers, score.RoundTwoIncorrectAnswers, score.RoundTwoCorrectAnswerPercent, score.RoundTwoDailyDouble1, score.RoundTwoDailyDouble2, score.RoundTwoScore,
			score.FinalJeopardyStartingScore, score.FinalJeopardyWager, score.FinalJeopardyScore,
			score.TotalGameAttempts, score.TotalGameBuzzes, score.TotalGameBuzzPercent, score.RoundTwoCorrectAnswers, score.TotalGameIncorrectAnswers, score.TotalGameCorrectAnswerPercent,
			score.TotalGameDailyDoublesCorrect, score.TotalGameDailyDoublesIncorrect, score.TotalGameDailyDoubleWinnings,
			score.TotalGameScore, score.TotalTripleStumpers,
		}

		for j, value := range values {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	if err := f.SaveAs(filePath); err != nil {
		log.Printf("Failed to save the Excel file: %v", err)
		return err
	}

	return nil
}
