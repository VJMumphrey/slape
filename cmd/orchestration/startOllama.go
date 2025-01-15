package orchestration

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

func Startup(prompt string) error {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return err
	}

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model:  "phi3.5",
		Prompt: prompt,
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Print(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return err
	}
	fmt.Println()

	return nil
}
