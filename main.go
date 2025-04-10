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
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StoneG24/slape/internal/logging"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/StoneG24/slape/pkg/pipeline"
)

var (
    isGPU = pipeline.IsGPU()
    image = pipeline.PickImage()

	s = pipeline.SimplePipeline{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: image,
		DockerClient:   nil,
		GPU:            isGPU,
	}

	c = pipeline.ChainofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: image,
		DockerClient:   nil,
		GPU:            isGPU,
	}

	d = pipeline.DebateofModels{
		// updates after created
		Models:         []string{},
		ContextBox:     pipeline.ContextBox{},
		Tools:          pipeline.Tools{},
		Active:         true,
		ContainerImage: image,
		DockerClient:   nil,
		GPU:            isGPU,
	}

	e = pipeline.EmbeddingPipeline{
		// updates after created
		DockerClient:   nil,
		ContainerImage: image,
		GPU:            isGPU,
	}
)

// @title My API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /
func main() {

	// Change to Debug so we get debug logs
	slog.SetLogLoggerLevel(slog.LevelDebug)

	go logging.CreateLogFile()

	// channel for managing pipelines
	// keystone := make(chan pipeline.Pipeline)

	// Default Mux for our server.
	// For auth in the future we will want to setup a different set.
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

	err := shutdownPipelines()
	if err != nil {
		slog.Error("Error", "ErrorString", err)
	}
	/*
	   s.DockerClient.Close()
	   c.DockerClient.Close()
	   d.DockerClient.Close()
	   e.DockerClient.Close()
	*/

	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %s", err)
	}

	// Close the pipeline to stop adding new pipelines
	// close(keystone)

	slog.Info("[+] Server gracefully stopped")
}

type Coors struct {
	handler http.Handler
}

// ServeHTTP handles the request by passing it to the real
// handler and logging the request details
func (c *Coors) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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

// need to shutdown pipelines with
// a remote request since the shutdown functions now need a req struct
func shutdownPipelines() error {

	url := "http://localhost:8080/%s/shutdown"
	pipelines := []string{"simple", "cot", "deb", "emb"}

	for _, pipeline := range pipelines {
		requrl := fmt.Sprintf(url, pipeline)
		resp, err := http.Get(requrl)
		if err != nil {
			return err
		}
		resp.Body.Close()
	}

	return nil
}
