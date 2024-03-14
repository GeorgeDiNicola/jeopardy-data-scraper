package main

import (
	"gorm.io/gorm"
)

type JeopardyGameBoxScore struct {
	gorm.Model
	EpisodeNumber string `gorm:"type:varchar(100)" json:"episode_number"`
	EpisodeTitle  string `gorm:"type:varchar(100)" json:"episode_title"`
	Date          string `gorm:"type:date" json:"date"`
	LastName      string `gorm:"type:varchar(100)" json:"last_name"`
	FirstName     string `gorm:"type:varchar(100)" json:"first_name"`
	HomeCity      string `gorm:"type:varchar(100)" json:"city"`
	HomeState     string `gorm:"type:varchar(100)" json:"state"`
	GameWinner    bool   `json:"game_winner"` // true or false
	//JeopardyRound          string  `gorm:"type:varchar(50)" json:"jeopardy_round"` // Enum: jeopardy_round, double_jeopardy, final_jeopardy, game_totals
	R1Att                  int `json:"r1_att"`
	R1Buz                  int `json:"r1_buz"`
	R1BuzPercentage        int `gorm:"type:decimal(5,2)" json:"r1_buz_percentage"`
	R1Correct              int `json:"r1_correct"`
	R1Incorrect            int `json:"r1_incorrect"`
	R1CorrectPercentage    int `gorm:"type:decimal(5,2)" json:"r1_correct_percentage"`
	R1Dd1                  int `json:"r1_dd1"`
	R1Dd2                  int `json:"r1_dd2"`
	R1Eor                  int `json:"r1_eor"`
	R2Att                  int `json:"r2_att"`
	R2Buz                  int `json:"r2_buz"`
	R2BuzPercentage        int `gorm:"type:decimal(5,2)" json:"r2_buz_percentage"`
	R2Correct              int `json:"r2_correct"`
	R2Incorrect            int `json:"r2_incorrect"`
	R2CorrectPercentage    int `gorm:"type:decimal(5,2)" json:"r2_correct_percentage"`
	R2Dd1                  int `json:"r2_dd1"`
	R2Dd2                  int `json:"r2_dd2"`
	R2Eor                  int `json:"r2_eor"`
	StartingFjScore        int `json:"starting_fj_score"`
	FjWager                int `json:"fj_wager"`
	FinalScore             int `json:"final_score"`
	AttTotal               int `json:"att_total"`
	BuzTotal               int `json:"buz_total"`
	BuzPercentageTotal     int `gorm:"type:decimal(5,2)" json:"buz_percentage_total"`
	CorrectTotal           int `json:"correct_total"`
	IncorrectTotal         int `json:"incorrect_total"`
	CorrectPercentageTotal int `gorm:"type:decimal(5,2)" json:"correct_percentage_total"`
	DdTotal                int `json:"dd_total"`
	DdPercentageTotal      int `gorm:"type:decimal(5,2)" json:"dd_percentage_total"`
	TotalFjCorrect         int `json:"total_fj_correct"`
	TotalFjIncorrect       int `json:"total_fj_incorrect"`
	TotalFjPercentage      int `gorm:"type:decimal(5,2)" json:"total_fj_percentage"`
	CoryatScore            int `json:"coryat_score"`
}

type JeopardyGameBoxScoreTotal struct {
	gorm.Model
	EpisodeNumber       string `gorm:"type:varchar(100)" json:"episode_number"`
	EpisodeTitle        string `gorm:"type:varchar(100)" json:"episode_title"`
	Date                string `gorm:"type:date" json:"date"`
	LastName            string `gorm:"type:varchar(100)" json:"last_name"`
	FirstName           string `gorm:"type:varchar(100)" json:"first_name"`
	City                string `gorm:"type:varchar(100)" json:"city"`
	State               string `gorm:"type:varchar(100)" json:"state"`
	GameWinner          bool   `json:"game_winner"` // true or false
	TotalAtt            int    `json:"total_att"`
	TotalBuz            int    `json:"total_buz"`
	TotalBuzPercentage  int    `gorm:"type:decimal(5,2)" json:"total_buz_percentage"`
	TotalCorrect        int    `json:"total_correct"`
	TotalIncorrect      int    `json:"total_incorrect"`
	CorrectPercentage   int    `gorm:"type:decimal(5,2)" json:"correct_percentage"`
	TotalDdCorrect      int    `json:"total_dd_correct"`
	TotalDdIncorrect    int    `json:"total_dd_incorrect"`
	TotalDdWinnings     int    `gorm:"type:decimal(5,2)" json:"total_dd_winnings"`
	FinalScore          int    `json:"final_score"`
	TotalTripleStumpers int    `json:"total_triple_stumpers"`
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