package internetsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coder/hnsw"
	"github.com/gocolly/colly"
	"github.com/openai/openai-go"
)

type embeddingResponse struct {
	Response openai.CreateEmbeddingResponse
}

type embeddingPrompt struct {
	prompt []string
}

// InternetSearch is used to search the internet with an models query request
//
// First the model should generate a query.
// Then the tool should check the query to make sure its a valid url.
// Then the search tool should make a request to a privacy oriented search engine.
// Once this is done it should web scrape the top five websites, and store it in the contex box.
// Internet search should not be compared with the rest of tools because it can be
// dangerous if not used properly, hence why it is serperate.
func InternetSearch(query string) *hnsw.Graph[string] {

	g := hnsw.NewGraph[string]()

	collyCollector := colly.NewCollector()
	maxNumLinks := 0

	query = strings.ReplaceAll(query, " ", "+")

	queryurl := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s&atb=v477-4&ia=web", query)

	// Find and visit all links
	collyCollector.OnHTML("result__url", func(element *colly.HTMLElement) {
		//counter used to limit the number of websites
		if maxNumLinks < 3 {
			//searches for links on duckduckgo and creates links to individual sites to be scraped
			link := "https://" + strings.TrimSpace(element.Text)
			maxNumLinks++

			scrape(link, g)
		}
	})

	err := collyCollector.Visit(queryurl)
	if err != nil {
		log.Println("error while scraping duckduckgo")
	}

	return g
}

// used to scrape individual sites
func scrape(link string, g *hnsw.Graph[string]) {

	collyCollector := colly.NewCollector()

	//scrapes all paragraph elements from each webpage
	collyCollector.OnHTML("p", func(element *colly.HTMLElement) {
		//concatenates all text from paragraph elements
		text := []string{element.Text}
		var prompt embeddingPrompt
		prompt.prompt = text
		embeddingPrompt, err := json.Marshal(prompt)
		if err != nil {
			log.Println("The fucking json Marshal didn't fucking work you STUPID MOTHERFUCKER")
		}
		resp, err := http.Post("http://localhost:8082/emb/generate", "application/json", bytes.NewReader(embeddingPrompt))
		if err != nil {
			log.Println("Sorry Vito what was that? You cut out.")
		}
		defer resp.Body.Close()

		var embeddingResponse embeddingResponse

		err = json.NewDecoder(resp.Body).Decode(&embeddingResponse)
		if err != nil {
			log.Println("Failed to Read response body in internet search")
		}

		vectorizedText := embeddingResponse.Response.Data[0].Embedding
		textVector32 := make([]float32, len(vectorizedText))
		for i, val := range vectorizedText {
			textVector32[i] = float32(val)
		}

		g.Add(
			hnsw.MakeNode(element.Text, textVector32),
		)

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
