package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()

	// get all dates

	// get all last names
	var lastNames []string
	c.OnHTML(".name-1", func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		lastNames = append(lastNames, e.Text)
	})

	// get all first names
	var firstNames []string
	c.OnHTML(".name-0", func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		firstNames = append(firstNames, e.Text)
	})

	// get all ATTs
	var atts []string
	c.OnHTML(`td[data-header="ATT"]`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		atts = append(atts, e.Text)
	})

	// get all Buz
	var buzs []string
	c.OnHTML(`td[data-header="BUZ"]`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		buzs = append(buzs, e.Text)
	})

	// get all BUZ %
	var buzPercentages []string
	c.OnHTML(`td[data-header="BUZ %"]`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		buzPercentages = append(buzPercentages, e.Text)
	})

	// get all COR/INC
	var corrects []string
	c.OnHTML(`td[data-header="COR/INC"] span.cor`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		corrects = append(corrects, e.Text)
	})

	// get all inc
	var incorrects []string
	c.OnHTML(`td[data-header="COR/INC"] span.inc`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		incorrects = append(incorrects, e.Text)
	})

	// get all CORRECT %
	var correctPercentages []string
	c.OnHTML(`td[data-header="CORRECT %"]`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		correctPercentages = append(correctPercentages, e.Text)
	})

	// get EOR SCORE
	var eorScores []string
	c.OnHTML(`td[data-header="EOR SCORE"]`, func(e *colly.HTMLElement) {
		// e.Text gets the text content of the found element
		eorScores = append(eorScores, e.Text)
	})

	c.Visit("https://www.jeopardy.com/track/jeopardata")

	fmt.Println(lastNames)

	// for i := 0; i < len(lastNames); i++ {
	// 	lastName, firstName, att, buz := lastNames[i], firstNames[i], atts[i], buzs[i]
	// 	buzPercentage, correct, eor := buzPercentages[i], corrects[i], eorScores[i]
	// 	incorrect, correctPercentage := incorrects[i], correctPercentages[i]
	// 	fmt.Printf("%s, %s, %s, %s, %s, %s, %s, %s, %s", lastName, firstName, att, buz, buzPercentage, correct, incorrect, correctPercentage, eor)
	// 	fmt.Println("\n\n")
	// }
}
