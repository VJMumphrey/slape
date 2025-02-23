/*
Package SLaPE is a binary that starts a pod on the local computer using as socket to podman.

Usage:

	./slape

Containerized models are spawned as needed adhering to a pipeline system.
*/
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StoneG24/slape/cmd/pipeline"
	"github.com/StoneG24/slape/cmd/prompt"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	openaiClient = openai.NewClient(
		option.WithBaseURL("http://localhost:8000/v1"),
	)
)

type simpleRequest struct {

	// Prompt is the string that
	// will be appended to the prompt
	// string chosen.
	Prompt string `json:"prompt"`
	Model  string `json:"model ,omitempty"`

	// Options are strings matching
	// the names of prompt types
	Mode string `json:"mode"`
}

type simpleResponse struct {
	Answer string `json:"answer"`
}

// cors is used to handle cors for each HandleFunc that we create.
func cors(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Option", "GET, POST, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	}
}

// simplerequest is used to handle simple requests as needed.
func simplerequest(w http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	cors(w, req)

	s := pipeline.SimplePipeline{
		// updates after created
		Model:      "",
		ContextBox: pipeline.ContextBox{},
		Tools:      pipeline.Tools{},
		Active:     true,
	}
	go s.Setup(ctx, apiClient)

	w.Header().Set("Content-Type", "application/json")

	var simplePayload simpleRequest

	err = json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	var promptChoice string
	var maxtokens int64

	switch simplePayload.Mode {
	case "simple":
		promptChoice = prompt.SimplePrompt
		maxtokens = 100
	case "cot":
		promptChoice = prompt.CoTPrompt
		maxtokens = 4096
	case "tot":
		promptChoice = prompt.ToTPrompt
		maxtokens = 32768
	case "got":
		promptChoice = prompt.GoTPrompt
		maxtokens = 32768
	case "moe":
		promptChoice = prompt.MoEPrompt
		maxtokens = 32768
	case "thinkinghats":
		promptChoice = prompt.SixThinkingHats
		maxtokens = 32768
	default:
		promptChoice = prompt.SimplePrompt
		maxtokens = 100
	}

	// for debugging
	color.Yellow(promptChoice)

	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))
		resp, err := http.Get("http://localhost:8000/health")
		if err != nil {
			color.Red("%s", err)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			break
		} else if resp.StatusCode == http.StatusServiceUnavailable {
			continue
		}
	}

	// generate a response
	result, err := s.Generate(simplePayload.Prompt, promptChoice, maxtokens, openaiClient)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error getting generation from model"))
		return
	}

	go s.Shutdown(ctx, apiClient)

	// for debugging streaming
	color.Green(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error marshaling your response from model"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// request GET for backend check to make sure llamacpp is ready for requests.
// returns 200 ok when things are ready
func upDog(w http.ResponseWriter, req *http.Request) {

	cors(w, req)

	resp, err := http.Get("http://localhost:8000/health")
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("error while checking model load status..."))
		return
	}

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	} else if resp.StatusCode == http.StatusServiceUnavailable {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("loading model..."))
		return
	}
}

// CheckGPU is a helper func that is meant to find a gpu on the system.
// If a gpu is found we should use the gpu since it will be a faster alternative to CPU.
// This means we have to download a GPU image of llamacpp and change the runtime to nvidia or something else.
func CheckGPU() {

}

func main() {
	CheckGPU()

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	reader, err := pipeline.PullImage(apiClient, ctx)
	if err != nil {
		color.Red("%s", err)
		return
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	http.HandleFunc("/simple", simplerequest)
	http.HandleFunc("/up", upDog)

	// Create a new HTTP server.
	srv := &http.Server{
		Addr: ":3069",
	}

	// Start the server in a goroutine.
	color.Green("[+] Server started on :3069")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// Create a channel to listen for interrupt signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Block until a signal is received.
	<-sigChan

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %s", err)
	}

	log.Println("Server gracefully stopped")
}
