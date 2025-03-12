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

	"github.com/StoneG24/slape/pkg/api"
	"github.com/StoneG24/slape/internal/vars"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// ChainofModels is the next step above smallest pipeline.
// This pipeline contains a ContextBox and the models in squential order.
// ChainofModels forces the models to talk in sequential order
// like the name suggests.
type ChainofModels struct {
	Models []string
	ContextBox
	Tools
	Active         bool
	ContainerImage string
	DockerClient   *client.Client
	GPU            bool

	// for internal use to store the models in
	containers []container.CreateResponse
}

type chainRequest struct {
	// Prompt is the string that
	// will be appended to the prompt
	// string chosen.
	Prompt string `json:"prompt"`

	// Options are strings matching
	// the names of prompt types
	Mode string `json:"mode"`
}

type chainSetupPayload struct {
	Models []string `json:"models"`
}

type chainResponse struct {
	Answer string `json:"answer"`
}

// ChainPipeline, handlerfunc expects POST method and returns nothing
func (c *ChainofModels) ChainPipelineSetupRequest(w http.ResponseWriter, req *http.Request) {
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

	var setupPayload chainSetupPayload

	err = json.NewDecoder(req.Body).Decode(&setupPayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	c.PickImage()

	c.Models = setupPayload.Models
	c.DockerClient = apiClient
	c.GPU = IsGPU()

	go c.Setup(context.Background())

	w.WriteHeader(http.StatusOK)
}

// ChainPipelineRequest is used to handle requests for chain of models pipelines.
// The json expected is
// - prompt string, prompt from the user.
// - models array of strings, an array of strings containing three models to use.
// - mode string, mode of prompt struture to use.
func (c *ChainofModels) ChainPipelineGenerateRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	go api.Cors(w, req)

	var payload chainRequest

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	promptChoice, maxtokens := processPrompt(payload.Mode)

	// generate a response
	result, err := c.Generate(payload.Prompt, promptChoice, maxtokens)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error getting generation from model"))
		go c.Shutdown(ctx)
		return
	}

	go c.Shutdown(ctx)

	// for debugging streaming
	color.Green(result)

	respPayload := chainResponse{
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

func (c *ChainofModels) Setup(ctx context.Context) error {
	reader, err := PullImage(c.DockerClient, ctx, c.ContainerImage)
	if err != nil {
		color.Red("%s", err)
		return err
	}
	color.Green("Pulling Image...")
	// prints out the status of the download
	// worth while for big images
	io.Copy(os.Stdout, reader)

	for i, model := range c.Models {
		createResponse, err := CreateContainer(
			c.DockerClient,
			"800"+strconv.Itoa(i),
			"",
			ctx,
			model,
			c.ContainerImage,
			c.GPU,
		)

		if err != nil {
			color.Yellow("%s", createResponse.Warnings)
			color.Red("%s", err)
			return err
		}

		color.Green("%s", createResponse.ID)
		c.containers = append(c.containers, createResponse)
	}

	return nil
}

// ChainofModels.Generate is the facilitator of model orchestration based on the chain of model pipeline.
// Since the pipeline is based on the Chan of Thought prompting technique, it follows this style, mimicing its behavior.
func (c *ChainofModels) Generate(prompt string, systemprompt string, maxtokens int64) (string, error) {
	var result string

	for i, model := range c.containers {
		// start container
		err := (c.DockerClient).ContainerStart(context.Background(), model.ID, container.StartOptions{})
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
			Model:       openai.String(c.Models[i]),
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
		(c.DockerClient).ContainerStop(context.Background(), model.ID, container.StopOptions{})
	}

	return result, nil
}

// ChainofModels.Shutdown handles the shutdown of the pipelines models.
func (c *ChainofModels) Shutdown(ctx context.Context) {
	// turn off the containers if they aren't already off
	for _, model := range c.containers {
		(c.DockerClient).ContainerStop(ctx, model.ID, container.StopOptions{})
	}

	// remove the containers, seperate incase it's already stopped
	for _, model := range c.containers {
		(c.DockerClient).ContainerRemove(ctx, model.ID, container.RemoveOptions{})
	}

	color.Green("Shutting Down...")
}

func (c *ChainofModels) PickImage() {
	gpuTrue := IsGPU()
	if gpuTrue {
		gpus, err := GatherGPUs()
		if err != nil {
			c.ContainerImage = vars.CpuImage
			return
		}
		for _, gpu := range gpus {
			if gpu.DeviceInfo.Vendor.Name == "NVIDIA Corporation" {
				c.ContainerImage = vars.CudagpuImage
				break
			}

			// BUG(v,t): fix idk what the value is.
			// After reading upstream, he reads the devices mounted
			// with $ ll /sys/class/drm/
			if gpu.DeviceInfo.Vendor.Name == "Advanced Micro Devices, Inc. [AMD/ATI]" {
				c.ContainerImage = vars.RocmgpuImage
				break
			}
		}
	} else {
		c.ContainerImage = vars.CpuImage
	}

	fmt.Println(c.ContainerImage)
}
