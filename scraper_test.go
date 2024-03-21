package main

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/html"
)

type ScraperSuite struct {
	suite.Suite
	DefaultDoc       *goquery.Document
	DefaultEpisodeId string
}

func (suite *ScraperSuite) SetupTest() {
	htmlNode, err := html.Parse(strings.NewReader(HtmlForTests))
	if err != nil {
		panic(err)
	}

	suite.DefaultDoc = goquery.NewDocumentFromNode(htmlNode)
	suite.DefaultEpisodeId = "ep-9151"
}

func (s *ScraperSuite) TestGetNumberOfTripleStumpers() {
	numberOfTripleStrumpersFound, err := getNumberOfTripleStumpers(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)
	s.Equal(4, numberOfTripleStrumpersFound)
}

func (s *ScraperSuite) TestGetGameTotals() {
	gameTotals, err := getGameTotals(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)

	// Test for Yogesh
	s.Equal("RAUT", gameTotals[0].LastName)
	s.Equal(true, gameTotals[0].GameWinner) // TODO: FIX THIS
	s.Equal(46, gameTotals[0].TotalAttempts)
	s.Equal(21, gameTotals[0].TotalBuzzes)
	s.Equal(46, gameTotals[0].TotalBuzzPercentage)
	s.Equal(21, gameTotals[0].TotalCorrect)
	s.Equal(0, gameTotals[0].TotalIncorrect)
	s.Equal(0, gameTotals[0].TotalIncorrect)
	s.Equal(100, gameTotals[0].CorrectPercentage)
	s.Equal(0, gameTotals[0].TotalDailyDoubleCorrect)
	s.Equal(0, gameTotals[0].TotalDailyDoubleIncorrect)
	s.Equal(0, gameTotals[0].TotalDailyDoubleWinnings)
	s.Equal(13399, gameTotals[0].FinalScore)

	// Test for Troy
	s.Equal("MEYER", gameTotals[1].LastName)
	s.Equal(false, gameTotals[1].GameWinner) // TODO: FIX THIS
	s.Equal(47, gameTotals[1].TotalAttempts)
	s.Equal(18, gameTotals[1].TotalBuzzes)
	s.Equal(38, gameTotals[1].TotalBuzzPercentage)
	s.Equal(19, gameTotals[1].TotalCorrect)
	s.Equal(0, gameTotals[1].TotalIncorrect)
	s.Equal(0, gameTotals[1].TotalIncorrect)
	s.Equal(100, gameTotals[1].CorrectPercentage)
	s.Equal(1, gameTotals[1].TotalDailyDoubleCorrect)
	s.Equal(0, gameTotals[1].TotalDailyDoubleIncorrect)
	s.Equal(2800, gameTotals[1].TotalDailyDoubleWinnings)
	s.Equal(6399, gameTotals[1].FinalScore)
}

func (s *ScraperSuite) TestGetContestantInformation() {
	contestantInfo, err := getContestantInformation(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)

	// Test for Yogesh Raut
	s.Equal("RAUT", contestantInfo[0].LastName)
	s.Equal("YOGESH", contestantInfo[0].FirstName)
	s.Equal("VANCOUVER", contestantInfo[0].HomeCity)
	s.Equal("WASHINGTON", contestantInfo[0].HomeState)
	s.Equal(true, contestantInfo[0].GameWinner)

	// Test for Troy Meyer
	s.Equal("MEYER", contestantInfo[1].LastName)
	s.Equal("TROY", contestantInfo[1].FirstName)
	s.Equal("TAMPA", contestantInfo[1].HomeCity)
	s.Equal("FLORIDA", contestantInfo[1].HomeState)
	s.Equal(false, contestantInfo[1].GameWinner)

	// Test for Ben Chan
	s.Equal("CHAN", contestantInfo[2].LastName)
	s.Equal("BEN", contestantInfo[2].FirstName)
	s.Equal("GREEN BAY", contestantInfo[2].HomeCity)
	s.Equal("WISCONSIN", contestantInfo[2].HomeState)
	s.Equal(false, contestantInfo[2].GameWinner)
}

func (s *ScraperSuite) TestGetJeopardyRound() {
	jeopardyRound, err := getJeopardyRound(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)

	// Test for Yogesh Raut
	s.Equal("RAUT", jeopardyRound[0].LastName)
	s.Equal("YOGESH", jeopardyRound[0].FirstName)
	s.Equal(21, jeopardyRound[0].Attempts)
	s.Equal(10, jeopardyRound[0].Buzzes)
	s.Equal(48, jeopardyRound[0].BuzzPercentage)
	s.Equal(10, jeopardyRound[0].Correct)
	s.Equal(0, jeopardyRound[0].Incorrect)
	s.Equal(100, jeopardyRound[0].CorrectPercentage)
	s.Equal(0, jeopardyRound[0].DailyDouble)
	s.Equal(5400, jeopardyRound[0].EndOfRoundScore)

	// Test for Troy Meyer
	s.Equal("MEYER", jeopardyRound[1].LastName)
	s.Equal("TROY", jeopardyRound[1].FirstName)
	s.Equal(27, jeopardyRound[1].Attempts)
	s.Equal(11, jeopardyRound[1].Buzzes)
	s.Equal(41, jeopardyRound[1].BuzzPercentage)
	s.Equal(12, jeopardyRound[1].Correct)
	s.Equal(0, jeopardyRound[1].Incorrect)
	s.Equal(100, jeopardyRound[1].CorrectPercentage)
	s.Equal(2800, jeopardyRound[1].DailyDouble)
	s.Equal(9400, jeopardyRound[1].EndOfRoundScore)

	// Test for Ben Chan
	s.Equal("CHAN", jeopardyRound[2].LastName)
	s.Equal("BEN", jeopardyRound[2].FirstName)
	s.Equal(22, jeopardyRound[2].Attempts)
	s.Equal(8, jeopardyRound[2].Buzzes)
	s.Equal(36, jeopardyRound[2].BuzzPercentage)
	s.Equal(7, jeopardyRound[2].Correct)
	s.Equal(1, jeopardyRound[2].Incorrect)
	s.Equal(88, jeopardyRound[2].CorrectPercentage)
	s.Equal(0, jeopardyRound[2].DailyDouble)
	s.Equal(3400, jeopardyRound[2].EndOfRoundScore)
}

func (s *ScraperSuite) TestGetDoubleJeopardyRound() {
	doubleJeopardyRound, err := getDoubleJeopardyRound(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)

	// Test for Yogesh Raut
	s.Equal("RAUT", doubleJeopardyRound[0].LastName)
	s.Equal("YOGESH", doubleJeopardyRound[0].FirstName)
	s.Equal(25, doubleJeopardyRound[0].Attempts)
	s.Equal(11, doubleJeopardyRound[0].Buzzes)
	s.Equal(44, doubleJeopardyRound[0].BuzzPercentage)
	s.Equal(11, doubleJeopardyRound[0].Correct)
	s.Equal(0, doubleJeopardyRound[0].Incorrect)
	s.Equal(100, doubleJeopardyRound[0].CorrectPercentage)
	s.Equal(0, doubleJeopardyRound[0].DailyDouble1)
	s.Equal(0, doubleJeopardyRound[0].DailyDouble2)
	s.Equal(16600, doubleJeopardyRound[0].EndOfRoundScore)

	// Test for Troy Meyer
	s.Equal("MEYER", doubleJeopardyRound[1].LastName)
	s.Equal("TROY", doubleJeopardyRound[1].FirstName)
	s.Equal(20, doubleJeopardyRound[1].Attempts)
	s.Equal(7, doubleJeopardyRound[1].Buzzes)
	s.Equal(35, doubleJeopardyRound[1].BuzzPercentage)
	s.Equal(7, doubleJeopardyRound[1].Correct)
	s.Equal(0, doubleJeopardyRound[1].Incorrect)
	s.Equal(100, doubleJeopardyRound[1].CorrectPercentage)
	s.Equal(0, doubleJeopardyRound[1].DailyDouble1)
	s.Equal(0, doubleJeopardyRound[1].DailyDouble2)
	s.Equal(19800, doubleJeopardyRound[1].EndOfRoundScore)

	// Test for Ben Chan
	s.Equal("CHAN", doubleJeopardyRound[2].LastName)
	s.Equal("BEN", doubleJeopardyRound[2].FirstName)
	s.Equal(13, doubleJeopardyRound[2].Attempts)
	s.Equal(8, doubleJeopardyRound[2].Buzzes)
	s.Equal(62, doubleJeopardyRound[2].BuzzPercentage)
	s.Equal(8, doubleJeopardyRound[2].Correct)
	s.Equal(2, doubleJeopardyRound[2].Incorrect)
	s.Equal(80, doubleJeopardyRound[2].CorrectPercentage)
	s.Equal(4200, doubleJeopardyRound[2].DailyDouble1)
	s.Equal(-9600, doubleJeopardyRound[2].DailyDouble2)
	s.Equal(3200, doubleJeopardyRound[2].EndOfRoundScore)
}

func (s *ScraperSuite) TestGetFinalJeopardyRound() {
	finalJeopardyRound, err := getFinalJeopardyRound(s.DefaultDoc, s.DefaultEpisodeId)
	s.NoError(err)

	// Test for Yogesh Raut
	s.Equal("RAUT", finalJeopardyRound[0].LastName)
	s.Equal("YOGESH", finalJeopardyRound[0].FirstName)
	s.Equal(16600, finalJeopardyRound[0].StartingFjScore)
	s.Equal(-3201, finalJeopardyRound[0].FjWager)
	s.Equal(13399, finalJeopardyRound[0].FinalScore)

	// Test for Troy Meyer
	s.Equal("MEYER", finalJeopardyRound[1].LastName)
	s.Equal("TROY", finalJeopardyRound[1].FirstName)
	s.Equal(19800, finalJeopardyRound[1].StartingFjScore)
	s.Equal(-13401, finalJeopardyRound[1].FjWager)
	s.Equal(6399, finalJeopardyRound[1].FinalScore)

	// Test for Ben Chan
	s.Equal("CHAN", finalJeopardyRound[2].LastName)
	s.Equal("BEN", finalJeopardyRound[2].FirstName)
	s.Equal(3200, finalJeopardyRound[2].StartingFjScore)
	s.Equal(3200, finalJeopardyRound[2].FjWager)
	s.Equal(6400, finalJeopardyRound[2].FinalScore)

}

func (s *ScraperSuite) TestGetJeopardataWebPage() {
	doc, err := getJeopardataWebPage(1)
	s.NoError(err)
	s.NotNil(doc)
}

func TestScraperSuite(t *testing.T) {
	suite.Run(t, new(ScraperSuite))
}
