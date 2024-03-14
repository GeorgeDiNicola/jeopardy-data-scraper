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
)

func getContestantInformation(doc *goquery.Document, episode string) ([]Contestant, error) {
	var contestants []Contestant
	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .contestant", episode)
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		lastName := strings.TrimSpace(s.Find(".name-1").Text())
		firstName := strings.TrimSpace(s.Find(".name-0").Text())
		home := strings.TrimSpace(s.Find(".home").Text())
		homeCityState := strings.Split(home, ", ")
		homeCity, homeState := strings.TrimSpace(homeCityState[0]), strings.TrimSpace(homeCityState[1])

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

		finalScore, _ := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(round.Find("td[data-header='Final Score']").Text()), "$"))

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

	response, err := http.Get("https://www.jeopardy.com/track/jeopardata")
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
		EpisodeNumber string
		Date          string
	}

	// Find each episode and extract its number and date
	doc.Find(".episode").Each(func(i int, s *goquery.Selection) {
		episodeNumber, _ := s.Attr("id")
		date, _ := s.Attr("data-weekday")
		episodes = append(episodes, struct {
			EpisodeNumber string
			Date          string
		}{EpisodeNumber: episodeNumber, Date: date})
	})

	episodeId := "ep-9121"

	contestants, _ := getContestantInformation(doc, episodeId)

	boxScoreTotals, _ := getGameTotals(doc, episodeId)

	for i := 0; i < len(contestants); i++ {
		boxScoreTotals[i].City = contestants[i].HomeCity
		boxScoreTotals[i].State = contestants[i].HomeState
		boxScoreTotals[i].GameWinner = contestants[i].GameWinner

	}

	date := doc.Find("table[aria-labelledby='ep-9121-label'] .date").Text()
	title := doc.Find("table[aria-labelledby='ep-9121-label'] .title").Text()
	episodeNumber := strings.Split(episodeId, "-")[1]
	for i := 0; i < len(boxScoreTotals); i++ {
		// TODO: maybe assign these elsewhere
		boxScoreTotals[i].EpisodeNumber = episodeNumber
		boxScoreTotals[i].Date = date
		boxScoreTotals[i].EpisodeTitle = title
		// fmt.Printf("EpisodeNumber: %s\nDate: %s\nLastName: %s\nFirstName: %s\nCity: %s\nState: %s\nGameWinner: %v\nTotalAtt: %d\nTotalBuz: %d\nTotalBuzPercentage: %.2f\nTotalCorrect: %d\nTotalIncorrect: %d\nCorrectPercentage: %.2f\nTotalDdCorrect: %d\nTotalDdIncorrect: %d\nTotalDdWinnings: %d\nFinalScore: %d\nTotalTripleStumpers: %d\n",
		// 	boxScoreTotals[i].EpisodeNumber,
		// 	boxScoreTotals[i].Date,
		// 	boxScoreTotals[i].LastName,
		// 	boxScoreTotals[i].FirstName,
		// 	boxScoreTotals[i].City,
		// 	boxScoreTotals[i].State,
		// 	boxScoreTotals[i].GameWinner,
		// 	boxScoreTotals[i].TotalAtt,
		// 	boxScoreTotals[i].TotalBuz,
		// 	boxScoreTotals[i].TotalBuzPercentage,
		// 	boxScoreTotals[i].TotalCorrect,
		// 	boxScoreTotals[i].TotalIncorrect,
		// 	boxScoreTotals[i].CorrectPercentage,
		// 	boxScoreTotals[i].TotalDdCorrect,
		// 	boxScoreTotals[i].TotalDdIncorrect,
		// 	boxScoreTotals[i].TotalDdWinnings,
		// 	boxScoreTotals[i].FinalScore,
		// 	boxScoreTotals[i].TotalTripleStumpers,
		// )
	}

	writeGameTotalsOutToCsv(boxScoreTotals)

}
