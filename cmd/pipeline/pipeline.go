// Package pipline is meant to help facilitate the running of several models in sequential order.
//
// For machines with small amounts of resources the pipelines will have to manage models by orchestrating models as need be.
// Pipelines also dictate the structure of conversations between the models. See the following docs for more info.
package pipeline

// ContextBox is a struct that contains a
// vector store of information usable by models.
// The vector store is of openai spec.
type ContextBox struct{}

// SimplePipeline is the smallest pipeline.
// It contains only a model with a ContextBox.
type SimplePipeline struct{}

// InitSimplePipeline creates a SimplePipeline.
func InitSimplePipeline() {}

// ChainofModels is the next step above smallest pipeline.
// This pipeline contains a ContextBox and the models in squential order.
// ChainofModels forces the models to talk in sequential order
// like the name suggests.
type ChainofModels struct{}

// InitChainofModels creates a ChainofModels pipeline.
// Includes a ContextBox and all models needed in squential order.
func InitChainofModels() {}

// DebateofModels is pipeline for debate structured prompting.
// Models talk in a round robin style.
type DebateofModels struct{}

// InitDebateofModels creates a DebateofModels pipeline for debates.
// Includes a ContextBox and all models needed.
func InitDebateofModels() {}
