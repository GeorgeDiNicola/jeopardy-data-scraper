package main

import (
	"gorm.io/gorm"
)

type JeopardyGameBoxScore struct {
	gorm.Model
	EpisodeNumber                  string `gorm:"type:varchar(100)" json:"episode_number"`
	EpisodeTitle                   string `gorm:"type:varchar(100)" json:"episode_title"`
	EpisodeDate                    string `gorm:"type:date" json:"episode_date"`
	ContestantLastName             string `gorm:"type:varchar(100)" json:"contestant_last_name"`
	ContestantFirstName            string `gorm:"type:varchar(100)" json:"contestant_first_name"`
	HomeCity                       string `gorm:"type:varchar(100)" json:"home_city"`
	HomeState                      string `gorm:"type:varchar(100)" json:"home_state"`
	IsWinner                       bool   `json:"is_winner"`
	RoundOneAttempts               int    `json:"round_one_attempts"`
	RoundOneBuzzes                 int    `json:"round_one_buzzes"`
	RoundOneBuzzPercent            int    `json:"round_one_buzz_percent"`
	RoundOneCorrectAnswers         int    `json:"round_one_correct_answers"`
	RoundOneIncorrectAnswers       int    `json:"round_one_incorrect_answers"`
	RoundOneCorrectAnswerPercent   int    `json:"round_one_correct_answer_percent"`
	RoundOneDailyDoubles           int    `json:"round_one_daily_doubles"`
	RoundOneScore                  int    `json:"round_one_score"`
	RoundTwoAttempts               int    `json:"round_two_attempts"`
	RoundTwoBuzzes                 int    `json:"round_two_buzzes"`
	RoundTwoBuzzPercent            int    `json:"round_two_buzz_percent"`
	RoundTwoCorrectAnswers         int    `json:"round_two_correct_answers"`
	RoundTwoIncorrectAnswers       int    `json:"round_two_incorrect_answers"`
	RoundTwoCorrectAnswerPercent   int    `json:"round_two_correct_answer_percent"`
	RoundTwoDailyDouble1           int    `json:"round_two_daily_double_1"`
	RoundTwoDailyDouble2           int    `json:"round_two_daily_double_2"`
	RoundTwoScore                  int    `json:"round_two_score"`
	FinalJeopardyStartingScore     int    `json:"final_jeopardy_starting_score"`
	FinalJeopardyWager             int    `json:"final_jeopardy_wager"`
	FinalJeopardyScore             int    `json:"final_jeopardy_score"`
	TotalGameAttempts              int    `json:"total_game_attempts"`
	TotalGameBuzzes                int    `json:"total_game_buzzes"`
	TotalGameBuzzPercent           int    `json:"total_game_buzz_percent"`
	TotalGameCorrectAnswers        int    `json:"total_game_correct_answers"`
	TotalGameIncorrectAnswers      int    `json:"total_game_incorrect_answers"`
	TotalGameCorrectAnswerPercent  int    `json:"total_game_correct_answer_percent"`
	TotalGameDailyDoublesCorrect   int    `json:"total_daily_doubles_correct"`
	TotalGameDailyDoublesIncorrect int    `json:"total_game_daily_doubles_incorrect"`
	TotalGameDailyDoubleWinnings   int    `json:"total_game_daily_double_winnings"`
	TotalGameScore                 int    `json:"total_game_score"`
	TotalTripleStumpers            int    `json:"total_triple_stumpers"`
}

type Episode struct {
	EpisodeID string
	Date      string
}

type Contestant struct {
	FirstName  string
	LastName   string
	HomeCity   string
	HomeState  string
	GameWinner bool
}

type JeopardyRound struct {
	LastName          string
	FirstName         string
	Attempts          int
	Buzzes            int
	BuzzPercentage    int
	Correct           int
	Incorrect         int
	CorrectPercentage int
	DailyDouble       int
	EndOfRoundScore   int
}

type DoubleJeopardyRound struct {
	LastName          string
	FirstName         string
	Attempts          int
	Buzzes            int
	BuzzPercentage    int
	Correct           int
	Incorrect         int
	CorrectPercentage int
	DailyDouble1      int
	DailyDouble2      int
	EndOfRoundScore   int
}

type FinalJeopardyRound struct {
	LastName        string
	FirstName       string
	StartingFjScore int
	FjWager         int
	FinalScore      int
}

type JeopardyGameBoxScoreTotal struct {
	EpisodeNumber             string `json:"episode_number"`
	EpisodeTitle              string `json:"episode_title"`
	Date                      string `json:"date"`
	LastName                  string `json:"last_name"`
	FirstName                 string `json:"first_name"`
	City                      string `json:"city"`
	State                     string `json:"state"`
	GameWinner                bool   `json:"game_winner"` // true or false
	TotalAttempts             int    `json:"total_att"`
	TotalBuzzes               int    `json:"total_buz"`
	TotalBuzzPercentage       int    `json:"total_buz_percentage"`
	TotalCorrect              int    `json:"total_correct"`
	TotalIncorrect            int    `json:"total_incorrect"`
	CorrectPercentage         int    `json:"correct_percentage"`
	TotalDailyDoubleCorrect   int    `json:"total_daily_double_correct"`
	TotalDailyDoubleIncorrect int    `json:"total_daily_double_incorrect"`
	TotalDailyDoubleWinnings  int    `json:"total_daily_double_winnings"`
	FinalScore                int    `json:"final_score"`
}
