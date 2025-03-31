package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/StoneG24/slape/internal/vars"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	// change if you want to make things go faster for testing
	rounds = 3
)

type (
	// DebateofModels is pipeline for debate structured prompting.
	// Models talk in a round robin style.
	// According to the paper, Improving Factuality and Reasoning in Language Models through Multiagent Debate, pg8, https://arxiv.org/abs/2305.14325,
	// 3-4 rounds was the best range. There wasn't much of an improvement from 3 to 4 and greater. Since we are constrained on resources and compute time, we'll use 3.
	DebateofModels struct {
		Models []string
		ContextBox
		Tools
		Active         bool
		ContainerImage string
		DockerClient   *client.Client
		GPU            bool
		Thinking       bool

		// for internal use only
		containers []container.CreateResponse
	}

	debateRequest struct {
		// Prompt is the string that
		// will be appended to the prompt
		// string chosen.
		Prompt string `json:"prompt"`

		// Options are strings matching
		// the names of prompt types
		Mode string `json:"mode"`

		// Should thinking be a step in the process
		Thinking bool `json:"thinking"`
	}

	debateSetupPayload struct {
		Models []string `json:"models"`
	}

	debateResponse struct {
		Answer string `json:"answer"`
	}
)

// DebatePipelineSetupRequest, handlerfunc expects POST method and returns nothing
func (d *DebateofModels) DebatePipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("%s", err)
		return
	}
	go api.Cors(w, req)

	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var setupPayload debateSetupPayload

	err = json.NewDecoder(req.Body).Decode(&setupPayload)
	if err != nil {
		slog.Error("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	d.Models = setupPayload.Models
	d.DockerClient = apiClient

	go d.Setup(context.Background())

	w.WriteHeader(http.StatusOK)
}

// DebatePipelineGenerateRequest is used to handle the request for a debate style thought process.
func (d *DebateofModels) DebatePipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	api.Cors(w, req)

	var payload debateRequest

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		slog.Error("%s", err)
		http.Error(w, "Error unexpected request format", http.StatusUnprocessableEntity)
		return
	}

	promptChoice, maxtokens := processPrompt(payload.Mode)
	d.ContextBox.SystemPrompt = promptChoice
	d.ContextBox.Prompt = payload.Prompt

	if d.Thinking {
		thoughts, err := d.getThoughts()
		if err != nil {
			slog.Error("Error", "errorstring", err)
		}
		d.ContextBox.Thoughts = thoughts
	}

	// generate a response
	result, err := d.Generate(payload.Prompt, promptChoice, maxtokens)
	if err != nil {
		slog.Error("%s", err)
		http.Error(w, "Error getting generation from model", http.StatusOK)
		return
	}

	// for debugging streaming
	log.Print(result)

	respPayload := debateResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		slog.Error("%s", err)
		http.Error(w, "Error marshaling your response from model", http.StatusOK)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

// InitDebateofModels creates a DebateofModels pipeline for debates.
// Includes a ContextBox and all models needed.
func (d *DebateofModels) Setup(ctx context.Context) error {
	reader, err := PullImage(d.DockerClient, ctx, d.ContainerImage)
	if err != nil {
		slog.Error("%s", err)
		return err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	for i := 0; i < len(d.Models)-1; i++ {
		createResponse, err := CreateContainer(d.DockerClient, "800"+strconv.Itoa(i), "", ctx, d.Models[i], d.ContainerImage, d.GPU)
		if err != nil {
			slog.Warn("%s", createResponse.Warnings)
			slog.Error("%s", err)
			return err
		}

		slog.Info("%s", createResponse.ID)
		d.Models[i] = createResponse.ID
	}

	return nil
}

func (d *DebateofModels) Generate(prompt string, systemprompt string, maxtokens int64) (string, error) {
	var result string

	for j := 0; j < rounds; j++ {

		for i, model := range d.containers {
			// start container
			err := (d.DockerClient).ContainerStart(context.Background(), model.ID, container.StartOptions{})
			if err != nil {
				slog.Error("Error", "errorstring", err)
				return "", err
			}
			slog.Info("Starting container %d...", i)

			for {
				// sleep and give server guy a break
				time.Sleep(time.Duration(5 * time.Second))

				if api.UpDog("800" + strconv.Itoa(i)) {
					break
				}
			}

			openaiClient := openai.NewClient(
				option.WithBaseURL("http://localhost:800" + strconv.Itoa(i) + "/v1"),
			)

			slog.Debug("Debug: %s%s", systemprompt, prompt)

			err = d.PromptBuilder(result)
			if err != nil {
				return "", err
			}

			// Answer the initial question.
			// If it's the first model, there will not be any questions from the previous model.
			param := openai.ChatCompletionNewParams{
				Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(d.SystemPrompt),
					openai.UserMessage(d.Prompt),
					openai.UserMessage(d.FutureQuestions),
				}),
				Seed:        openai.Int(0),
				Model:       openai.String(d.Models[i]),
				Temperature: openai.Float(vars.ModelTemperature),
				MaxTokens:   openai.Int(maxtokens),
			}

			result, err = GenerateCompletion(param, "", *openaiClient)
			if err != nil {
				slog.Error("%s", err)
				return "", err
			}

			// Summarize the answer generate.
			// This apparently makes it easier for the next models to digest the information.
			summarizePrompt := fmt.Sprintf("Given this answer %s, can you summarize it", result)
			param = openai.ChatCompletionNewParams{
				Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(d.SystemPrompt),
					openai.UserMessage(summarizePrompt),
					//openai.UserMessage(s.FutureQuestions),
				}),
				Seed:        openai.Int(0),
				Model:       openai.String(d.Models[i]),
				Temperature: openai.Float(vars.ModelTemperature),
				MaxTokens:   openai.Int(maxtokens),
			}

			result, err = GenerateCompletion(param, "", *openaiClient)
			if err != nil {
				slog.Error("%s", err)
				return "", err
			}

			// Ask the model to generate questions for the model to answer.
			// Then store this answer in the contextbox for the next go around.
			/*
				askFutureQuestions := fmt.Sprintf("Given this answer, %s, can you make some further questions to ask the next model in order to aid in answering the question?", result)
				param = openai.ChatCompletionNewParams{
					Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
						openai.SystemMessage(d.SystemPrompt),
						openai.UserMessage(askFutureQuestions),
						//openai.UserMessage(s.FutureQuestions),
					}),
					Seed:        openai.Int(0),
					Model:       openai.String(d.Models[i]),
					Temperature: openai.Float(vars.ModelTemperature),
					MaxTokens:   openai.Int(maxtokens),
				}

				result, err = GenerateCompletion(param, "", *openaiClient)
				if err != nil {
					slog.Error("%s", err)
					return "", err
				}
			*/

			d.FutureQuestions = result

			slog.Info("Stopping container %d...", i)
			(d.DockerClient).ContainerStop(context.Background(), model.ID, container.StopOptions{})
		}
	}

	return result, nil
}

func (d *DebateofModels) Shutdown(w http.ResponseWriter, req *http.Request) {
	// turn off the containers if they aren't already off
	for i := range d.Models {
		(d.DockerClient).ContainerStop(context.Background(), d.Models[i], container.StopOptions{})
	}

	// remove the containers
	for i := range d.Models {
		(d.DockerClient).ContainerRemove(context.Background(), d.Models[i], container.RemoveOptions{})
	}

	slog.Info("Shutting Down...")
}
