package main

import "log"

func main() {
	mode := "FULL"
	numberOfPages := 73 // TODO: make this dynamic
	var jeopardyBoxScores []JeopardyGameBoxScore

	if mode == "FULL" {
		jeopardyBoxScores = ScrapeAllJeopardata(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			writeToPostgresDB(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel(jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}

	} else if mode == "INCREMENTAL" {
		jeopardyBoxScores = ScrapeIncrementalJeopardata(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			writeToPostgresDB(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel(jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}
	}
}
