package pipeline

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

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
func PullImage(apiClient *client.Client, ctx context.Context, containerImage string) (io.ReadCloser, error) {
	reader, err := apiClient.ImagePull(ctx, containerImage, image.PullOptions{All: false, RegistryAuth: ""})

	return reader, err
}

func CreateContainer(apiClient *client.Client, portNum string, name string, ctx context.Context, modelName string, containerImage string) (container.CreateResponse, error) {

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

	var mountString string

	if runtime.GOOS == "windows" {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Vito are less gay")
		}

		currentPath := filepath.Dir(ex)

		fmt.Println(currentPath)

		mountString = currentPath + "\\models"
	}

	if runtime.GOOS == "linux" {
		mountString = os.Getenv("PWD") + "/models"
	}

	// create container
	createResponse, err := apiClient.ContainerCreate(ctx, &container.Config{
		ExposedPorts: portSet,
		Image:        containerImage,
		Cmd:          []string{"-m", "/models/" + modelName, "--port", "8000", "--host", "0.0.0.0", "-n", "32678", "-fa"},
	}, &container.HostConfig{
		// TODO check for gpu, if true set to use nvidia runtime, rocm, or cdi
		//Runtime: "nvidia",
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
//
// prompt comes from a user and is the question being asked.
// systemprompt is the systemprompt chosen based on the prompting style requested.
func GenerateCompletion(param openai.ChatCompletionNewParams, followupQuestion string, openaiClient openai.Client) (string, error) {

	var result string

	stream := openaiClient.Chat.Completions.NewStreaming(context.Background(), param)

	// optionally, an accumulator helper can be used
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		//w.Write([]byte(chunk.Choices[0].Delta.Content))
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
		}

		// if using tool calls
		//if tool, ok := acc.JustFinishedToolCall(); ok {
		//	println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		//}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
		}

		// it's best to use chunks after handling JustFinished events
		if len(chunk.Choices) > 0 {
			println(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		return "", err
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	result = acc.Choices[0].Message.Content

	// Adding this for later
	//param.Messages.Value = append(param.Messages.Value, acc.Choices[0].Message)
	//param.Messages.Value = append(param.Messages.Value, openai.UserMessage(followupQuestion))

	return result, nil
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
