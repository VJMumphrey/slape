package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

var (
	// for ollama, change as needed
	u, err = url.Parse("http://localhost:11434")

	client = api.NewClient(u, http.DefaultClient)
)

// Startup allows us to pre-load models before we send a request.
// This is done by sending an empty request to /api/generate.
func Startup() {
	req := &api.GenerateRequest{
		Model: "smollm2:135m",
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
	client.Generate(ctx, req, respFunc)
}

func GetSimpleAnswer(prompt string) error {
	byteSimple := []byte(`
        {"respone": "value"}
    `)

	// By default, GenerateRequest is streaming.
	req := &api.GenerateRequest{
		Model:  "smollm2:135m",
		Prompt: prompt,
		// TODO(v) once we have websockets we can
		// re-enable this
		Stream: &[]bool{false}[0],
		Format: json.RawMessage(byteSimple),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.

		// In streaming mode, responses are partial so we call fmt.Print (and not
		// Println) in order to avoid spurious newlines being introduced. The
		// model will insert its own newlines if it wants.
		fmt.Println(resp.Response)
		return nil
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		return err
	}
	fmt.Println()

	return nil
}
