package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/StoneG24/slape/cmd/api"
	"github.com/StoneG24/slape/cmd/vars"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/openai/openai-go"
)

// SimplePipeline is the smallest pipeline.
// It contains only a model with a ContextBox.
// This is good for a giving a single model access to tools
// like internet search.
type SimplePipeline struct {
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

type simpleRequest struct {
	// Prompt is the string that
	// will be appended to the prompt
	// string chosen.
	Prompt string `json:"prompt"`
	Model  string `json:"model"`

	// Options are strings matching
	// the names of prompt types
	Mode string `json:"mode"`
}

type simpleResponse struct {
	Answer string `json:"answer"`
}

// simplerequest is used to handle simple requests as needed.
func (s *SimplePipeline) SimplePipelineRequest(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		color.Red("%s", err)
		return
	}
	defer apiClient.Close()

	go api.Cors(w, req)

	w.Header().Set("Content-Type", "application/json")

	var simplePayload simpleRequest

	err = json.NewDecoder(req.Body).Decode(&simplePayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Error unexpected request format"))
		return
	}

	s.PickImage()

	s.Model = simplePayload.Model
	s.DockerClient = apiClient
	s.GPU = IsGPU()

	go s.Setup(ctx)

	promptChoice, maxtokens := processPrompt(simplePayload.Mode)

	// generate a response
	result, err := s.Generate(simplePayload.Prompt, promptChoice, maxtokens, vars.OpenaiClient)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error getting generation from model"))
		err := s.Shutdown(ctx)
		if err != nil {
			color.Red("%s", err)
		}
		return
	}

	go s.Shutdown(ctx)

	// for debugging streaming
	color.Green(result)

	respPayload := simpleResponse{
		Answer: result,
	}

	json, err := json.Marshal(respPayload)
	if err != nil {
		color.Red("%s", err)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Error marshaling your response from model"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func (s *SimplePipeline) Setup(ctx context.Context) error {

	reader, err := PullImage(s.DockerClient, ctx, s.ContainerImage)
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
		s.DockerClient,
		"8000",
		"",
		ctx,
		s.Model,
		s.ContainerImage,
		s.GPU,
	)

	if err != nil {
		color.Yellow("%s", createResponse.Warnings)
		color.Red("%s", err)
		return err
	}

	// start container
	err = (s.DockerClient).ContainerStart(context.Background(), createResponse.ID, container.StartOptions{})
	if err != nil {
		color.Red("%s", err)
		return err
	}

	// For debugging
	color.Green("%s", createResponse.ID)
	s.container = createResponse

	return nil

}

func (s *SimplePipeline) Generate(prompt string, systemprompt string, maxtokens int64, openaiClient *openai.Client) (string, error) {
	// take care of upDog on our own
	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	color.Yellow("Debug: %s%s", systemprompt, prompt)

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemprompt),
			openai.UserMessage(prompt),
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

func (s *SimplePipeline) Shutdown(ctx context.Context) error {
	err := (s.DockerClient).ContainerStop(ctx, s.container.ID, container.StopOptions{})
	if err != nil {
		color.Red("%s", err)
        return  nil
	}

	err = (s.DockerClient).ContainerRemove(ctx, s.container.ID, container.RemoveOptions{})
	if err != nil {
		color.Red("%s", err)
        return  nil
	}

	color.Green("Shutting Down...")

	return nil
}

func (s *SimplePipeline) PickImage() {
	gpuTrue := IsGPU()
	if gpuTrue {
		gpus, err := GatherGPUs()
		if err != nil {
			s.ContainerImage = vars.CpuImage
			return
		}
		for _, gpu := range gpus {
			if gpu.DeviceInfo.Vendor.Name == "NVIDIA Corporation" {
				s.ContainerImage = vars.CudagpuImage
				break
			}

			// BUG(v,t): fix idk what the value is.
			// After reading upstream, he reads the devices mounted
			// with $ ll /sys/class/drm/
			if gpu.DeviceInfo.Vendor.Name == "Advanced Micro Devices, Inc. [AMD/ATI]" {
				s.ContainerImage = vars.RocmgpuImage
				break
			}
		}
	} else {
		s.ContainerImage = vars.CpuImage
	}

	fmt.Println(s.ContainerImage)
}
