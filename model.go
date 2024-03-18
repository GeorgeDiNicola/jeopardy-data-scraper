package main

import (
	"gorm.io/gorm"
)

// TODO: give all of these full names
type JeopardyGameBoxScore struct {
	gorm.Model
	EpisodeNumber                         string `gorm:"type:varchar(100)" json:"episode_number"`
	EpisodeTitle                          string `gorm:"type:varchar(100)" json:"episode_title"`
	EpisodeDate                           string `gorm:"type:date" json:"episode_date"`
	ContestantLastName                    string `gorm:"type:varchar(100)" json:"contestant_last_name"`
	ContestantFirstName                   string `gorm:"type:varchar(100)" json:"contestant_first_name"`
	HomeCity                              string `gorm:"type:varchar(100)" json:"home_city"`
	HomeState                             string `gorm:"type:varchar(100)" json:"home_state"`
	IsWinner                              bool   `json:"is_winner"`
	JeopardyAttempts                      int    `json:"jeopardy_attempts"`
	JeopardyBuzzes                        int    `json:"jeopardy_buzzes"`
	JeopardyBuzzPercentage                int    `json:"jeopardy_buzz_percentage"`
	JeopardyCorrectAnswers                int    `json:"jeopardy_correct_answers"`
	JeopardyIncorrectAnswers              int    `json:"jeopardy_incorrect_answers"`
	JeopardyCorrectAnswerPercentage       int    `json:"jeopardy_correct_answer_percentage"`
	JeopardyDailyDoublesFound             int    `json:"jeopardy_daily_doubles_found"`
	JeopardyScore                         int    `json:"jeopardy_score"`
	DoubleJeopardyAttempts                int    `json:"double_jeopardy_attempts"`
	DoubleJeopardyBuzzes                  int    `json:"double_jeopardy_buzzes"`
	DoubleJeopardyBuzzPercentage          int    `json:"double_jeopardy_buzz_percentage"`
	DoubleJeopardyCorrectAnswers          int    `json:"double_jeopardy_correct_answers"`
	DoubleJeopardyIncorrectAnswers        int    `json:"double_jeopardy_incorrect_answers"`
	DoubleJeopardyCorrectAnswerPercentage int    `json:"double_jeopardy_correct_answer_percentage"`
	DoubleJeopardyDailyDouble1Found       int    `json:"double_jeopardy_daily_double_1_found"`
	DoubleJeopardyDailyDouble2Found       int    `json:"double_jeopardy_daily_double_2_found"`
	DoubleJeopardyScore                   int    `json:"double_jeopardy_score"`
	FinalJeopardyStartingScore            int    `json:"final_jeopardy_starting_score"`
	FinalJeopardyWager                    int    `json:"final_jeopardy_wager"`
	FinalJeopardyScore                    int    `json:"final_jeopardy_score"`
	TotalAttempts                         int    `json:"total_attempts"`
	TotalBuzzes                           int    `json:"total_buzzes"`
	TotalBuzzPercentage                   int    `json:"total_buzz_percentage"`
	TotalCorrectAnswers                   int    `json:"total_correct_answers"`
	TotalIncorrectAnswers                 int    `json:"total_incorrect_answers"`
	TotalCorrectAnswerPercentage          int    `json:"total_correct_answer_percentage"`
	TotalDailyDoublesCorrect              int    `json:"total_daily_doubles_correct"`
	TotalDailyDoublesIncorrect            int    `json:"total_daily_doubles_incorrect"`
	TotalDailyDoubleWinnings              int    `json:"total_daily_double_winnings"`
	TotalScore                            int    `json:"total_score"`
	TotalTripleStumpers                   int    `json:"total_triple_stumpers"`
	CoryatScore                           int    `json:"coryat_score"`
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
	Att               int
	Buz               int
	BuzPercentage     int
	Correct           int
	Incorrect         int
	CorrectPercentage int
	DailyDouble       int
	EorScore          int
}

type DoubleJeopardyRound struct {
	LastName          string
	FirstName         string
	Att               int
	Buz               int
	BuzPercentage     int
	Correct           int
	Incorrect         int
	CorrectPercentage int
	DailyDouble1      int
	DailyDouble2      int
	EorScore          int
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
	TotalAtt                  int    `json:"total_att"`
	TotalBuz                  int    `json:"total_buz"`
	TotalBuzPercentage        int    `json:"total_buz_percentage"`
	TotalCorrect              int    `json:"total_correct"`
	TotalIncorrect            int    `json:"total_incorrect"`
	CorrectPercentage         int    `json:"correct_percentage"`
	TotalDailyDoubleCorrect   int    `json:"total_daily_double_correct"`
	TotalDailyDoubleIncorrect int    `json:"total_daily_double_incorrect"`
	TotalDailyDoubleWinnings  int    `json:"total_daily_double_winnings"`
	FinalScore                int    `json:"final_score"`
}
