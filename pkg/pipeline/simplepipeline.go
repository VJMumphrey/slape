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

	"github.com/StoneG24/slape/internal/logging"
	"github.com/StoneG24/slape/internal/vars"
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
		// models used in the pipeline
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

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("ErrorString", "Errorstring", err)
		return
	}


	err = json.NewDecoder(req.Body).Decode(&setupPayload)
	if err != nil {
		logger.Error("Error", "ErrorString", err)
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

    // Create the logger for this
	logger := logging.CreateLogger()

	var simplePayload simpleRequest

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		logger.Error("Error", "ErrorString", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	promptChoice, maxtokens := processPrompt(simplePayload.Mode)

	s.ContextBox.SystemPrompt = promptChoice
	s.ContextBox.Prompt = simplePayload.Prompt
	s.Thinking, err = strconv.ParseBool(simplePayload.Thinking)
	if err != nil {
		logger.Error("Error", "Errorstring", err)
		http.Error(w, "Error parsing thinking value. Expecting sound boolean definitions.", http.StatusBadRequest)
	}
	if s.Thinking {
		thoughts, err := s.getThoughts()
		if err != nil {
			logger.Error("Error", "Errorstring", err)
			http.Error(w, "Error gathering thoughts", http.StatusInternalServerError)
		}
		s.ContextBox.Thoughts = thoughts
	}

	// generate a response
	result, err := s.Generate(maxtokens, vars.OpenaiClient, logger)
	if err != nil {
		logger.Error("Error", "ErrorString", err)
		http.Error(w, "Error getting generation from model", http.StatusOK)

		return
	}

	// for debugging streaming
	logger.Debug(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		logger.Error("Error", "ErrorString", err)
		http.Error(w, "Error marshaling your response from model", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *SimplePipeline) Setup(ctx context.Context, logger *slog.Logger) error {

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
		ctx,
		s.Models[0],
		s.ContainerImage,
		s.GPU,
	)

	if err != nil {
		logger.Warn("Warn", "WarningString", createResponse.Warnings)
		logger.Error("Error", "ErrorString", err)
		return err
	}

	// start container
	err = (s.DockerClient).ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
	if err != nil {
		logger.Error("Error", "ErrorString", err)
		return err
	}

	// For debugging
	logger.Info(createResponse.ID)
	s.container = createResponse

	return nil

}

func (s *SimplePipeline) Generate(maxtokens int64, openaiClient *openai.Client, logger *slog.Logger) (string, error) {
	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	logger.Debug("SystemPrompt", s.ContextBox.SystemPrompt, "Prompt", s.ContextBox.Prompt)

	err := s.PromptBuilder("")
	if err != nil {
		return "", err
	}

	logger.Debug("SystemPrompt", s.SystemPrompt)

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

	result, err := GenerateCompletion(param, "", *openaiClient)
	if err != nil {
		return "", err
	}

	return result, nil
}

// BUG(v) if the server is shut off then the container ids are also lost leaving us with orphan containers. This is for all pipelines.
func (s *SimplePipeline) Shutdown(w http.ResponseWriter, req *http.Request) {
    logger := logging.CreateLogger()

	err := (s.DockerClient).ContainerStop(context.Background(), s.container.ID, container.StopOptions{})
	if err != nil {
		logger.Error("Error", "ErrorString", err)
	}

	err = (s.DockerClient).ContainerRemove(context.Background(), s.container.ID, container.RemoveOptions{})
	if err != nil {
		logger.Error("Error", "ErrorString", err)
	}

	logger.Info("Shutting Down...")
}
