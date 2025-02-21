// Package pipline is meant to help facilitate the running of several models in sequential order.
//
// For machines with small amounts of resources the pipelines will have to manage models by orchestrating models as need be.
// Pipelines also dictate the structure of conversations between the models. See the following docs for more info.
package pipeline

// SimplePipeline is the smallest pipeline.
// It contains only a model with a ContextBox.
// This is useful giving the model access to tools.
// like internet search
type SimplePipeline struct {
	Model string
	ContextBox
	Tools
}

// InitSimplePipeline creates a SimplePipeline.
func InitSimplePipeline() {}

// ChainofModels is the next step above smallest pipeline.
// This pipeline contains a ContextBox and the models in squential order.
// ChainofModels forces the models to talk in sequential order
// like the name suggests.
type ChainofModels struct {
	Model1 string
	Model2 string
	Model3 string
	ContextBox
	Tools
}

// InitChainofModels creates a ChainofModels pipeline.
// Includes a ContextBox and all models needed in squential order.
func InitChainofModels() {}

// DebateofModels is pipeline for debate structured prompting.
// Models talk in a round robin style.
type DebateofModels struct{
    Models []string
    ContextBox
    Tools
}

// InitDebateofModels creates a DebateofModels pipeline for debates.
// Includes a ContextBox and all models needed.
func InitDebateofModels() {}
