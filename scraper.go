package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/xuri/excelize/v2"
)

func getContestantInformation(doc *goquery.Document, episode string) ([]Contestant, error) {
	var contestants []Contestant
	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .contestant", episode)
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		lastName := strings.TrimSpace(s.Find(".name-1").Text())
		firstName := strings.TrimSpace(s.Find(".name-0").Text())

		home := strings.TrimSpace(s.Find(".home").Text())
		var homeCity, homeState string
		if len(home) > 1 {
			homeCityState := strings.Split(home, ", ")

			if len(homeCityState) > 1 {
				homeCity, homeState = strings.TrimSpace(homeCityState[0]), strings.TrimSpace(homeCityState[1])
			}
		}

		contestant := Contestant{
			FirstName: firstName,
			LastName:  lastName,
			HomeCity:  homeCity,
			HomeState: homeState,
		}
		contestants = append(contestants, contestant)
	})

	var winnerLastName, winnerFirstName string
	query = fmt.Sprintf("table[aria-labelledby='%s-label'] .winner", episode)
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		winnerLastName = strings.TrimSpace(s.Find(".name-1").Text())
		winnerFirstName = strings.TrimSpace(s.Find(".name-0").Text())
	})

	for i := 0; i < len(contestants); i++ {
		if contestants[i].FirstName == winnerFirstName && contestants[i].LastName == winnerLastName {
			contestants[i].GameWinner = true
		} else {
			contestants[i].GameWinner = false
		}
	}

	return contestants, nil

}

func getGameTotals(doc *goquery.Document, episode string) ([]JeopardyGameBoxScoreTotal, error) {
	var boxScoreTotals []JeopardyGameBoxScoreTotal

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .game-totals", episode)
	doc.Find(query).Each(func(index int, round *goquery.Selection) {
		// first row is the header, skipping it
		if index%4 == 0 {
			return
		}

		// extract the fields for each column
		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()
		att, _ := strconv.Atoi(strings.TrimSpace(round.Find("td[data-header='ATT']").Text()))
		buz, _ := strconv.Atoi(strings.TrimSpace(round.Find("td[data-header='BUZ']").Text()))
		buzPercent, _ := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text()), "%"), 64)

		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		var correct, incorrect int
		if len(corInc) > 1 {
			correctAndIncorrect := strings.Split(corInc, "/")
			correct, _ = strconv.Atoi(correctAndIncorrect[0])
			incorrect, _ = strconv.Atoi(correctAndIncorrect[1])
		}

		correctPercent, _ := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text()), "%"), 64)

		dd := strings.TrimSpace(round.Find("td[data-header='DD (COR/INC)']").Text())
		var ddCorrect, ddIncorrect, ddWinnings int
		if len(dd) > 1 {
			// Normalize the string by replacing tabs and newlines with a space
			normalizedData := strings.ReplaceAll(dd, "\n", " ")
			normalizedData = strings.ReplaceAll(normalizedData, "\t", " ")

			parts := strings.Fields(normalizedData)
			fractionPart := parts[0]
			fractionParts := strings.Split(fractionPart, "/")
			ddCorrect, _ = strconv.Atoi(fractionParts[0])
			ddIncorrect, _ = strconv.Atoi(fractionParts[1])

			if len(parts) > 1 {
				monetaryPart := parts[1]
				monetaryValue := strings.Replace(monetaryPart, "$", "", -1)
				monetaryValue = strings.Replace(monetaryValue, ",", "", -1)
				ddWinnings, _ = strconv.Atoi(monetaryValue)
			}

		}

		finalScoreDirty := strings.TrimSpace(round.Find("td[data-header='Final Score']").Text())
		finalScoreDirty = strings.Replace(finalScoreDirty, "$", "", -1)
		finalScoreClean := strings.Replace(finalScoreDirty, ",", "", -1)
		finalScore, _ := strconv.Atoi(finalScoreClean)

		tempTotal := JeopardyGameBoxScoreTotal{
			EpisodeNumber:       episode,
			Date:                "",
			LastName:            lastName,
			FirstName:           firstName,
			City:                "",
			State:               "",
			TotalAtt:            att,
			TotalBuz:            buz,
			TotalBuzPercentage:  buzPercent,
			TotalCorrect:        correct,
			TotalIncorrect:      incorrect,
			CorrectPercentage:   correctPercent,
			TotalDdCorrect:      ddCorrect,
			TotalDdIncorrect:    ddIncorrect,
			TotalDdWinnings:     ddWinnings,
			FinalScore:          finalScore,
			TotalTripleStumpers: 0,
		}
		index += 1
		boxScoreTotals = append(boxScoreTotals, tempTotal)
	})

	return boxScoreTotals, nil
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
			strconv.FormatFloat(game.TotalBuzPercentage, 'f', 2, 64),
			strconv.Itoa(game.TotalCorrect),
			strconv.Itoa(game.TotalIncorrect),
			strconv.FormatFloat(game.CorrectPercentage, 'f', 2, 64),
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

func main() {
	// TODO: loop to next page
	// https://www.jeopardy.com/track/jeopardata?page=1
	var allEpisodeBoxScoreTotals []JeopardyGameBoxScoreTotal

	totalNumberOfWebPages := 73
	for i := 1; i < totalNumberOfWebPages+1; i++ {
		fmt.Printf("scraping data from page: %d", i)
		url := fmt.Sprintf("https://www.jeopardy.com/track/jeopardata?page=%d", i)
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			log.Fatalf("Error fetching the page: %s", response.Status)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		htmlPageContent := string(body)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlPageContent))
		if err != nil {
			log.Fatal(err)
		}

		// collect all of the episodes on the webpage
		var episodes []struct {
			EpisodeID string
			Date      string
		}

		// Find each episode and extract its number and date
		doc.Find(".episode").Each(func(i int, s *goquery.Selection) {
			episodeID, _ := s.Attr("id")
			date, _ := s.Attr("data-weekday")
			episodes = append(episodes, struct {
				EpisodeID string
				Date      string
			}{EpisodeID: episodeID, Date: date})
		})

		// Get the data for each episode
		for _, episode := range episodes {
			// TODO: maybe use the date above
			contestants, _ := getContestantInformation(doc, episode.EpisodeID)

			boxScoreTotals, _ := getGameTotals(doc, episode.EpisodeID)

			for i := 0; i < len(contestants); i++ {
				boxScoreTotals[i].City = contestants[i].HomeCity
				boxScoreTotals[i].State = contestants[i].HomeState
				boxScoreTotals[i].GameWinner = contestants[i].GameWinner

			}

			date := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .date", episode.EpisodeID)).Text()
			title := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .title", episode.EpisodeID)).Text()
			episodeNumber := strings.Split(episode.EpisodeID, "-")[1]
			for i := 0; i < len(boxScoreTotals); i++ {
				// TODO: maybe assign these elsewhere
				boxScoreTotals[i].EpisodeNumber = episodeNumber
				boxScoreTotals[i].Date = date
				boxScoreTotals[i].EpisodeTitle = title
				fmt.Printf("EpisodeNumber: %s\nDate: %s\nLastName: %s\nFirstName: %s\nCity: %s\nState: %s\nGameWinner: %v\nTotalAtt: %d\nTotalBuz: %d\nTotalBuzPercentage: %.2f\nTotalCorrect: %d\nTotalIncorrect: %d\nCorrectPercentage: %.2f\nTotalDdCorrect: %d\nTotalDdIncorrect: %d\nTotalDdWinnings: %d\nFinalScore: %d\nTotalTripleStumpers: %d\n",
					boxScoreTotals[i].EpisodeNumber,
					boxScoreTotals[i].Date,
					boxScoreTotals[i].LastName,
					boxScoreTotals[i].FirstName,
					boxScoreTotals[i].City,
					boxScoreTotals[i].State,
					boxScoreTotals[i].GameWinner,
					boxScoreTotals[i].TotalAtt,
					boxScoreTotals[i].TotalBuz,
					boxScoreTotals[i].TotalBuzPercentage,
					boxScoreTotals[i].TotalCorrect,
					boxScoreTotals[i].TotalIncorrect,
					boxScoreTotals[i].CorrectPercentage,
					boxScoreTotals[i].TotalDdCorrect,
					boxScoreTotals[i].TotalDdIncorrect,
					boxScoreTotals[i].TotalDdWinnings,
					boxScoreTotals[i].FinalScore,
					boxScoreTotals[i].TotalTripleStumpers,
				)

				// add it to the final output for the CSV
				allEpisodeBoxScoreTotals = append(allEpisodeBoxScoreTotals, boxScoreTotals[i])
			}

		}

	}

	writeGameTotalsOutToCsv(allEpisodeBoxScoreTotals)

	writeGameTotalsOutToExcel(allEpisodeBoxScoreTotals)

}
