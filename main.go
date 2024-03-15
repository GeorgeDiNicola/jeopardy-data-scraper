package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	var allEpisodeJeopardyGameBoxScores []JeopardyGameBoxScore

	// TODO: make this dynamic
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

			episodeDate := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .date", episode.EpisodeID)).Text()
			episodeTitle := doc.Find(fmt.Sprintf("table[aria-labelledby='%s-label'] .title", episode.EpisodeID)).Text()
			episodeNumber := strings.Split(episode.EpisodeID, "-")[1]

			for i := 0; i < len(contestants); i++ {
				var jeopardyGameBoxScore JeopardyGameBoxScore
				jeopardyGameBoxScore.EpisodeNumber = episodeNumber
				jeopardyGameBoxScore.EpisodeTitle = episodeTitle
				jeopardyGameBoxScore.Date = episodeDate
				jeopardyGameBoxScore.LastName = contestants[i].LastName
				jeopardyGameBoxScore.FirstName = contestants[i].FirstName
				jeopardyGameBoxScore.HomeCity = contestants[i].HomeCity
				jeopardyGameBoxScore.HomeState = contestants[i].HomeState
				jeopardyGameBoxScore.GameWinner = contestants[i].GameWinner

				// fill in R1
				jeopardyGameBoxScore.R1Att = jeopardyRounds[i].Att
				jeopardyGameBoxScore.R1Buz = jeopardyRounds[i].Buz
				jeopardyGameBoxScore.R1BuzPercentage = jeopardyRounds[i].BuzPercentage
				jeopardyGameBoxScore.R1Correct = jeopardyRounds[i].Correct
				jeopardyGameBoxScore.R1Incorrect = jeopardyRounds[i].Incorrect
				jeopardyGameBoxScore.R1CorrectPercentage = jeopardyRounds[i].CorrectPercentage
				jeopardyGameBoxScore.R1DailyDouble = jeopardyRounds[i].DailyDouble
				jeopardyGameBoxScore.R1Eor = jeopardyRounds[i].EorScore

				// fill in R2
				jeopardyGameBoxScore.R2Att = doubleJeopardyRounds[i].Att
				jeopardyGameBoxScore.R2Buz = doubleJeopardyRounds[i].Buz
				jeopardyGameBoxScore.R2BuzPercentage = doubleJeopardyRounds[i].BuzPercentage
				jeopardyGameBoxScore.R2Correct = doubleJeopardyRounds[i].Correct
				jeopardyGameBoxScore.R2Incorrect = doubleJeopardyRounds[i].Incorrect
				jeopardyGameBoxScore.R2CorrectPercentage = doubleJeopardyRounds[i].CorrectPercentage
				jeopardyGameBoxScore.R2DailyDouble1 = doubleJeopardyRounds[i].DailyDouble1
				jeopardyGameBoxScore.R2DailyDouble2 = doubleJeopardyRounds[i].DailyDouble2
				jeopardyGameBoxScore.R2Eor = doubleJeopardyRounds[i].EorScore

				// fil in r3
				jeopardyGameBoxScore.StartingFjScore = finalJeopardyRounds[i].StartingFjScore
				jeopardyGameBoxScore.FjWager = finalJeopardyRounds[i].FjWager
				jeopardyGameBoxScore.FjFinalScore = finalJeopardyRounds[i].FinalScore

				// totals
				jeopardyGameBoxScore.AttTotal = boxScoreTotals[i].TotalAtt
				jeopardyGameBoxScore.BuzTotal = boxScoreTotals[i].TotalBuz
				jeopardyGameBoxScore.BuzPercentageTotal = boxScoreTotals[i].TotalBuzPercentage
				jeopardyGameBoxScore.CorrectTotal = boxScoreTotals[i].TotalCorrect
				jeopardyGameBoxScore.IncorrectTotal = boxScoreTotals[i].TotalIncorrect
				jeopardyGameBoxScore.CorrectPercentageTotal = boxScoreTotals[i].CorrectPercentage
				jeopardyGameBoxScore.DailyDoubleCorrectTotal = boxScoreTotals[i].TotalDailyDoubleCorrect
				jeopardyGameBoxScore.DailyDoubleIncorrectTotal = boxScoreTotals[i].TotalDailyDoubleIncorrect
				jeopardyGameBoxScore.DailyDoubleWinningsTotal = boxScoreTotals[i].TotalDailyDoubleWinnings
				jeopardyGameBoxScore.FinalScoreTotal = boxScoreTotals[i].FinalScore

				allEpisodeJeopardyGameBoxScores = append(allEpisodeJeopardyGameBoxScores, jeopardyGameBoxScore)
			}
		}
	}

	writeBoxScoreHistoryToExcel(allEpisodeJeopardyGameBoxScores)

	writeToPostgresDB(allEpisodeJeopardyGameBoxScores)
}
