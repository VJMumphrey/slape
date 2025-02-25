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

// ChainofModels is the next step above smallest pipeline.
// This pipeline contains a ContextBox and the models in squential order.
// ChainofModels forces the models to talk in sequential order
// like the name suggests.
type ChainofModels struct {
	Model1 string
	Model2 string
	Model3 string
	ContextBox
	Tools
	Active         bool
	ContainerImage string
	DockerClient   *client.Client
	GPU            bool
}

// InitChainofModels creates a ChainofModels pipeline.
// Includes a ContextBox and all models needed in squential order.
func (c *ChainofModels) Setup(ctx context.Context) error {
	reader, err := PullImage(c.DockerClient, ctx, c.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	createResponse, err := CreateContainer(c.DockerClient, "8000", "", ctx, "Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf", c.ContainerImage, c.GPU)
	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return err
	}

	// For debugging
	color.Green("%s", createResponse.ID)

	c.Model1 = createResponse.ID

	return nil
}

// TODO need to finish the model spin up and down later
func (c *ChainofModels) Generate(prompt string, systemprompt string, maxtokens int64, openaiClient *openai.Client, cli *client.Client) (string, error) {
	// start container
	err := cli.ContainerStart(context.Background(), c.Model1, container.StartOptions{})
	if err != nil {
		color.Red("%s", err)
		return "", err
	}

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

func (c *ChainofModels) Shutdown(ctx context.Context, cli *client.Client) {
	// turn off the containers if they aren't already off
	cli.ContainerStop(ctx, c.Model1, container.StopOptions{})
	cli.ContainerStop(ctx, c.Model2, container.StopOptions{})
	cli.ContainerStop(ctx, c.Model3, container.StopOptions{})

	// remove the containers
	cli.ContainerRemove(ctx, c.Model1, container.RemoveOptions{})
	cli.ContainerRemove(ctx, c.Model2, container.RemoveOptions{})
	cli.ContainerRemove(ctx, c.Model3, container.RemoveOptions{})
}
