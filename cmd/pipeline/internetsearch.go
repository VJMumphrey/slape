package pipeline

// InternetSearch is used to search the internet with an models query request
//
// First the model should generate a query.
// Then the tool should check the query to make sure its a valid url.
// Then the search tool should make a request to a privacy oriented search engine.
// Once this is done it should scrape the top five websites,
// index the information into the vectore store file and store it in the contex box.
func InternetSearch() {}
