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
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StoneG24/slape/pkg/api"
	"github.com/StoneG24/slape/pkg/pipeline"
)

var (
	s = pipeline.SimplePipeline{
		// updates after created
		Models:         "",
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

	// This is info by default
	var programLevel = new(slog.LevelVar)
	// Change to Debug so we get debug logs
	programLevel.Set(slog.LevelDebug)

	// channel for managing pipelines
	// keystone := make(chan pipeline.Pipeline)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /simple/generate", s.SimplePipelineGenerateRequest)
	mux.HandleFunc("POST /simple/setup", s.SimplePipelineSetupRequest)
	mux.HandleFunc("GET /simple/shutdown", s.Shutdown)
	mux.HandleFunc("POST /cot/generate", c.ChainPipelineGenerateRequest)
	mux.HandleFunc("POST /cot/setup", c.ChainPipelineSetupRequest)
	mux.HandleFunc("GET /cot/shutdown", c.Shutdown)
	mux.HandleFunc("POST /deb/setup", d.DebatePipelineSetupRequest)
	mux.HandleFunc("POST /deb/generate", d.DebatePipelineGenerateRequest)
	mux.HandleFunc("GET /deb/shutdown", d.Shutdown)
	mux.HandleFunc("GET  /emb/setup", e.EmbeddingPipelineSetupRequest)
	mux.HandleFunc("POST /emb/generate", e.EmbeddingPipelineGenerateRequest)
	mux.HandleFunc("GET /emb/shutdown", e.Shutdown)
	//mux.HandleFunc("/moe", simplerequest)
	//mux.HandleFunc("/up", upDog)
	mux.HandleFunc("GET /getmodels", api.GetModels)

	// This is against my religion
	wrappingMux := NewCoors(mux)

	// Create a new HTTP server.
	srv := &http.Server{
		Addr:    ":8080",
		Handler: wrappingMux,
	}

	// Start the server in a goroutine.
	slog.Info("[+] Server started on :8080")
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

	slog.Info("[+] Server gracefully stopped")
}

type Coors struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (c *Coors) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//start := time.Now()

	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
		return
	}

	c.handler.ServeHTTP(w, req)
	//log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
}

// NewCoors constructs a new Logger middleware handler
func NewCoors(handlerToWrap http.Handler) *Coors {
	return &Coors{handlerToWrap}
}

// TODO rewrite to model middleware
func Cors(w http.ResponseWriter, req *http.Request) {
}
