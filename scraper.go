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

func getJeopardyRound(doc *goquery.Document, episode string) ([]JeopardyRound, error) {
	var jeopardyRound []JeopardyRound

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .jeopardy-round", episode)
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
		buzPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text()), "%"))

		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		var correct, incorrect int
		if len(corInc) > 1 {
			correctAndIncorrect := strings.Split(corInc, "/")
			correct, _ = strconv.Atoi(correctAndIncorrect[0])
			incorrect, _ = strconv.Atoi(correctAndIncorrect[1])
		}

		correctPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text()), "%"))

		ddDirty := strings.TrimSpace(round.Find("td[data-header='DD']").Text())
		ddDirty = strings.Replace(ddDirty, "$", "", -1)
		ddClean := strings.Replace(ddDirty, ",", "", -1)
		dd, _ := strconv.Atoi(ddClean)

		eorScoreDirty := strings.TrimSpace(round.Find("td[data-header='EOR SCORE']").Text())
		eorScoreDirty = strings.Replace(eorScoreDirty, "$", "", -1)
		eorScoreClean := strings.Replace(eorScoreDirty, ",", "", -1)
		eorScore, _ := strconv.Atoi(eorScoreClean)

		newRound := JeopardyRound{
			LastName:          lastName,
			FirstName:         firstName,
			Att:               att,
			Buz:               buz,
			BuzPercentage:     buzPercent,
			Correct:           correct,
			Incorrect:         incorrect,
			CorrectPercentage: correctPercent,
			DailyDouble:       dd,
			EorScore:          eorScore,
		}
		index += 1
		jeopardyRound = append(jeopardyRound, newRound)
	})

	return jeopardyRound, nil
}

func getDoubleJeopardyRound(doc *goquery.Document, episode string) ([]DoubleJeopardyRound, error) {
	var doubleJeopardyRound []DoubleJeopardyRound

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .double-jeopardy", episode)
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
		buzPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text()), "%"))

		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		var correct, incorrect int
		if len(corInc) > 1 {
			correctAndIncorrect := strings.Split(corInc, "/")
			correct, _ = strconv.Atoi(correctAndIncorrect[0])
			incorrect, _ = strconv.Atoi(correctAndIncorrect[1])
		}

		correctPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text()), "%"))

		dd1Dirty := strings.TrimSpace(round.Find("td[data-header='DD#1']").Text())
		fmt.Println("HERE")
		fmt.Println(dd1Dirty)
		dd1Dirty = strings.Replace(dd1Dirty, "$", "", -1)
		dd1Clean := strings.Replace(dd1Dirty, ",", "", -1)
		dd1, _ := strconv.Atoi(dd1Clean)

		dd2Dirty := strings.TrimSpace(round.Find("td[data-header='DD#2']").Text())
		dd2Dirty = strings.Replace(dd2Dirty, "$", "", -1)
		dd2Clean := strings.Replace(dd2Dirty, ",", "", -1)
		dd2, _ := strconv.Atoi(dd2Clean)

		eorScoreDirty := strings.TrimSpace(round.Find("td[data-header='EOR SCORE']").Text())
		eorScoreDirty = strings.Replace(eorScoreDirty, "$", "", -1)
		eorScoreClean := strings.Replace(eorScoreDirty, ",", "", -1)
		eorScore, _ := strconv.Atoi(eorScoreClean)

		newRound := DoubleJeopardyRound{
			LastName:          lastName,
			FirstName:         firstName,
			Att:               att,
			Buz:               buz,
			BuzPercentage:     buzPercent,
			Correct:           correct,
			Incorrect:         incorrect,
			CorrectPercentage: correctPercent,
			DailyDouble1:      dd1,
			DailyDouble2:      dd2,
			EorScore:          eorScore,
		}
		index += 1
		doubleJeopardyRound = append(doubleJeopardyRound, newRound)
	})

	return doubleJeopardyRound, nil
}

func getFinalJeopardyRound(doc *goquery.Document, episode string) ([]FinalJeopardyRound, error) {
	var finalJeopardyRound []FinalJeopardyRound

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .final-jeopardy", episode)
	doc.Find(query).Each(func(index int, round *goquery.Selection) {
		// first row is the header, skipping it
		if index%4 == 0 {
			return
		}

		// extract the fields for each column
		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()

		startingFjScoreDirty := strings.TrimSpace(round.Find("td[data-header='Starting']").Text())
		startingFjScoreDirty = strings.Replace(startingFjScoreDirty, "$", "", -1)
		startingFjScoreClean := strings.Replace(startingFjScoreDirty, ",", "", -1)
		startingFjScore, _ := strconv.Atoi(startingFjScoreClean)

		fjWagerDirty := strings.TrimSpace(round.Find("td[data-header='FJ! Wager']").Text())
		fjWagerDirty = strings.Replace(fjWagerDirty, "$", "", -1)
		fjWagerClean := strings.Replace(fjWagerDirty, ",", "", -1)
		fjWager, _ := strconv.Atoi(fjWagerClean)

		finalScoreDirty := strings.TrimSpace(round.Find("td[data-header='Final Score']").Text())
		finalScoreDirty = strings.Replace(finalScoreDirty, "$", "", -1)
		finalScoreClean := strings.Replace(finalScoreDirty, ",", "", -1)
		finalScore, _ := strconv.Atoi(finalScoreClean)

		newRound := FinalJeopardyRound{
			LastName:        lastName,
			FirstName:       firstName,
			StartingFjScore: startingFjScore,
			FjWager:         fjWager,
			FinalScore:      finalScore,
		}
		index += 1
		finalJeopardyRound = append(finalJeopardyRound, newRound)

		fmt.Printf("Contestant: %s, %s\nStarting Final Jeopardy Score: %d\nFinal Jeopardy Wager: %d\nFinal Score: %d\n",
			newRound.LastName, newRound.FirstName,
			newRound.StartingFjScore,
			newRound.FjWager,
			newRound.FinalScore,
		)
	})

	return finalJeopardyRound, nil
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
		buzPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text()), "%"))

		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		var correct, incorrect int
		if len(corInc) > 1 {
			correctAndIncorrect := strings.Split(corInc, "/")
			correct, _ = strconv.Atoi(correctAndIncorrect[0])
			incorrect, _ = strconv.Atoi(correctAndIncorrect[1])
		}

		correctPercent, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text()), "%"))

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

func main() {
	var allEpisodeBoxScoreTotals []JeopardyGameBoxScoreTotal
	//var allEpisodeJeopardyGameBoxScores []JeopardyGameBoxScore

	totalNumberOfWebPages := 73
	for i := 0; i <= totalNumberOfWebPages; i++ {
		fmt.Printf("scraping data from page: %d", i)
		fmt.Println("...")
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

			jeopardyRounds, _ := getJeopardyRound(doc, episode.EpisodeID)
			fmt.Println(jeopardyRounds)

			doubleJeopardyRounds, _ := getDoubleJeopardyRound(doc, episode.EpisodeID)
			fmt.Println(doubleJeopardyRounds)

			finalJeopardyRounds, _ := getFinalJeopardyRound(doc, episode.EpisodeID)
			fmt.Println(finalJeopardyRounds)

			// get all of the Game totals
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

				// add it to the final output for the CSV
				allEpisodeBoxScoreTotals = append(allEpisodeBoxScoreTotals, boxScoreTotals[i])
			}

			// get all of the round details
		}
	}

	writeGameTotalsOutToCsv(allEpisodeBoxScoreTotals)

	writeGameTotalsOutToExcel(allEpisodeBoxScoreTotals)

}
