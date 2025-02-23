// Package pipline is meant to help facilitate the running of several models in sequential order.
//
// For machines with small amounts of resources the pipelines will have to manage models by orchestrating models as need be.
// Pipelines also dictate the structure of conversations between the models. See the following docs for more info.
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
}

func (s *SimplePipeline) Setup(ctx context.Context) error {

	_, err := PullImage(s.DockerClient, ctx, s.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	color.Green("Pulling Image...")
	// prints out the status of the download
	// worth while for big images
	// io.Copy(os.Stdout, reader)

	createResponse, err := CreateContainer(s.DockerClient, "8000", "", ctx, "Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf", s.ContainerImage)
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
}

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

	createResponse, err := CreateContainer(c.DockerClient, "8000", "", ctx, "Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf", c.ContainerImage)
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

// DebateofModels is pipeline for debate structured prompting.
// Models talk in a round robin style.
type DebateofModels struct {
	Models []string
	ContextBox
	Tools
	Active         bool
	ContainerImage string
}

// InitDebateofModels creates a DebateofModels pipeline for debates.
// Includes a ContextBox and all models needed.
func (d *DebateofModels) Setup(ctx context.Context, cli *client.Client) error {
	reader, err := PullImage(cli, ctx, d.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	return nil
}

func (d *DebateofModels) Generate(prompt string, systemprompt string, maxtokens int64, openaiClient *openai.Client, cli *client.Client) {
}

func (d *DebateofModels) Shutdown(ctx context.Context, cli *client.Client) {}
