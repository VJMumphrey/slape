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
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	openaiClient = openai.NewClient(
		option.WithBaseURL("http://localhost:11434/v1/"),
	)
)

type simple struct {
	Prompt string `json:"prompt"`
}

type generate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream string `json:"stream"`
}

type pullModel struct {
	Model    string `json:"model"`
	Insecure string `json:"insecure"`
	Stream   string `json:"stream"`
}

func simplerequest(w http.ResponseWriter, req *http.Request) {

	var simplePayload simple

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		return
	}

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	// This checks for the image before pulling
	reader, err := apiClient.ImagePull(ctx, "ghcr.io/ggerganov/llama.cpp:server", image.PullOptions{All: false, RegistryAuth: ""})
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error pulling the image"))
		return
	}
	io.Copy(os.Stdout, reader)

	portSet := nat.PortSet{
		nat.Port("8000/tcp"): struct{}{}, // map 11434 TCP port
	}

	portBindings := nat.PortMap{
		"8000/tcp": []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "8000",
			},
		},
	}

	// create container
	var createResponse container.CreateResponse
	createResponse, err = apiClient.ContainerCreate(context.Background(), &container.Config{
		ExposedPorts: portSet,
		//Cmd:          []string{"ollama", "run", "llama3.2:1b"},
		Image: "ghcr.io/ggerganov/llama.cpp:server",
	}, &container.HostConfig{
		//Runtime: "nvidia",
		PortBindings: portBindings,
		Mounts: []mount.Mount{{
			Type:     mount.TypeVolume,
			Source:   "models",
			Target:   "/models",
			ReadOnly: false,
		}},
	}, nil, nil, "llamacpp")
	if err != nil {
		log.Println(err)
		w.Write([]byte("Error creating the container"))
		return
	}

	// start container
	if err := apiClient.ContainerStart(context.Background(), createResponse.ID, container.StartOptions{}); err != nil {
		log.Println(err)
		w.Write([]byte("Error starting the container"))
		return
	}

	log.Println(createResponse.ID)

	// generate a response
	chatCompletion, err := openaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(simplePayload.Prompt),
		}),
		Model: openai.String("llama3.2"),
	})
	if err != nil {
		panic(err.Error())
	}

	// For debugging
	log.Println(chatCompletion.Choices[0].Message.Content)

	// TODO json the response
	w.Write([]byte(chatCompletion.Choices[0].Message.Content))
}

func main() {

	http.HandleFunc("/simple", simplerequest)
	color.Green("[+] Server started on :3069")

	http.ListenAndServe(":3069", nil)
}
