package main

import "fmt"

func main() {

	mode := "FULL"
	numberOfPages := 73 // TODO: make this dynamic

	var allEpisodeJeopardyGameBoxScores []JeopardyGameBoxScore
	// full web scrape
	if mode == "FULL" {
		allEpisodeJeopardyGameBoxScores = ScrapeAllJeopardata(numberOfPages)
		writeToPostgresDB(allEpisodeJeopardyGameBoxScores)
		writeBoxScoreHistoryToExcel(allEpisodeJeopardyGameBoxScores)
	} else if mode == "INCREMENTAL" {
		// TODO: ideas
		//  * get more recent DB by date or episode most recently known
		//  * could use today's date and work backwards until the date last known?
		//allEpisodeJeopardyGameBoxScores = IncrementalScrape()
		fmt.Println(mode)
	}
}
