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

type simple struct {
	Prompt string `json:"prompt"`
}

func simplerequest(w http.ResponseWriter, req *http.Request) {

	var simplePayload simple

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	reader, err := PullImage(apiClient, ctx)
	if err != nil {
		color.Red("%s", err)
		w.Write([]byte("Error pulling the image"))
		return
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

    createResponse, err := CreateContainer(apiClient, "8000", "llamacpp", ctx)
	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		w.Write([]byte("Error creating the container"))
		return
	}

	// start container
    err = apiClient.ContainerStart(ctx, createResponse.ID, container.StartOptions{})
    if err != nil {
		color.Red("%s", err)
		w.Write([]byte("Error starting the container"))
		return
	}

	// For debugging
	log.Println(createResponse.ID)

	// generate a response
	chatCompletion, err := GenerateCompletion(simplePayload.Prompt)
	if err != nil {
		color.Red("%s", err)
		w.Write([]byte("Error getting generation from model"))
		return
	}

	// For debugging
	color.Green(chatCompletion.Choices[0].Message.Content)

	// TODO json the response
	w.Write([]byte(chatCompletion.Choices[0].Message.Content))
}

func main() {

	http.HandleFunc("/simple", simplerequest)
	color.Green("[+] Server started on :3069")

	http.ListenAndServe(":3069", nil)
}
