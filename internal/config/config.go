package config

import "time"

// Constants and configurations
var OutputFileName := "./data/jeopardata.xlsx"
var DateFormat string = "January 2, 2006"
var DelayBetweenRequests time.Duration = 1 // seconds
