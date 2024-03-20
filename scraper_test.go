package main

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestGetNumberOfTripleStumpers(t *testing.T) {
	episodeId := "ep-9151"

	// Parse the HTML string into an *html.Node
	htmlNode, err := html.Parse(strings.NewReader(htmlForTests))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(htmlNode)

	numberOfTripleStrumpersFound, err := getNumberOfTripleStumpers(doc, episodeId)
	assert.NoError(t, err)
	assert.Equal(t, 4, numberOfTripleStrumpersFound)
}

func TestGetGameTotals(t *testing.T) {
	episodeId := "ep-9151"

	// Parse the HTML string into an *html.Node
	htmlNode, err := html.Parse(strings.NewReader(htmlForTests))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(htmlNode)

	gameTotals, err := getGameTotals(doc, episodeId)
	assert.NoError(t, err)

	// Test for Yogesh
	assert.Equal(t, "RAUT", gameTotals[0].LastName)
	assert.Equal(t, true, gameTotals[0].GameWinner) // TODO: FIX THIS
	assert.Equal(t, 46, gameTotals[0].TotalAttempts)
	assert.Equal(t, 21, gameTotals[0].TotalBuzzes)
	assert.Equal(t, 46, gameTotals[0].TotalBuzzPercentage)
	assert.Equal(t, 21, gameTotals[0].TotalCorrect)
	assert.Equal(t, 0, gameTotals[0].TotalIncorrect)
	assert.Equal(t, 0, gameTotals[0].TotalIncorrect)
	assert.Equal(t, 100, gameTotals[0].CorrectPercentage)
	assert.Equal(t, 0, gameTotals[0].TotalDailyDoubleCorrect)
	assert.Equal(t, 0, gameTotals[0].TotalDailyDoubleIncorrect)
	assert.Equal(t, 0, gameTotals[0].TotalDailyDoubleWinnings)
	assert.Equal(t, 13399, gameTotals[0].FinalScore)

	// Test for Troy
	assert.Equal(t, "MEYER", gameTotals[1].LastName)
	assert.Equal(t, false, gameTotals[1].GameWinner) // TODO: FIX THIS
	assert.Equal(t, 47, gameTotals[1].TotalAttempts)
	assert.Equal(t, 18, gameTotals[1].TotalBuzzes)
	assert.Equal(t, 38, gameTotals[1].TotalBuzzPercentage)
	assert.Equal(t, 19, gameTotals[1].TotalCorrect)
	assert.Equal(t, 0, gameTotals[1].TotalIncorrect)
	assert.Equal(t, 0, gameTotals[1].TotalIncorrect)
	assert.Equal(t, 100, gameTotals[1].CorrectPercentage)
	assert.Equal(t, 1, gameTotals[1].TotalDailyDoubleCorrect)
	assert.Equal(t, 0, gameTotals[1].TotalDailyDoubleIncorrect)
	assert.Equal(t, 2800, gameTotals[1].TotalDailyDoubleWinnings)
	assert.Equal(t, 6399, gameTotals[1].FinalScore)
}
