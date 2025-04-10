package pipeline

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

const (
	// TODO add ?q={query encoded} and maybe some headers for User-agent
	duckduckgoUrl = "https://html.duckduckgo.com/html/"
)

// InternetSearch is used to search the internet with an models query request
//
// First the model should generate a query.
// Then the tool should check the query to make sure its a valid url.
// Then the search tool should make a request to a privacy oriented search engine.
// Once this is done it should web scrape the top five websites, and store it in the contex box.
// Internet search should not be compared with the rest of tools because it can be
// dangerous if not used properly, hence why it is serperate.
func InternetSearch(query string) {
	c := colly.NewCollector()

	query = strings.ReplaceAll(query, " ", "+")

	// Find and visit all links
	c.OnHTML("result__url", func(e *colly.HTMLElement) {
		//counter used to limit the number of websites
		maxLinks := 0
		for maxLinks < 3 {
			//searches for links on duckduckgo and creates links to individual sites to be scraped
			link := e.Attr("href")
			index := strings.Index(link, "https")
			link = link[index:]
			maxLinks++

			scrape(link)
		}
	})

	err := c.Visit(fmt.Sprintf("https://html.duckduckgo/html/?q=%s", query))
	if err != nil {
		fmt.Println("error while scraping duckduckgo")
	}
}

// used to scrape individual sites
func scrape(link string) {

	c := colly.NewCollector()

	//scrapes all paragraph elements from each webpage
	c.OnHTML("p", func(e *colly.HTMLElement) {
		//concatenates all text from paragraph elements
		text := ""
		text += e.Text
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	//start scraping by visiting the page
	err := c.Visit(link)
	if err != nil {
		fmt.Println("error while scraping webpage")
	}

	//TODO add functionality to embed text
}
