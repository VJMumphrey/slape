package vars

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	CpuImage = "ghcr.io/ggml-org/llama.cpp:server"

	CudagpuImage = "ghcr.io/ggml-org/llama.cpp:server-cuda"
	RocmgpuImage = "ghcr.io/ggml-org/llama.cpp:server-rocm"
)

var (
	OpenaiClient = openai.NewClient(
		option.WithBaseURL("http://localhost:8000/v1"),
	)

	GenerationClient = openai.NewClient(
		option.WithBaseURL("http://localhost:8081/v1"),
	)

	EmbeddingClient = openai.NewClient(
		option.WithBaseURL("http://localhost:8082/v1"),
	)

	ModelTemperature = 0.1
)
