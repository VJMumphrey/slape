// Package api is used to create general functions needed for the api of slape.
// Some of which are defined in the pipelines and therefore are not included in this package.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/StoneG24/slape/pkg/vars"
)

func UpDog(port string) bool {
	resp, err := http.Get("http://localhost:" + port + "/health")
	if err != nil {
		log.Println("Error checking model", err)
		return false
	}

	switch resp.StatusCode {
	case http.StatusOK:
		log.Println("Model is ready")
		return true
	case http.StatusServiceUnavailable:
		log.Println("Model is not ready...")
		return false
	default:
		log.Println("Model is not ready...")
		return false
	}
}

// swagger:model ModelsResponse
type ModelsResponse struct {
	// List of models to be returned to the client
	// required: true
	Models []string `json:"models"`
}

func GetModels(w http.ResponseWriter, req *http.Request) {
	files, err := os.ReadDir("./models")
	if err != nil {
		log.Println("Error reading models folder", err)
		w.Write([]byte("Error while trying to read models files"))
	}

	bundle := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			bundle = append(bundle, file.Name())
		}
	}

	respBundle := ModelsResponse{
		Models: bundle,
	}

	json, err := json.Marshal(respBundle)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

	return
}

type LogRequest struct {
	Logs []byte `json:"logs"`
}

func GetLogs(w http.ResponseWriter, req *http.Request) {
	contents, err := os.ReadFile("./logs/" + vars.Logfilename)
	if err != nil {
		log.Panicln("Error getting logs for frontend", err)
	}

	logs := LogRequest{
		Logs: contents,
	}

	jsonLogs, err := json.Marshal(logs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonLogs)
}
