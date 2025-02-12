package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/openai/openai-go"
)

// PullImage uses the docker api to pull an image down.
// The function also checks for the image locally before pulling.
func PullImage(apiClient *client.Client, ctx context.Context) (io.ReadCloser, error) {
	reader, err := apiClient.ImagePull(ctx, "ghcr.io/ggerganov/llama.cpp:server", image.PullOptions{All: false, RegistryAuth: ""})

	return reader, err
}

func CreateContainer(apiClient *client.Client, portNum string, name string, ctx context.Context) (container.CreateResponse, error) {

	portSet := nat.PortSet{
		nat.Port(portNum + "/tcp"): struct{}{}, // map 11434 TCP port
	}

	portBindings := nat.PortMap{
		nat.Port(portNum + "/tcp"): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "8000",
			},
		},
	}

    mountString := os.Getenv("PWD") + "/models"

	// create container
	createResponse, err := apiClient.ContainerCreate(ctx, &container.Config{
		ExposedPorts: portSet,
		Image:        "ghcr.io/ggerganov/llama.cpp:server",
		Cmd:          []string{"-m", "/models/Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf", "--port", "8000", "--host", "0.0.0.0", "-n", "512"},
	}, &container.HostConfig{
		//Runtime: "nvidia",
		/*
			Binds: []string{
				"/models:/models",
			},
		*/
		PortBindings: portBindings,
		Mounts: []mount.Mount{{
			Type:   mount.TypeBind,
			Source: mountString,
			Target: "/models",
		}},
	}, nil, nil, name)

	return createResponse, err
}

// This is very simple for right now but when we add structured outputs it will
// get very complicated.
func GenerateCompletion(prompt string) (*openai.ChatCompletion, error) {
	chatCompletion, err := openaiClient.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.String("llama3.2"),
	})

	return chatCompletion, err
}

// FindContainer finds a specific container based on the nomenclature of /name.
// Useful making checks before
func FindContainer(apiClient *client.Client, ctx context.Context) (types.Container, bool) {
	status := false

	containers, err := apiClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if container.Names[len(container.Names)-1] == "/llamacpp" {
			return container, true
		}
	}

	return types.Container{}, status
}
