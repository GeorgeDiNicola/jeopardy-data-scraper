# jeopardy-data-scraper
Web scraper for parsing, transforming, and storing data from https://www.jeopardy.com/track/jeopardata.

## Usage Modes
1. Full data extract and storage
2. Incremental data extract and storage

## Design Goals
- Provide an API or database for data analyst Jeopardy fans to more easily explore Jeopardata. Making it into a better format
- Analyze/Identify patterns in Jeopardy gameplay over time and identify inflection points from influential players like James Holzhauer
- Provide a location to store historic Jeopardy! data in a more data analytical-friendly format than the website

## System Design
![Jeopardata System Design](/Users/georgedinicola/jeopardy-data-scraper/docs/Jeopardy-System-Design.png)

## Algorithm for Bulk Scrape
1. For each web page
2. Collect each episode by DATE
3. For each date
	4. Collect the names & home city info of the contestants
	5. Collect Jeopardy Round Data
	6. Collect Double Jeopardy Roud Data
	7. Collect Final Jeopardy Round Data
	8. Collect game totals data
		9. write to DB

## Algorithm for Incremental Scrape
1. Check the last date in the DB
2. Collect data each episode by DATE from last date until current
3. For each date
	4. Collect the names & home city info of the contestants
5. Collect Jeopardy Round Data
6. Collect Double Jeopardy Roud Data
7. Collect Final Jeopardy Round Data
8. Collect game totals data
	9. write to DB

## Definitions:
- EOR - Score at the end of the round
- ATT - Attempts to buzz in
- BUZ - number of times a contestant buzzed in
- BUZ % -  percentage individual contestant buzzed in vs. attempts
- COR/INC - How many correct/incorrect responses
- CORRECT % - percentage of correct responses
- DD - Daily Double/FJ Final Jeopardy!
- Triple Stumper - Clues (except DD) for which none of the contestants provide a correct response

## Jeopardata API Endpoints

### Basic Information
- List all episode numbers, dates, and titles: /v1/jeopardata/episodes
- Get an array of data from a specific episode by episode number or date: /v1/jeopardata/episodes/{<date>}|{<episode-number>}

### Contestant Info
- Get an array of contestants from an episode number or date:  /v1/jeopardata/episodes/{<date>}|{episodeNumber}/contestants
- Get a specific contestant's information/metadata from an episode number or date: /v1/jeopardata/episodes/{<date>}|{episodeNumber}/contestants/{last-name}

### Game Rounds and Scores
- /v1/jeopardata/episodes/{<date>}|{episodeNumber}/jeopardyRoundScores
- /v1/jeopardata/episodes/{<date>}|{episodeNumber}/doubleJeopardyRoundScores
- /v1/jeopardata/episodes/{<date>}|{episodeNumber}/finalJeopardyScores
- /v1/jeopardata//episodes/{<date>}|{episodeNumber}/gameTotals
- /v1/jeopardata//episodes/{<date>}|{episodeNumber}/tripleStumpers

### Get Game Notes for an Episode
- /v1/jeopardata/episodes/{<date>}|{episodeNumber}/notes