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
	"golang.org/x/net/websocket"
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
)

// @title My API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /
func main() {

	// channel for managing pipelines
	// keystone := make(chan pipeline.Pipeline)

	http.HandleFunc("POST /simple", websocket.Handler(s.SimplePipelineGenerateRequest))
	http.HandleFunc("POST /smplsetup", s.SimplePipelineSetupRequest)
	http.HandleFunc("POST /cot", c.ChainPipelineGenerateRequest)
	http.HandleFunc("POST /cotsetup", c.ChainPipelineSetupRequest)
	http.HandleFunc("POST /debate", d.DebatePipelineGenerateRequest)
	http.HandleFunc("POST /debsetup", d.DebatePipelineSetupRequest)
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

	go s.Shutdown(context.Background())
	go c.Shutdown(context.Background())
	go d.Shutdown(context.Background())

	log.Println("Server gracefully stopped")
}
