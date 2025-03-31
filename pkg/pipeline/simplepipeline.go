package pipeline

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

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
		Model string
		ContextBox
		Tools
		Active         bool
		ContainerImage string
		DockerClient   *client.Client
		GPU            bool

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
	}

	simpleSetupPayload struct {
		// Model is the name of the single
		// model used in the pipeline
		Models string `json:"model"`
	}

	simpleResponse struct {
		// Answer is a json string containing the answer is markdown format
		// along with the models thought process
		Answer string `json:"answer"`
	}
)

// SimplePipelineSetupRequest, handlerfunc expects POST method and returns no content
func (s *SimplePipeline) SimplePipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
        slog.Error("%s", err)
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
        slog.Error("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	s.Model = setupPayload.Models
	s.DockerClient = apiClient

	go s.Setup(context.Background())

	w.WriteHeader(http.StatusOK)
}

// simplerequest is used to handle simple requests as needed.
func (s *SimplePipeline) SimplePipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	go api.Cors(w, req)

	var simplePayload simpleRequest

	err := json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		slog.Error("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	promptChoice, maxtokens := processPrompt(simplePayload.Mode)

	s.ContextBox.SystemPrompt = promptChoice
	s.ContextBox.Prompt = simplePayload.Prompt
	thoughts, err := s.getThoughts()
    log.Println("%s", thoughts)
	s.ContextBox.Thoughts = thoughts

	// generate a response
	result, err := s.Generate(maxtokens, vars.OpenaiClient)
	if err != nil {
	    slog.Error("%s", err)
		http.Error(w, "Error getting generation from model", http.StatusOK)

		return
	}

	// for debugging streaming
    slog.Debug(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		slog.Error("Error", err)
		http.Error(w, "Error marshaling your response from model", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *SimplePipeline) Setup(ctx context.Context) error {

	reader, err := PullImage(s.DockerClient, ctx, s.ContainerImage)
	if err != nil {
		slog.Error("%s", err)
		return err
	}
	slog.Info("Pulling Image...")
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)
	defer reader.Close()

	createResponse, err := CreateContainer(
		s.DockerClient,
		"8000",
		"",
		ctx,
		s.Model,
		s.ContainerImage,
		s.GPU,
	)

	if err != nil {
		slog.Warn("%s", createResponse.Warnings)
		slog.Error("%s", err)
		return err
	}

	// start container
	err = (s.DockerClient).ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
	if err != nil {
		slog.Error("%s", err)
		return err
	}

	// For debugging
	slog.Info(createResponse.ID)
	s.container = createResponse

	return nil

}

func (s *SimplePipeline) Generate(maxtokens int64, openaiClient *openai.Client) (string, error) {
	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	slog.Debug("Debug: %s%s", s.ContextBox.SystemPrompt, s.ContextBox.Prompt)

	err := s.PromptBuilder("")
	if err != nil {
		return "", err
	}

	slog.Debug(s.SystemPrompt)

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(s.SystemPrompt),
			openai.UserMessage(s.Prompt),
			//openai.UserMessage(s.FutureQuestions),
		}),
		Seed:        openai.Int(0),
		Model:       openai.String(s.Model),
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
	err := (s.DockerClient).ContainerStop(context.Background(), s.container.ID, container.StopOptions{})
	if err != nil {
        slog.Error("%s", err)
	}

	err = (s.DockerClient).ContainerRemove(context.Background(), s.container.ID, container.RemoveOptions{})
	if err != nil {
        slog.Error("%s", err)
	}

	slog.Info("Shutting Down...")
}
