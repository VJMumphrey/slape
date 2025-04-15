package pipeline

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/StoneG24/slape/pkg/api"
	"github.com/StoneG24/slape/pkg/vars"
	"github.com/coder/hnsw"
	"github.com/openai/openai-go"
)

// ContextBox is a struct that contains a
// group of strings that contains context on a given problem.
// This is coupled with the system prompt chosen is what makes the models understand
// the gven situation more.

// This information should be kept within a pipeline for privacy and safety reasons.
type ContextBox struct {
	// Simple prompt components
	SystemPrompt string
	Thoughts     string
	Prompt       string
	// Currently not in use
	ConversationHistory *[]string
	FutureQuestions     string

	// These will come from the internet search package.
	InternetSearchResults *[]string
	VectorStore           *hnsw.Graph[string]

	// These will come from tool calls
	ToolResults *[]string
}

// PromptBuilder takes the ContextBox and builds the system prompt
func (c *ContextBox) promptBuilder(previousAnswer string) error {

	// since we are operating on a parameter its
	// safer to create a local copy
	prevAns := previousAnswer
	if len(previousAnswer) == 0 {
		prevAns = "None"
	}

	// information generated as prelinary thoughts
	// TODO(v) move to generation functions like thoughts
	var additionalContex string
	if c.InternetSearchResults != nil {
		for _, result := range *c.InternetSearchResults {
			additionalContex += result
		}
	} else {
		additionalContex = "None"
	}

	log.Println(c.Thoughts, additionalContex, prevAns)
	c.SystemPrompt = fmt.Sprintf(c.SystemPrompt, c.Thoughts, additionalContex, prevAns)

	// TODO(v) do something different for debate where we have question/idea and ask the hats after.
	return nil
}

// getThought is used to generate initial thoughts about a given question.
// This is supposed to create some guardrails for thought.
// This will not be good for slms but llms that are centered around reasoning
func (c *ContextBox) getThoughts(ctx context.Context) {

	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(vars.ThinkingPrompt),
			openai.UserMessage(c.Prompt),
			//openai.UserMessage(s.FutureQuestions),
		},
		Seed: openai.Int(0),
		//Model:       openai.String(pipeline.Model),
		Temperature: openai.Float(0.4),
		MaxTokens:   openai.Int(16348),
	}

	for {
		// sleep and give server guy a break
		time.Sleep(time.Duration(5 * time.Second))

		// Single model, single port, assuming one pipeline is running at a time
		if api.UpDog("8000") {
			break
		}
	}

	result, err := GenerateCompletion(ctx, param, "", vars.OpenaiClient)
	log.Println(result)
	if err != nil {
		c.Thoughts = "None"
	}

	log.Println("Debug Thinking result", result)

	c.Thoughts = result
}

// getInternetSearch is used to generate initial context about a given question.
func (c *ContextBox) getInternetSearch(ctx context.Context) error {

	// TODO need to convert this to f32 with iterative approach
	embCh := make(chan []float64)
	searchCh := make(chan *hnsw.Graph[string])

	// Generate embedding of prompt
	go func(context.Context, chan []float64) {
		embedparam := openai.EmbeddingNewParams{
			Input:      openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: []string{c.Prompt}},
			Model:      embedmodel,
			Dimensions: openai.Int(1024),
		}

		for {
			// sleep and give server guy a break
			time.Sleep(time.Duration(2 * time.Second))

			// Single model, single port, assuming one pipeline is running at a time
			if api.UpDog("8082") {
				break
			}
		}

		result, err := GenerateEmbedding(ctx, embedparam, vars.EmbeddingClient)
		log.Println(result)
		if err != nil {
			// Have to check length later
			embCh <- []float64{}
		}
		embCh <- result.Data[0].Embedding
	}(ctx, embCh)

	// Generate the query and run internetsearch
	go func(context.Context, chan *hnsw.Graph[string]) {
		// take care of upDog on our own
		for {
			// sleep and give server guy a break
			time.Sleep(time.Duration(1 * time.Second))

			// Single model, single port, assuming one pipeline is running at a time
			if api.UpDog("8000") {
				break
			}
		}

		queryPrompt := `
        Act as a internet guru who knows how to search up anything on duckduckgo.com.
        you are given a request and your job is to create a search query that captures that request to the fullest.
        `

		param := openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(queryPrompt),
				openai.UserMessage(c.Prompt),
				//openai.UserMessage(s.FutureQuestions),
			},
			Seed: openai.Int(0),
			//Model:       s.Models[0],
			Temperature: openai.Float(vars.ModelTemperature),
			MaxTokens:   openai.Int(100),
		}

		result, err := GenerateCompletion(ctx, param, "", vars.OpenaiClient)
		if err != nil {
            log.Println("Error generating query for internet search", err)
		}

        // take the result and run the internetsearch
        graph := 
        

		searchCh <- graph
	}(ctx, searchCh)

	// Combine the two and search the graph for relative neighbors
	var embedding []float64
	var graph *hnsw.Graph[string]
	for i := 0; i < 2; i++ {

		select {
            case embedding = <-embCh:
            case graph = <-searchCh:
		}

	}




	log.Println("Internet Search result", result)

	c.InternetSearchResults = result
}
