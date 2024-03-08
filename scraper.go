package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// func getHtmlForForDate(htmlContent, date string) string {
// 	re := regexp.MustCompile(fmt.Sprintf(`(?s)<div class="date">%s</div>(.*?)<div class="date">`, date))
// 	matches := re.FindStringSubmatch(htmlContent)
// 	if len(matches) <= 1 {
// 		fmt.Printf("No HTML for date %s", date)
// 		return ""
// 	}

// 	return matches[1]
// }

func main() {
	response, err := http.Get("https://www.jeopardy.com/track/jeopardata")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Error fetching the page: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	htmlContent := string(body)

	// HTML content between 2 date tags
	//htmlForDate := getHtmlForForDate(htmlContent, "March 6, 2024")

	// TODO: only match for the current date!
	htmlOneGame := htmlContent

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlOneGame))
	if err != nil {
		log.Fatal(err)
	}

	// get the game's contestant information
	doc.Find(".contestant").Each(func(i int, s *goquery.Selection) {
		lastName := strings.TrimSpace(s.Find(".name-1").Text())
		firstName := strings.TrimSpace(s.Find(".name-0").Text())
		home := strings.TrimSpace(s.Find(".home").Text())
		homeCityState := strings.Split(home, ", ")
		homeCity, homeState := strings.TrimSpace(homeCityState[0]), strings.TrimSpace(homeCityState[1])
		fmt.Printf("%s %s from %s, %s", lastName, firstName, homeCity, homeState)
		fmt.Println()
	})

	// get the Jeopardy Round information for each game
	doc.Find(".jeopardy-round").Each(func(index int, round *goquery.Selection) {
		// first row is the header and skipping it
		if index == 0 {
			return
		}

		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()
		name := firstName + " " + lastName

		att := strings.TrimSpace(round.Find("td[data-header='ATT']").Text())
		buz := strings.TrimSpace(round.Find("td[data-header='BUZ']").Text())
		buzPercent := strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text())
		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		correctPercent := strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text())
		dd := strings.TrimSpace(round.Find("td[data-header='DD']").Text())
		eorScore := strings.TrimSpace(round.Find("td[data-header='EOR SCORE']").Text())

		fmt.Printf("J Contestant: %d, %s\n", index, name)
		fmt.Printf("ATT: %s, BUZ: %s, BUZ %%: %s, COR/INC: %s, CORRECT %%: %s, DD: %s, EOR SCORE: %s\n", att, buz, buzPercent, corInc, correctPercent, dd, eorScore)
		fmt.Println("-----")
		index += 1
	})

	// TODO: get the double jeopary round info for each game
	doc.Find(".jeopardy-round").Each(func(index int, round *goquery.Selection) {
		// first row is the header and skipping it
		if index == 0 {
			return
		}

		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()
		name := firstName + " " + lastName

		att := strings.TrimSpace(round.Find("td[data-header='ATT']").Text())
		buz := strings.TrimSpace(round.Find("td[data-header='BUZ']").Text())
		buzPercent := strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text())
		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		correctPercent := strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text())
		dd := strings.TrimSpace(round.Find("td[data-header='DD']").Text())
		eorScore := strings.TrimSpace(round.Find("td[data-header='EOR SCORE']").Text())

		fmt.Printf("DD Contestant: %d, %s\n", index, name)
		fmt.Printf("ATT: %s, BUZ: %s, BUZ %%: %s, COR/INC: %s, CORRECT %%: %s, DD: %s, EOR SCORE: %s\n", att, buz, buzPercent, corInc, correctPercent, dd, eorScore)
		fmt.Println("-----")
		index += 1
	})

	// TODO: get the final jeopary round info for each game

	// get the Game Totals info for each game
	doc.Find(".game-totals").Each(func(index int, round *goquery.Selection) {
		// first row is the header and skipping it
		if index == 0 {
			return
		}

		firstName := round.Find(".name-0").Text()
		lastName := round.Find(".name-1").Text()
		name := firstName + " " + lastName

		att := strings.TrimSpace(round.Find("td[data-header='ATT']").Text())
		buz := strings.TrimSpace(round.Find("td[data-header='BUZ']").Text())
		buzPercent := strings.TrimSpace(round.Find("td[data-header='BUZ %']").Text())
		corInc := strings.TrimSpace(round.Find("td[data-header='COR/INC']").Text())
		correctPercent := strings.TrimSpace(round.Find("td[data-header='CORRECT %']").Text())
		dd := strings.TrimSpace(round.Find("td[data-header='DD (COR/INC)']").Text())
		finalScore := strings.TrimSpace(round.Find("td[data-header='Final Score']").Text())

		fmt.Printf("F Contestant: %d, %s\n", index, name)
		fmt.Printf("ATT: %s, BUZ: %s, BUZ %%: %s, COR/INC: %s, CORRECT %%: %s, DD (COR/INC): %s, Final Score: %s\n", att, buz, buzPercent, corInc, correctPercent, dd, finalScore)
		fmt.Println("-----")
		index += 1
	})
}
