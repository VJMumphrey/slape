// Package api is used to create general functions needed for the api of slape.
// Some of which are defined in the pipelines and therefore are not included in this package.
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/StoneG24/slape/pkg/vars"
)

func UpDog(port string) bool {
	// swagger:operation GET /updog UpDog
	//
	// UpDog is a check to make sure llamacpp is ready for requests.
	//
	// Responses:
	//	200: StatusOk
	resp, err := http.Get("http://localhost:" + port + "/health")
	if err != nil {
		slog.Error("%s", err)
		return false
	}

	switch resp.StatusCode {
	case http.StatusOK:
		slog.Info("Model is ready")
		return true
	case http.StatusServiceUnavailable:
		slog.Warn("Model is not ready...")
		return false
	default:
		slog.Warn("Model is not ready...")
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
		slog.Error("Error", "Errorstring", err)
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

func ShutdownPipes(w http.ResponseWriter, req *http.Request) {
	ShutdownPipelines()

	w.WriteHeader(http.StatusOK)
	return
}

// shutdownPipelines is used to shutdown pipelines with
// a remote request since the shutdown functions are now http.HandleFunc
func ShutdownPipelines() error {

	url := "http://localhost:8080/%s/shutdown"
	pipelines := []string{"simple", "cot", "deb", "emb"}

	for _, pipeline := range pipelines {
		requrl := fmt.Sprintf(url, pipeline)
		resp, err := http.Get(requrl)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}
