// Package api is used to create general functions needed for the api of slape.
// Some of which are defined in the pipelines and therefore are not included in this package.
package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func Cors(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Option", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

func UpDog(port string) bool {
	// swagger:operation GET /updog UpDog
	//
	// UpDog is a check to make sure llamacpp is ready for requests.
	//
	// Responses:
	//	200: StatusOk
	resp, err := http.Get("http://localhost:" + port + "/health")
	if err != nil {
		color.Red("%s", err)
		return false
	}

	switch resp.StatusCode {
	case http.StatusOK:
		color.Green("Model is ready")
		return true
	case http.StatusServiceUnavailable:
		color.Yellow("Model is not ready...")
		return false
	default:
		color.Yellow("Model is not ready...")
		return false
	}
}

// swagger:model ModelsResponse
type ModelsResponse struct {
	// List of models to be returned to the client
	// required: true
	Models []string `json:"models"`
}

// swagger:route GET /getmodels
//
// get models that the backend can see in the ./models folder
//
// Responses:
//
//	200: ModelsResponse
func GetModels(w http.ResponseWriter, req *http.Request) {
	files, err := os.ReadDir("./models")
	if err != nil {
		color.Red("%s", err)
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

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	return
}

// TODO(t) implement this. Download the gguf model into the ./models folder
func DownloadHuggingFaceModel(w http.ResponseWriter, req *http.Request) {}
