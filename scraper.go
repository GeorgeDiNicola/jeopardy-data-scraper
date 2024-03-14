package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getGameTotals(doc *goquery.Document, episode string) ([]JeopardyGameBoxScoreTotal, error) {
	var boxScoreTotals []JeopardyGameBoxScoreTotal

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .game-totals", episode)
	doc.Find(query).Each(func(index int, round *goquery.Selection) {
		// first row is the header and skipping it
		if index%4 == 0 {
			return
		}

		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()
		name := firstName + " " + lastName

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
			fmt.Println(dd)
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

		fmt.Printf("F Contestant: %d, %s\n", index, name)
		tempTotal := JeopardyGameBoxScoreTotal{
			EpisodeNumber:       "9126",
			Date:                "",
			LastName:            lastName,
			FirstName:           firstName,
			City:                "",
			State:               "",
			GameChampion:        "",
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

	boxScoreTotals, _ := getGameTotals(doc, "ep-9121")

	for _, boxScore := range boxScoreTotals {
		fmt.Printf("EpisodeNumber: %s\nDate: %s\nLastName: %s\nFirstName: %s\nCity: %s\nState: %s\nGameChampion: %s\nTotalAtt: %d\nTotalBuz: %d\nTotalBuzPercentage: %.2f\nTotalCorrect: %d\nTotalIncorrect: %d\nCorrectPercentage: %.2f\nTotalDdCorrect: %d\nTotalDdIncorrect: %d\nTotalDdWinnings: %d\nFinalScore: %d\nTotalTripleStumpers: %d\n",
			boxScore.EpisodeNumber,
			boxScore.Date,
			boxScore.LastName,
			boxScore.FirstName,
			boxScore.City,
			boxScore.State,
			boxScore.GameChampion,
			boxScore.TotalAtt,
			boxScore.TotalBuz,
			boxScore.TotalBuzPercentage,
			boxScore.TotalCorrect,
			boxScore.TotalIncorrect,
			boxScore.CorrectPercentage,
			boxScore.TotalDdCorrect,
			boxScore.TotalDdIncorrect,
			boxScore.TotalDdWinnings,
			boxScore.FinalScore,
			boxScore.TotalTripleStumpers,
		)
	}
}
