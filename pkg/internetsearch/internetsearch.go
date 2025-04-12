package pipeline

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/openai/openai-go"
)

const (
	// TODO add ?q={query encoded} and maybe some headers for User-agent
	duckduckgoUrl = "https://html.duckduckgo.com/html/"
)

type embeddingResponse struct {
	Response openai.CreateEmbeddingResponse
}

// InternetSearch is used to search the internet with an models query request
//
// First the model should generate a query.
// Then the tool should check the query to make sure its a valid url.
// Then the search tool should make a request to a privacy oriented search engine.
// Once this is done it should web scrape the top five websites, and store it in the contex box.
// Internet search should not be compared with the rest of tools because it can be
// dangerous if not used properly, hence why it is serperate.
func InternetSearch(query string) {
	collyCollector := colly.NewCollector()

	query = strings.ReplaceAll(query, " ", "+")

	// Find and visit all links
	collyCollector.OnHTML("result__url", func(element *colly.HTMLElement) {
		//counter used to limit the number of websites
		maxLinks := 0
		for maxLinks < 3 {
			//searches for links on duckduckgo and creates links to individual sites to be scraped
			link := element.Attr("href")
			index := strings.Index(link, "https")
			link = link[index:]
			maxLinks++

			scrape(link)
		}
	})

	err := collyCollector.Visit(fmt.Sprintf("https://html.duckduckgo/html/?q=%s", query))
	if err != nil {
		log.Println("error while scraping duckduckgo")
	}
}

// used to scrape individual sites
func scrape(link string) {

	collyCollector := colly.NewCollector()

	//scrapes all paragraph elements from each webpage
	collyCollector.OnHTML("p", func(element *colly.HTMLElement) {
		//concatenates all text from paragraph elements
		text := element.Text
		embeddingText := strings.NewReader(text)
		resp, err := http.Post("http://localhost:8082/emb/generate", "application/json", embeddingText)
		if err != nil {
			log.Println("Sorry Vito what was that? You cut out.")
		}
		defer resp.Body.Close()

		var embeddingResponse embeddingResponse

		err = json.NewDecoder(resp.Body).Decode(&embeddingResponse)
		if err != nil {
			log.Println("Failed to Read response body in internet search")
		}

		vectorizedText := embeddingResponse.Response

		fmt.Println(vectorizedText)

	})
	collyCollector.OnRequest(func(req *colly.Request) {
		log.Println("Visiting", req.URL)
	})

	//start scraping by visiting the page
	err := collyCollector.Visit(link)
	if err != nil {
		log.Println("error while scraping webpage")
	}

}
