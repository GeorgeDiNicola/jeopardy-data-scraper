# jeopardy-data-scraper
Web scraper for parsing, transforming, and storing data from https://www.jeopardy.com/track/jeopardata.

It also powers this [Tableau Public Dashboard](https://public.tableau.com/app/profile/george.dinicola/viz/JeopardyStatistics).

## Usage Modes
- Full data extract and storage
- Incremental data extract and storage

## Design Goals
- Provide an API and file export mechanisms for data analyst Jeopardy fans to easily explore Jeopardata.
- Analyze and identify patterns in Jeopardy gameplay over time, such as presenting inflection points in data from influential players like James Holzhauer.
- Provide a location to store historic Jeopardy! data in a transformed format more suitable for data-driven analysis and applications.

## System Design
![Jeopardata System Design](docs/Jeopardy-System-Design.png)
*Application scheduled to execute once per day after the Jeopardata posts to extract the most recent Jeopardy episode data it does not know of, then saves it to the DB.

## Data Persistence Layer
- Database Product: **PostgreSQL**
- Database Name: **jeaopardata**
- Table Name: **jeopardy_game_box_scores**
- Description: **The historical box scores for each game, player in each Jeopardy game.**

### Data Dictionary:
| Field                               | Type    | PostgreSQL Type        | JSON Key                          | Description |
|-------------------------------------|---------|------------------------|-----------------------------------|-------------|
| EpisodeNumber                       | string  | `varchar(100)`         | `episode_number`                  | -           |
| EpisodeTitle                        | string  | `varchar(100)`         | `episode_title`                   | -           |
| EpisodeDate                         | string  | `date`                 | `episode_date`                    | -           |
| ContestantLastName                  | string  | `varchar(100)`         | `contestant_last_name`            | -           |
| ContestantFirstName                 | string  | `varchar(100)`         | `contestant_first_name`           | -           |
| HomeCity                            | string  | `varchar(100)`         | `home_city`                       | -           |
| HomeState                           | string  | `varchar(100)`         | `home_state`                      | -           |
| IsWinner                            | bool    | -                      | `is_winner`                       | -           |
| JeopardyAttempts                    | int     | -                      | `jeopardy_attempts`               | -           |
| JeopardyBuzzes                      | int     | -                      | `jeopardy_buzzes`                 | -           |
| JeopardyBuzzPercentage              | int     | -                      | `jeopardy_buzz_percentage`        | -           |
| JeopardyCorrectAnswers              | int     | -                      | `jeopardy_correct_answers`        | -           |
| JeopardyIncorrectAnswers            | int     | -                      | `jeopardy_incorrect_answers`      | -           |
| JeopardyCorrectAnswerPercentage     | int     | -                      | `jeopardy_correct_answer_percentage` | -     |
| JeopardyDailyDoubles           | int     | -                      | `jeopardy_daily_doubles`    | -           |
| JeopardyScore                       | int     | -                      | `jeopardy_score`                  | -           |
| DoubleJeopardyAttempts              | int     | -                      | `double_jeopardy_attempts`        | -           |
| DoubleJeopardyBuzzes                | int     | -                      | `double_jeopardy_buzzes`          | -           |
| DoubleJeopardyBuzzPercentage        | int     | -                      | `double_jeopardy_buzz_percentage` | -           |
| DoubleJeopardyCorrectAnswers        | int     | -                      | `double_jeopardy_correct_answers` | -           |
| DoubleJeopardyIncorrectAnswers      | int     | -                      | `double_jeopardy_incorrect_answers` | -       |
| DoubleJeopardyCorrectAnswerPercentage | int  | -                      | `double_jeopardy_correct_answer_percentage` | - |
| DoubleJeopardyDailyDouble1     | int     | -                      | `double_jeopardy_daily_double_1` | -  |
| DoubleJeopardyDailyDouble2     | int     | -                      | `double_jeopardy_daily_double_2` | -  |
| DoubleJeopardyScore                 | int     | -                      | `double_jeopardy_score`          



## Tableau Dashboard Project - Visualization of Jeopardy game data trends over time
### Project Link: [Tableau Public Dashboard](https://public.tableau.com/app/profile/george.dinicola/viz/JeopardyStatistics)
### Design: 
![Tableau Dashboard Project Design](docs/Tableau-Dashboard-Design.png)

## API Endpoints
| Operation | Endpoint                           | Description |
|-----------|------------------------------------|-------------|
| GET       | `/v1/episodes`                     | Basic Information. Retrieves a list of episodes, optionally filtered by date range or specific attributes like episode number.<br>**Params:**<ul><li>`startDate` (optional)</li><li>`endDate` (optional)</li><li>`episodeNumber` (optional)</li></ul> |
| GET       | `/v1/episodes/{episodeNumber}`     | Game Episode Information. Retrieves detailed information about a specific episode, including contestant details and scores.<br>**Params:**<ul><li>`episodeNumber`</li></ul> |
| GET       | `/v1/episodes/{episodeNumber}/performance` | Game Episode Information. Retrieves information about contestants, potentially filtered by name, home city, or state.<br>**Params:**<ul><li>`episodeNumber`</li><li>`gameWinner` - filters for game winners (i.e. game champion stats only)</li></ul> |
| GET       | `/v1/contestants`                  | Contestant Information. Retrieves information about contestants, potentially filtered by name, home city, or state.<br>**Params:**<ul><li>`lastName` (optional)</li><li>`firstName` (optional)</li><li>`homeCity` (optional)</li><li>`homeState` (optional)</li></ul> |
| GET       | `/v1/export`                       | Data Export. Exports all of the data to the user's web browser.<br>**Params:**<ul><li>`fileType` (default: csv)</li></ul>Supported data types: CSV, XLSX, JSON, Google Sheets |


## Algorithms for Scraping the Jeopardy Web Data
### Algorithm for Bulk Scrape
1. For each web page
2. Collect each episode by DATE
3. For each date
<br>&nbsp; &nbsp;4. Collect the names & home city info of the contestants
<br>&nbsp; &nbsp;5. Collect Jeopardy Round Data
<br>&nbsp; &nbsp;6. Collect Double Jeopardy <br>&nbsp; &nbsp;7. Collect Final Jeopardy Round Data
<br>&nbsp; &nbsp;8. Collect game totals data
<br>&nbsp; &nbsp; &nbsp; &nbsp; 9. write to DB

### Algorithm for Incremental Scrape
1. Check the last date in the DB
2. Collect data each episode by DATE from last date until current
3. For each date
<br>&nbsp; &nbsp;4. Collect the names & home city info of the contestants
5. Collect Jeopardy Round Data
6. Collect Double Jeopardy Roud Data
7. Collect Final Jeopardy Round Data
8. Collect game totals data
<br>&nbsp; &nbsp;9. write to DB