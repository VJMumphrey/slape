package main

import (
	"context"
	"fmt"
	"io"

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

func CreateContainer(apiClient *client.Client, portNum string, name string) (container.CreateResponse, error) {
	port := fmt.Sprintf("%s/tcp", portNum)

	portSet := nat.PortSet{
		nat.Port(port): struct{}{}, // map 11434 TCP port
	}

	portBindings := nat.PortMap{
		nat.Port(port): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: portNum,
			},
		},
	}

	// create container
	createResponse, err := apiClient.ContainerCreate(context.Background(), &container.Config{
		ExposedPorts: portSet,
		Image:        "ghcr.io/ggerganov/llama.cpp:server",
	}, &container.HostConfig{
		//Runtime: "nvidia",
		PortBindings: portBindings,
		Mounts: []mount.Mount{{
			Type:     mount.TypeVolume,
			Source:   "models",
			Target:   "/models",
			ReadOnly: false,
		}},
	}, nil, nil, name)

	return createResponse, err
}

// This is very simple for right now but when we add structured outputs it will
// get very complicated.
func GenerateCompletion(prompt string) (*openai.ChatCompletion, error) {
	chatCompletion, err := openaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.String("llama3.2"),
	})

	return chatCompletion, err
}
