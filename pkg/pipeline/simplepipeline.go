package pipeline

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/StoneG24/slape/pkg/vars"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/openai/openai-go"
)

type (
	// SimplePipeline is the smallest pipeline.
	// It contains only a model with a ContextBox.
	// This is good for a giving a single model access to tools
	// like internet search.
	SimplePipeline struct {
		Models []string
		ContextBox
		Tools
		Active         bool
		ContainerImage string
		DockerClient   *client.Client
		GPU            bool
		Thinking       bool

		// for internal use
		container container.CreateResponse
	}

	simpleRequest struct {
		// Prompt is the string that
		// will be appended to the prompt
		// string chosen.
		Prompt string `json:"prompt"`

		// Options are strings matching
		// the names of prompt types
		Mode string `json:"mode"`

		// Should thinking be included in the process
		Thinking string `json:"thinking"`
	}

	simpleSetupPayload struct {
		// Models is the name of the single
		// model used in the pipeline
		Models []string `json:"models"`
	}

	simpleResponse struct {
		// Answer is a json string containing the answer is markdown format
		// along with the models thought process
		Answer string `json:"answer"`
	}
)

// SimplePipelineSetupRequest, handlerfunc expects POST method and returns no content
func (s *SimplePipeline) SimplePipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
	var setupPayload simpleSetupPayload

	ctx, cancel := context.WithDeadline(req.Context(), time.Now().Add(30*time.Second))
	defer cancel()

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&setupPayload)
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	s.Models = setupPayload.Models
	s.DockerClient = apiClient

	s.Setup(ctx)

	w.WriteHeader(http.StatusOK)
}

// simplerequest is used to handle simple requests as needed.
func (s *SimplePipeline) SimplePipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	var simplePayload simpleRequest

	// use this to scope the context to the request
	ctx, cancel := context.WithDeadline(req.Context(), time.Now().Add(5*time.Minute))
	defer cancel()

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	promptChoice, maxtokens := processPrompt(simplePayload.Mode)

	s.ContextBox.SystemPrompt = promptChoice
	s.ContextBox.Prompt = simplePayload.Prompt
	s.Thinking, err = strconv.ParseBool(simplePayload.Thinking)
	if err != nil {
		slog.Error("Error", "Errorstring", err)
		http.Error(w, "Error parsing thinking value. Expecting sound boolean definitions.", http.StatusBadRequest)
	}
	if s.Thinking {
		s.getThoughts(ctx)
	}

	// gather here and then generate
	result, err := s.Generate(ctx, maxtokens, vars.OpenaiClient)
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		http.Error(w, "Error getting generation from model", http.StatusInternalServerError)

		return
	}

	// for debugging streaming
	slog.Debug(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		http.Error(w, "Error marshaling your response from model", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *SimplePipeline) Setup(ctx context.Context) error {

	childctx, cancel := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
	defer cancel()

	slog.Debug("Debug", "PullingImage", s.ContainerImage)

	reader, err := PullImage(s.DockerClient, childctx, s.ContainerImage)
	if err != nil {
		slog.Error("Error", "ErrorPullingContainerImage", err)
		return err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	createResponse, err := CreateContainer(
		s.DockerClient,
		"8000",
		"",
		childctx,
		s.Models[0],
		s.ContainerImage,
		s.GPU,
	)

	if err != nil {
		slog.Warn("Warn", "WarningString", createResponse.Warnings)
		slog.Error("Error", "ErrorString", err)
		return err
	}

	// start container
	err = (s.DockerClient).ContainerStart(childctx, createResponse.ID, container.StartOptions{})
	if err != nil {
		slog.Error("Error", "ErrorString", err)
		return err
	}

	slog.Info("Info", "Starting Container", createResponse.ID)
	s.container = createResponse

	return nil
}

func (s *SimplePipeline) Generate(ctx context.Context, maxtokens int64, openaiClient *openai.Client) (string, error) {
	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(2 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	slog.Debug("Debug", "SystemPrompt", s.ContextBox.SystemPrompt, "Prompt", s.ContextBox.Prompt)

	err := s.PromptBuilder("")
	if err != nil {
		return "", err
	}

	slog.Debug("Debug", "SystemPrompt", s.SystemPrompt)

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(s.SystemPrompt),
			openai.UserMessage(s.Prompt),
			//openai.UserMessage(s.FutureQuestions),
		}),
		Seed:        openai.Int(0),
		Model:       openai.String(s.Models[0]),
		Temperature: openai.Float(vars.ModelTemperature),
		MaxTokens:   openai.Int(maxtokens),
	}

	result, err := GenerateCompletion(ctx, param, "", *openaiClient)
	if err != nil {
		return "", err
	}

	return result, nil
}

// BUG(v) if the server is shut off then the container ids are also lost leaving us with orphan containers. This is for all pipelines.
func (s *SimplePipeline) Shutdown(w http.ResponseWriter, req *http.Request) {

	childctx, cancel := context.WithDeadline(req.Context(), time.Now().Add(30*time.Second))
	defer cancel()

	err := (s.DockerClient).ContainerStop(childctx, s.container.ID, container.StopOptions{})
	if err != nil {
		slog.Error("Error", "ErrorString", err)
	}

	err = (s.DockerClient).ContainerRemove(childctx, s.container.ID, container.RemoveOptions{})
	if err != nil {
		slog.Error("Error", "ErrorString", err)
	}

	slog.Info("Shutting Down...")
}
