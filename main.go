package main

import "log"

func main() {
	mode := "FULL"      // TODO: make this an input to the main function
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

		jeopardyBoxScores = ScrapeGameDataFull(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			saveJeopardyGameBoxScore(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel("jeopardata_box_scores_sample.xlsx", jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}

	} else if mode == "INCREMENTAL" {
		jeopardyBoxScores = ScrapeGameDataIncremental(numberOfPages)

		if len(jeopardyBoxScores) > 0 {
			saveJeopardyGameBoxScore(jeopardyBoxScores)
			writeBoxScoreHistoryToExcel("jeopardata_box_scores_sample.xlsx", jeopardyBoxScores)
		} else {
			log.Println("No new jeopardata records to extract")
		}
	} else {
		log.Fatalf("invalid mode: %s", mode)
		return
	}
}
