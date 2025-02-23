package pipeline

import (
	"fmt"

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
func InternetSearch() {}

func scrape() {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	// TODO this string comes from the duckduckgo query.
	c.Visit("http://example.com")
}
