/*
Package SLaPE is a binary that orchestrates containers using docker on the local computer.

Usage:

	./slape

Containerized models are spawned as needed adhering to a pipeline system.
*/
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
		ContainerImage: "",
		DockerClient:   nil,
		GPU:            false,
	}

	c = pipeline.ChainofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: "",
		DockerClient:   nil,
		GPU:            false,
	}

	d = pipeline.DebateofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: "",
		DockerClient:   nil,
		GPU:            false,
	}
)

func main() {

	// channel for managing pipelines
	// keystone := make(chan pipeline.Pipeline)

	http.HandleFunc("/simple", s.SimplePipelineGenerateRequest)
	http.HandleFunc("/smplsetup", s.SimplePipelineSetupRequest)
	http.HandleFunc("/cot", c.ChainPipelineGenerateRequest)
	http.HandleFunc("/cotsetup", c.ChainPipelineSetupRequest)
	http.HandleFunc("/debate", d.DebatePipelineGenerateRequest)
	http.HandleFunc("/debsetup", d.DebatePipelineSetupRequest)
	//http.HandleFunc("/moe", simplerequest)
	//http.HandleFunc("/up", upDog)
	http.HandleFunc("/getmodels", api.GetModels)

	// Create a new HTTP server.
	srv := &http.Server{
		Addr: ":3069",
	}

	// Start the server in a goroutine.
	color.Green("[+] Server started on :3069")
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
