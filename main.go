package main

import "log"

func main() {
	mode := "FULL"
	numberOfPages := 73 // TODO: make this dynamic
	var jeopardyBoxScores []JeopardyGameBoxScore

	if mode == "FULL" {
		// create the DB if it DNE
		err := createDatabaseIfDoesNotExist()
		if err != nil {
			log.Fatalf("failed to create the DB: %v", err)
			return
		}

		err = createJeopardyGameBoxScoreTable()
		if err != nil {
			log.Fatalf("failed to create the table: %v", err)
			return
		}

		jeopardyBoxScores = ScrapeJeopardataFull(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			saveJeopardyGameBoxScore(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel(jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}

	} else if mode == "INCREMENTAL" {
		jeopardyBoxScores = ScrapeJeopardataIncremental(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			saveJeopardyGameBoxScore(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel(jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}
	} else {
		log.Fatalf("invalid mode: %s", mode)
		return
	}
}
