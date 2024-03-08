package main

type ContestantStats struct {
	Attempts    int    `json:"attempts"`
	Buzzes      int    `json:"buzzes"`
	BuzzPct     string `json:"buzz_pct"`
	Correct     int    `json:"correct"`
	Incorrect   int    `json:"incorrect"`
	CorrectPct  string `json:"correct_pct"`
	DDCorrect   int    `json:"dd_correct,omitempty"`
	DDIncorrect int    `json:"dd_incorrect,omitempty"`
	DDWager     int    `json:"dd_wager,omitempty"`
	EORScore    string `json:"eor_score"`
	// TODO: Champion?
	Winner bool `json:"champion"`
}

type RoundStats struct {
	JeopardyRound       ContestantStats `json:"jeopardy_round"`
	DoubleJeopardyRound ContestantStats `json:"double_jeopardy_round"`
	FinalJeopardyRound  ContestantStats `json:"final_jeopardy_round,omitempty"`
}

type Contestant struct {
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	HomeCity   string     `json:"home_city"`
	HomeState  string     `json:"home_state"`
	RoundStats RoundStats `json:"round_stats"`
}

type JeopardyGame struct {
	Date           string       `json:"date"`
	Contestants    []Contestant `json:"contestants"`
	ChampionName   string       `json:"champion_name"`
	TripleStumpers int          `json:"triple_stumpers"`
	TournamentGame bool         `json:"tournament_game"`
}
