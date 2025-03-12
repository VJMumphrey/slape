// Package api is used to create general functions needed for the api of slape.
// Some of which are defined in the pipelines and therefore are not included in this package.
package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/fatih/color"
)

// cors is used to handle cors for each HandleFunc that we create.
func Cors(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Option", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

// request GET for backend check to make sure llamacpp is ready for requests.
// returns 200 ok when things are ready
func UpDog(port string) bool {

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

type modelsResponse struct {
	Models []string `json:"models"`
}

// request GET for client to get models that the backend can see
// in the ./models folder
// Returns a json string of key model and a value of array of strings
// - {models":["name"]}
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

	respBundle := modelsResponse{
		Models: bundle,
	}

	json, err := json.Marshal(respBundle)

	w.WriteHeader(http.StatusOK)
	w.Write(json)

	return
}

// TODO(t) implement this. Download the gguf model into the ./models folder
func DownloadHuggingFaceModel(w http.ResponseWriter, req *http.Request) {}
