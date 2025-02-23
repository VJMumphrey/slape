// Package pipline is meant to help facilitate the running of several models in sequential order.
//
// For machines with small amounts of resources the pipelines will have to manage models by orchestrating models as need be.
// Pipelines also dictate the structure of conversations between the models. See the following docs for more info.
package pipeline

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
)

type Pipeline interface {
	// ctx, docker client
	Setup(context.Context, *client.Client) error

	// userprompt, systemprompt, maxtokens, openaiClient
	Generate(string, string, int64, *openai.Client) (string, error)

	// ctx, docker client
	Shutdown(context.Context, *client.Client)
}

// SimplePipeline is the smallest pipeline.
// It contains only a model with a ContextBox.
// This is good for a giving a single model access to tools
// like internet search.
type SimplePipeline struct {
	// container.CreateResponse ID
	Model string
	ContextBox
	Tools
	Active bool
}

func (s *SimplePipeline) Setup(ctx context.Context, cli *client.Client) error {
	createResponse, err := CreateContainer(cli, "8000", "", ctx, "/models/Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf")
	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return err
	}

	// start container
	err = cli.ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
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
	Active bool
}

// InitChainofModels creates a ChainofModels pipeline.
// Includes a ContextBox and all models needed in squential order.
func (c *ChainofModels) Setup(ctx context.Context, cli *client.Client) error {
	createResponse, err := CreateContainer(cli, "8000", "", ctx, "/models/Dolphin3.0-Llama3.2-1B-Q4_K_M.gguf")
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
	Active bool
}

// InitDebateofModels creates a DebateofModels pipeline for debates.
// Includes a ContextBox and all models needed.
func (d *DebateofModels) Setup(ctx context.Context, cli *client.Client) {}

func (d *DebateofModels) Generate(prompt string, systemprompt string, maxtokens int64, openaiClient *openai.Client, cli *client.Client) {
}

func (d *DebateofModels) Shutdown(ctx context.Context, cli *client.Client) {}
