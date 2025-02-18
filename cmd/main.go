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

	"github.com/StoneG24/slape/cmd/prompt"
	"github.com/docker/docker/api/types/container"
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
	Mode string `json:"mode ,omitempty"`
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

	cors(w, req)

	w.Header().Set("Content-Type", "application/json")

	var simplePayload simpleRequest

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	var promptChoice string

	switch simplePayload.Mode {
	case "simple":
		promptChoice = prompt.SimplePrompt
	case "cot":
		promptChoice = prompt.CoTPrompt
	case "tot":
		promptChoice = prompt.ToTPrompt
	case "got":
		promptChoice = prompt.GoTPrompt
	case "thinkinghats":
		promptChoice = prompt.SixThinkingHats
	default:
		promptChoice = prompt.SimplePrompt
	}

	// for debugging
	color.Yellow(promptChoice)

	// generate a response
	chatCompletion, err := GenerateCompletion(simplePayload.Prompt, promptChoice)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error getting generation from model"))
		return
	}

	// For debugging
	//color.Green(chatCompletion.Choices[0].Message.Content)

	// for debugging streaming
	color.Green(chatCompletion)

	respPayload := simpleResponse{
		Answer: chatCompletion,
	}

	json, err := json.Marshal(respPayload)
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// init function runs on startup to spin up required resoruces.
func setup(ctx context.Context, cli *client.Client, conts *[]container.CreateResponse) (string, error) {
	reader, err := PullImage(cli, ctx)
	if err != nil {
		color.Red("%s", err)
		return "", err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	createResponse, err := CreateContainer(cli, "8000", "", ctx)
	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return "", err
	}

	*conts = append(*conts, createResponse)

	// start container
	err = cli.ContainerStart(ctx, createResponse.ID, container.StartOptions{})
	if err != nil {
		color.Red("%s", err)
		return "", err
	}

	// For debugging
	log.Println(createResponse.ID)

	return createResponse.ID, nil
}

// shutdown function runs on shutdown and cleans up app resources.
func Shutdown(ctx context.Context, cli *client.Client, conts *[]container.CreateResponse) {
	for _, containerGuy := range *conts {
		cli.ContainerStop(ctx, containerGuy.ID, container.StopOptions{})

		cli.ContainerRemove(ctx, containerGuy.ID, container.RemoveOptions{})
	}
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

func main() {
	conts := []container.CreateResponse{}

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	go setup(ctx, apiClient, &conts)

	http.HandleFunc("/simple", simplerequest)
	http.HandleFunc("/up", upcheck)

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

	Shutdown(ctx, apiClient, &conts)

	log.Println("Server gracefully stopped")
}
