/*
Package SLaPE is a binary that orchestrates containers using docker on the local computer.

Usage:

	./slape

Containerized models are spawned as needed adhering to a pipeline system.
*/

// go:generate go tool swagger generate spec -o ./swagger/swagger.yml --scan-models -c ./pkg --exclude-dep
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StoneG24/slape/pkg/api"
	"github.com/StoneG24/slape/pkg/pipeline"
	"github.com/fatih/color"
)

var (
	s = pipeline.SimplePipeline{
		// updates after created
		Model:          "",
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: pipeline.PickImage(),
		DockerClient:   nil,
		GPU:            pipeline.IsGPU(),
	}

	c = pipeline.ChainofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: pipeline.PickImage(),
		DockerClient:   nil,
		GPU:            pipeline.IsGPU(),
	}

	d = pipeline.DebateofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: pipeline.PickImage(),
		DockerClient:   nil,
		GPU:            pipeline.IsGPU(),
	}

	e = pipeline.EmbeddingPipeline{
		// updates after created
		DockerClient:   nil,
		ContainerImage: pipeline.PickImage(),
		GPU:            pipeline.IsGPU(),
	}
)

// @title My API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /
func main() {

	// channel for managing pipelines
	// keystone := make(chan pipeline.Pipeline)

	http.HandleFunc("POST /simple/generate", s.SimplePipelineGenerateRequest)
	http.HandleFunc("POST /simple/setup", s.SimplePipelineSetupRequest)
	http.HandleFunc("GET /simple/shutdown", s.Shutdown)
	http.HandleFunc("POST /cot/generate", c.ChainPipelineGenerateRequest)
	http.HandleFunc("POST /cot/setup", c.ChainPipelineSetupRequest)
	http.HandleFunc("GET /cot/shutdown", c.Shutdown)
	http.HandleFunc("POST /deb/setup", d.DebatePipelineSetupRequest)
	http.HandleFunc("POST /deb/generate", d.DebatePipelineGenerateRequest)
	http.HandleFunc("GET /deb/shutdown", d.Shutdown)
	http.HandleFunc("GET  /emb/setup", e.EmbeddingPipelineSetupRequest)
	http.HandleFunc("POST /emb/generate", e.EmbeddingPipelineGenerateRequest)
	http.HandleFunc("GET /emb/shutdown", e.Shutdown)
	//http.HandleFunc("/moe", simplerequest)
	//http.HandleFunc("/up", upDog)
	http.HandleFunc("GET /getmodels", api.GetModels)

	// Create a new HTTP server.
	srv := &http.Server{
		Addr: ":8080",
	}

	// Start the server in a goroutine.
	color.Green("[+] Server started on :8080")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// Create a channel to listen for interrupt signals.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Block until a signal is received.
	<-sigChan

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %s", err)
	}

	// Close the pipeline to stop adding new pipelines
	// close(keystone)

	s.Shutdown(nil, nil)
	c.Shutdown(nil, nil)
	d.Shutdown(nil, nil)

	log.Println("Server gracefully stopped")
}
