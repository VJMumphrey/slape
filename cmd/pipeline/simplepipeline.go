package pipeline

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
)

// SimplePipeline is the smallest pipeline.
// It contains only a model with a ContextBox.
// This is good for a giving a single model access to tools
// like internet search.
type SimplePipeline struct {
	// container.CreateResponse ID
	Model string
	ContextBox
	Tools
	Active         bool
	ContainerImage string
	DockerClient   *client.Client
	GPU            bool
}

func (s *SimplePipeline) Setup(ctx context.Context) error {

	reader, err := PullImage(s.DockerClient, ctx, s.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	color.Green("Pulling Image...")
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	createResponse, err := CreateContainer(s.DockerClient, "8000", "", ctx, "Phi-3.5-mini-instruct-Q4_K_M.gguf", s.ContainerImage, s.GPU)
	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return err
	}

	// start container
	err = (s.DockerClient).ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
	if err != nil {
		color.Red("%s", err)
		return err
	}

	// For debugging
	color.Green("%s", createResponse.ID)

	s.Model = createResponse.ID

	return nil

}

func (s *SimplePipeline) Generate(prompt string, systemprompt string, maxtokens int64, openaiClient *openai.Client) (string, error) {

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemprompt),
			openai.UserMessage(prompt),
		}),
		Seed:      openai.Int(0),
		Model:     openai.String("llama3.2"),
		MaxTokens: openai.Int(maxtokens),
	}

	result, err := GenerateCompletion(param, "", *openaiClient)
	if err != nil {
		color.Red("%s", err)
		return "", err
	}

	return result, nil
}

func (s *SimplePipeline) Shutdown(ctx context.Context, cli *client.Client) {
	cli.ContainerStop(ctx, s.Model, container.StopOptions{})
	cli.ContainerRemove(ctx, s.Model, container.RemoveOptions{})

	color.Green("Shutting Down...")
}
