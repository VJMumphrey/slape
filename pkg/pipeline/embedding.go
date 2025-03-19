package pipeline

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/StoneG24/slape/internal/vars"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/internal/param"
	"github.com/openai/openai-go/shared"
)

const (
	embedmodel = "snowflake-arctic-embed-l-v2.0-q4_k_m.gguf"
	genmodel   = "Phi-3.5-mini-instruct-Q4_K_M.gguf"
)

// This pipeline is meant to be used for indexing a RAG database.
// We are using MiniRAG for a size complexity balance.
type EmbeddingPipeline struct {
	DockerClient   *client.Client
	ContainerImage string
	GPU            bool

	// for internal use
	container container.CreateResponse
}

// SimplePipelineSetupRequest, handlerfunc expects POST method and returns no content
func (e *EmbeddingPipeline) EmbeddingPipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	go api.Cors(w, req)

	if req.Method != http.MethodPost {
		http.Error(w, "Wrong method used for endpoint", http.StatusBadRequest)
		return
	}

	var setupPayload simpleSetupPayload

	err = json.NewDecoder(req.Body).Decode(&setupPayload)
	if err != nil {
		color.Red("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	e.DockerClient = apiClient

	go e.Setup(context.Background())

	w.WriteHeader(http.StatusOK)
}

// simplerequest is used to handle simple requests as needed.
func (e *EmbeddingPipeline) EmbeddingPipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	go api.Cors(w, req)

	var simplePayload simpleRequest

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	// generate a response
	// TODO rewrite for embedding and rag
	result, err := e.Generate(simplePayload.Prompt, vars.OpenaiClient)
	if err != nil {
		color.Red("%s", err)
		http.Error(w, "Error getting generation from model", http.StatusOK)
		go e.Shutdown(ctx)

		return
	}

	go e.Shutdown(ctx)

	// for debugging streaming
	color.Green(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		color.Red("%s", err)
		http.Error(w, "Error marshaling your response from model", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (e *EmbeddingPipeline) Setup(ctx context.Context) error {

	reader, err := PullImage(e.DockerClient, ctx, e.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	color.Green("Pulling Image...")
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	createResponse, err := CreateContainer(
		e.DockerClient,
		"8000",
		"",
		ctx,
		embedmodel,
		e.ContainerImage,
		e.GPU,
	)

	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return err
	}

	// start container
	err = (e.DockerClient).ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
	if err != nil {
		color.Red("%s", err)
		return err
	}

	// For debugging
	color.Green("%s", createResponse.ID)
	e.container = createResponse

	return nil

}

func (e *EmbeddingPipeline) Generate(payload string, openaiClient *openai.Client) (string, error) {
	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	param := openai.EmbeddingNewParams{
		Input: []string {
			payload,
		},
		Model:      openai.String(embedmodel),
		Dimensions: openai.Int(1024),
	},

	// should return a type of openai.Embedding
	result, err := GenerateCompletion(param, "", *openaiClient)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (e *EmbeddingPipeline) Shutdown(ctx context.Context) error {
	err := (e.DockerClient).ContainerStop(ctx, e.container.ID, container.StopOptions{})
	if err != nil {
		color.Red("%s", err)
		return nil
	}

	err = (e.DockerClient).ContainerRemove(ctx, e.container.ID, container.RemoveOptions{})
	if err != nil {
		color.Red("%s", err)
		return nil
	}

	color.Green("Shutting Down...")

	return nil
}
