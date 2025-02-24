package pipeline

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
)

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
