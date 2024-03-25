package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// gets all of the Jeopardata from the Jeopardy.com website that the DB does not know about
func ScrapeGameDataIncremental(maxNumPages int) []JeopardyGameBoxScore {
	var jeopardyGameBoxScores []JeopardyGameBoxScore

	mostRecentEpisodeNum, err := getMostRecentEpisodeNumber()
	if err != nil {
		log.Fatal("Error querying for the most recent episode date: ", err)
	}

	// start on page 1 and scrape until that most recent date
	currentPageNumber := 0
	updateFinished := false
	for {
		if maxNumPages > currentPageNumber || updateFinished {
			break
		}

		doc, err := getJeopardataWebPage(currentPageNumber)
		if err != nil {
			log.Fatal("Error getting the Jeopardata web page: ", err)
		}

		// get all of the episodes on the page
		var episodes []Episode
		doc.Find(".episode").Each(func(i int, s *goquery.Selection) {
			episodeID, _ := s.Attr("id")
			date, _ := s.Attr("data-weekday")
			episodes = append(episodes, Episode{
				EpisodeID: episodeID,
				Date:      date,
			})
		})

		// Collect all relevant data for each episode
		for _, episode := range episodes {
			// only scrape up to the last known episode
			episodeNum := strings.Split(episode.EpisodeID, "-")[1]

			if mostRecentEpisodeNum == episodeNum {
				updateFinished = true
				break
			}

			jeopardyGameBoxScore := getJeopardyGameData(doc, episode)

			jeopardyGameBoxScores = append(jeopardyGameBoxScores, jeopardyGameBoxScore...)
		}

		// delay to avoid rate limiting from Jepoardy.com
		time.Sleep(DelayBetweenRequests * time.Second)

		currentPageNumber++ // next page
	}

	return jeopardyGameBoxScores

}

// gets all of the Jeopardata from the Jeopardy.com website
func ScrapeGameDataFull(totalNumberOfPages int) []JeopardyGameBoxScore {
	var allEpisodeJeopardyGameBoxScores []JeopardyGameBoxScore

	for i := 0; i <= totalNumberOfPages; i++ {

		doc, err := getJeopardataWebPage(i)
		if err != nil {
			log.Fatal("Error getting the Jeopardata web page: ", err)
		}

		// get all of the episodes on the page
		var episodes []struct {
			EpisodeID string
			Date      string
		}
		doc.Find(".episode").Each(func(i int, s *goquery.Selection) {
			episodeID, _ := s.Attr("id")
			date, _ := s.Attr("data-weekday")
			episodes = append(episodes, struct {
				EpisodeID string
				Date      string
			}{EpisodeID: episodeID, Date: date})
		})

		// Collect all relevant data for each episode
		for _, episode := range episodes {
			jeopardyGameBoxScore := getJeopardyGameData(doc, episode)

			allEpisodeJeopardyGameBoxScores = append(allEpisodeJeopardyGameBoxScores, jeopardyGameBoxScore...)
		}

		// delay to avoid rate limiting from Jepoardy.com
		time.Sleep(DelayBetweenRequests * time.Second)
	}

	return allEpisodeJeopardyGameBoxScores
}

func getJeopardataWebPage(currentPageNumber int) (*goquery.Document, error) {
	fmt.Printf("scraping data from page: %d ...\n", currentPageNumber)

	url := fmt.Sprintf("https://www.jeopardy.com/track/jeopardata?page=%d", currentPageNumber)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Error fetching the page: %s", response.Status)
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	htmlPageContent := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlPageContent))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getJeopardyGameData(doc *goquery.Document, episode Episode) []JeopardyGameBoxScore {
	// the 3 box scores for the Jeopardy game
	var jeopardyGameBoxScores []JeopardyGameBoxScore

	contestants, _ := getContestantInformation(doc, episode.EpisodeID)
	jeopardyRounds, _ := getJeopardyRound(doc, episode.EpisodeID)
	doubleJeopardyRounds, _ := getDoubleJeopardyRound(doc, episode.EpisodeID)

	boxScoreTotals, _ := getGameTotals(doc, episode.EpisodeID)
	for i := 0; i < len(contestants); i++ {
		boxScoreTotals[i].City = contestants[i].HomeCity
		boxScoreTotals[i].State = contestants[i].HomeState
		boxScoreTotals[i].GameWinner = contestants[i].GameWinner
	}

	finalJeopardyRounds, _ := getFinalJeopardyRound(doc, episode.EpisodeID)
	numberOfTripleStumpers, _ := getNumberOfTripleStumpers(doc, episode.EpisodeID)

	episodeDate := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .date", episode.EpisodeID)).Text()
	episodeTitle := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .title", episode.EpisodeID)).Text()
	episodeNumber := strings.Split(episode.EpisodeID, "-")[1]

	// fill in all of the collected data
	for i := 0; i < len(contestants); i++ {
		var jeopardyGameBoxScore JeopardyGameBoxScore

		// Jeopardy Metadata
		jeopardyGameBoxScore.EpisodeNumber = episodeNumber
		jeopardyGameBoxScore.EpisodeTitle = episodeTitle
		jeopardyGameBoxScore.EpisodeDate = episodeDate
		jeopardyGameBoxScore.ContestantLastName = contestants[i].LastName
		jeopardyGameBoxScore.ContestantFirstName = contestants[i].FirstName
		jeopardyGameBoxScore.HomeCity = contestants[i].HomeCity
		jeopardyGameBoxScore.HomeState = contestants[i].HomeState
		jeopardyGameBoxScore.IsWinner = contestants[i].GameWinner

		// Round 1
		jeopardyGameBoxScore.RoundOneAttempts = jeopardyRounds[i].Attempts
		jeopardyGameBoxScore.RoundOneBuzzes = jeopardyRounds[i].Buzzes
		jeopardyGameBoxScore.RoundOneBuzzPercent = jeopardyRounds[i].BuzzPercentage
		jeopardyGameBoxScore.RoundOneCorrectAnswers = jeopardyRounds[i].Correct
		jeopardyGameBoxScore.RoundOneIncorrectAnswers = jeopardyRounds[i].Incorrect
		jeopardyGameBoxScore.RoundOneCorrectAnswerPercent = jeopardyRounds[i].CorrectPercentage
		jeopardyGameBoxScore.RoundOneDailyDoubles = jeopardyRounds[i].DailyDouble
		jeopardyGameBoxScore.RoundOneScore = jeopardyRounds[i].EndOfRoundScore

		// Double Jeopardy
		jeopardyGameBoxScore.RoundTwoAttempts = doubleJeopardyRounds[i].Attempts
		jeopardyGameBoxScore.RoundTwoBuzzes = doubleJeopardyRounds[i].Buzzes
		jeopardyGameBoxScore.RoundTwoBuzzPercent = doubleJeopardyRounds[i].BuzzPercentage
		jeopardyGameBoxScore.RoundTwoCorrectAnswers = doubleJeopardyRounds[i].Correct
		jeopardyGameBoxScore.RoundTwoIncorrectAnswers = doubleJeopardyRounds[i].Incorrect
		jeopardyGameBoxScore.RoundTwoCorrectAnswerPercent = doubleJeopardyRounds[i].CorrectPercentage
		jeopardyGameBoxScore.RoundTwoDailyDouble1 = doubleJeopardyRounds[i].DailyDouble1
		jeopardyGameBoxScore.RoundTwoDailyDouble2 = doubleJeopardyRounds[i].DailyDouble2
		jeopardyGameBoxScore.RoundTwoScore = doubleJeopardyRounds[i].EndOfRoundScore

		// Final Jeopardy
		jeopardyGameBoxScore.FinalJeopardyStartingScore = finalJeopardyRounds[i].StartingFjScore
		jeopardyGameBoxScore.FinalJeopardyWager = finalJeopardyRounds[i].FjWager
		jeopardyGameBoxScore.FinalJeopardyScore = finalJeopardyRounds[i].FinalScore

		// Round Totals
		jeopardyGameBoxScore.TotalGameAttempts = boxScoreTotals[i].TotalAttempts
		jeopardyGameBoxScore.TotalGameBuzzes = boxScoreTotals[i].TotalBuzzes
		jeopardyGameBoxScore.TotalGameBuzzPercent = boxScoreTotals[i].TotalBuzzPercentage
		jeopardyGameBoxScore.TotalGameCorrectAnswers = boxScoreTotals[i].TotalCorrect
		jeopardyGameBoxScore.TotalGameIncorrectAnswers = boxScoreTotals[i].TotalIncorrect
		jeopardyGameBoxScore.TotalGameCorrectAnswerPercent = boxScoreTotals[i].CorrectPercentage
		jeopardyGameBoxScore.TotalGameDailyDoublesCorrect = boxScoreTotals[i].TotalDailyDoubleCorrect
		jeopardyGameBoxScore.TotalGameDailyDoublesIncorrect = boxScoreTotals[i].TotalDailyDoubleIncorrect
		jeopardyGameBoxScore.TotalGameDailyDoubleWinnings = boxScoreTotals[i].TotalDailyDoubleWinnings
		jeopardyGameBoxScore.TotalGameScore = boxScoreTotals[i].FinalScore

		// Triple Stumpers
		jeopardyGameBoxScore.TotalTripleStumpers = numberOfTripleStumpers

		jeopardyGameBoxScores = append(jeopardyGameBoxScores, jeopardyGameBoxScore)
	}

	return jeopardyGameBoxScores
}

func getContestantInformation(doc *goquery.Document, episode string) ([]Contestant, error) {
	var contestants []Contestant

	query := fmt.Sprintf("table[aria-labelledby='%s-label'] .contestant", episode)
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		lastName := strings.ToUpper(strings.TrimSpace(s.Find(".name-1").Text()))
		firstName := strings.ToUpper(strings.TrimSpace(s.Find(".name-0").Text()))

		home := strings.TrimSpace(s.Find(".home").Text())
		var homeCity, homeState string
		if len(home) > 1 {
			homeCityState := strings.Split(home, ", ")

			if len(homeCityState) > 1 {
				homeCity, homeState = strings.ToUpper(strings.TrimSpace(homeCityState[0])), strings.ToUpper(strings.TrimSpace(homeCityState[1]))
				homeState = getStateFullName(homeState)
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
		winnerLastName = strings.ToUpper(strings.TrimSpace(s.Find(".name-1").Text()))
		winnerFirstName = strings.ToUpper(strings.TrimSpace(s.Find(".name-0").Text()))
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
		firstName := strings.ToUpper(round.Find(".name-0").Text())
		lastName := strings.ToUpper(round.Find(".name-1").Text())
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
			Attempts:          att,
			Buzzes:            buz,
			BuzzPercentage:    buzPercent,
			Correct:           correct,
			Incorrect:         incorrect,
			CorrectPercentage: correctPercent,
			DailyDouble:       dd,
			EndOfRoundScore:   eorScore,
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
		firstName := strings.ToUpper(round.Find(".name-0").Text())
		lastName := strings.ToUpper(round.Find(".name-1").Text())
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
			Attempts:          att,
			Buzzes:            buz,
			BuzzPercentage:    buzPercent,
			Correct:           correct,
			Incorrect:         incorrect,
			CorrectPercentage: correctPercent,
			DailyDouble1:      dd1,
			DailyDouble2:      dd2,
			EndOfRoundScore:   eorScore,
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
		firstName := strings.ToUpper(round.Find(".name-0").Text())
		lastName := strings.ToUpper(round.Find(".name-1").Text())

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
		firstName := strings.ToUpper(round.Find(".name-0").Text())
		lastName := strings.ToUpper(round.Find(".name-1").Text())
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
			EpisodeNumber:             episode,
			LastName:                  lastName,
			FirstName:                 firstName,
			TotalAttempts:             att,
			TotalBuzzes:               buz,
			TotalBuzzPercentage:       buzPercent,
			TotalCorrect:              correct,
			TotalIncorrect:            incorrect,
			CorrectPercentage:         correctPercent,
			TotalDailyDoubleCorrect:   ddCorrect,
			TotalDailyDoubleIncorrect: ddIncorrect,
			TotalDailyDoubleWinnings:  ddWinnings,
			FinalScore:                finalScore,
		}
		index += 1
		boxScoreTotals = append(boxScoreTotals, tempTotal)
	})

	return boxScoreTotals, nil
}

func getNumberOfTripleStumpers(doc *goquery.Document, episode string) (int, error) {
	var numberOfTripleStumpers int

	query := fmt.Sprintf("td[headers='%s-notes']", episode)
	doc.Find(query).Each(func(i int, s *goquery.Selection) {
		gameNotesText := s.Find("p").Text()
		if strings.Contains(gameNotesText, "Triple Stumpers:") {
			re := regexp.MustCompile(`Triple Stumpers:\D*(\d+)(-day)?`)
			matches := re.FindStringSubmatch(gameNotesText)
			numberStr := matches[1]

			// Special case: Check if the "-day" pattern is present
			if len(matches) > 2 && matches[2] == "-day" {
				dayRe := regexp.MustCompile(`(\d+)-day`)
				dayMatches := dayRe.FindStringSubmatch(gameNotesText)
				if len(dayMatches) > 0 {
					numberStr = dayMatches[1][:1]
				}
			}

			numberOfTripleStumpers, _ = strconv.Atoi(numberStr)
		}
	})

	return numberOfTripleStumpers, nil
}
