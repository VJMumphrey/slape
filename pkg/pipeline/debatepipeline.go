package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/StoneG24/slape/internal/vars"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	// change if you want to make things go faster for testing
	rounds = 3
)

// DebateofModels is pipeline for debate structured prompting.
// Models talk in a round robin style.
// According to the paper, Improving Factuality and Reasoning in Language Models through Multiagent Debate, pg8, https://arxiv.org/abs/2305.14325,
// 3-4 rounds was the best range. There wasn't much of an improvement from 3 to 4 and greater. Since we are constrained on resources and compute time, we'll use 3.
type DebateofModels struct {
	Models []string
	ContextBox
	Tools
	Active         bool
	ContainerImage string
	DockerClient   *client.Client
	GPU            bool

	// for internal use only
	containers []container.CreateResponse
}

type debateRequest struct {
	// Prompt is the string that
	// will be appended to the prompt
	// string chosen.
	Prompt string `json:"prompt"`

	// Options are strings matching
	// the names of prompt types
	Mode string `json:"mode"`
}

type debateSetupPayload struct {
	Models []string `json:"models"`
}

type debateResponse struct {
	Answer string `json:"answer"`
}

// DebatePipelineSetupRequest, handlerfunc expects POST method and returns nothing
func (d *DebateofModels) DebatePipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
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
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	d.PickImage()

	d.Models = setupPayload.Models
	d.DockerClient = apiClient
	d.GPU = IsGPU()

	go d.Setup(context.Background())

	w.WriteHeader(http.StatusOK)
}

// DebatePipelineGenerateRequest is used to handle the request for a debate style thought process.
func (d *DebateofModels) DebatePipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	api.Cors(w, req)

	var payload debateRequest

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	promptChoice, maxtokens := processPrompt(payload.Mode)

	// generate a response
	result, err := d.Generate(payload.Prompt, promptChoice, maxtokens)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error getting generation from model"))
		go d.Shutdown(ctx)
		return
	}

	go d.Shutdown(ctx)

	// for debugging streaming
	color.Green(result)

	respPayload := debateResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error marshaling your response from model"))
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
		color.Red("%s", err)
		return err
	}
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	for i := 0; i < len(d.Models)-1; i++ {
		createResponse, err := CreateContainer(d.DockerClient, "800"+strconv.Itoa(i), "", ctx, d.Models[i], d.ContainerImage, d.GPU)
		if err != nil {
			color.Yellow("%s", createResponse.Warnings)
			color.Red("%s", err)
			return err
		}

		color.Green("%s", createResponse.ID)
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
				color.Red("%s", err)
				return "", err
			}
			color.Green("Starting container %d...", i)

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

			color.Yellow("Debug: %s%s", systemprompt, prompt)

			// get reponse
			param := openai.ChatCompletionNewParams{
				Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(systemprompt + result),
					openai.UserMessage(prompt),
				}),
				Seed:        openai.Int(0),
				Model:       openai.String(d.Models[i]),
				Temperature: openai.Float(vars.ModelTemperature),
				MaxTokens:   openai.Int(maxtokens),
			}

			result, err = GenerateCompletion(param, "", *openaiClient)
			if err != nil {
				color.Red("%s", err)
				return "", err
			}

			systemprompt = systemprompt + "\nAnswer from previous expert: " + result

			color.Green("Stopping container %d...", i)
			(d.DockerClient).ContainerStop(context.Background(), model.ID, container.StopOptions{})
		}
	}

	return result, nil
}

func (d *DebateofModels) Shutdown(ctx context.Context) {
	// turn off the containers if they aren't already off
	for i := range d.Models {
		(d.DockerClient).ContainerStop(ctx, d.Models[i], container.StopOptions{})
	}

	// remove the containers
	for i := range d.Models {
		(d.DockerClient).ContainerRemove(ctx, d.Models[i], container.RemoveOptions{})
	}

	color.Green("Shutting Down...")
}

func (d *DebateofModels) PickImage() {
	gpuTrue := IsGPU()
	if gpuTrue {
		gpus, err := GatherGPUs()
		if err != nil {
			d.ContainerImage = vars.CpuImage
			return
		}
		for _, gpu := range gpus {
			if gpu.DeviceInfo.Vendor.Name == "NVIDIA Corporation" {
				d.ContainerImage = vars.CudagpuImage
				break
			}

			// BUG(v,t): fix idk what the value is.
			// After reading upstream, he reads the devices mounted
			// with $ ll /sys/class/drm/
			if gpu.DeviceInfo.Vendor.Name == "Advanced Micro Devices, Inc. [AMD/ATI]" {
				d.ContainerImage = vars.RocmgpuImage
				break
			}
		}
	} else {
		d.ContainerImage = vars.CpuImage
	}

	fmt.Println(d.ContainerImage)
}
