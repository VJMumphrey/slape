package container

import (
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

var ollama_url = url.URL{Host: "localhost:1434"}

func setup() int {
	http_client := http.Client{}
	oclient := api.NewClient(&ollama_url, &http_client)

	oclient.Chat(ctx, req)

	return 0
}
