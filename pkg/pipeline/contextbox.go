package pipeline

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/StoneG24/slape/internal/vars"
	"github.com/StoneG24/slape/pkg/api"
	"github.com/openai/openai-go"
)

// ContextBox is a struct that contains a
// group of strings that contains context on a given problem.
// This is coupled with the system prompt chosen is what makes the models understand
// the gven situation more.

// This information should be kept within a pipeline for privacy and safety reasons.
type ContextBox struct {
	SystemPrompt          string
	Thoughts              string
	Prompt                string
	ConversationHistory   *[]string
	FutureQuestions       string
	InternetSearchResults *[]string
	ToolResults           *[]string
	//VectorStore           vectorstore.VectorStore{}
}

// PromptBuilder takes the ContextBox and builds the system prompt
func (c *ContextBox) PromptBuilder(previousAnswer string) error {

	// since we are operating on a parameter its
	// safer to create a local copy
	prevAns := previousAnswer

	// TODO(v,t) Go and gather the additional context from

	// TODO(v) vector store
	context := ""

	// minirag
	rag := ""

	// information generated as prelinary thoughts
	var additionalContex string
	if len(context) != 0 && len(rag) != 0 {
		additionalContex = context + rag
	} else {
		additionalContex = "None"
	}

	if len(previousAnswer) == 0 {
		prevAns = "None"
	}

	slog.Debug(c.Thoughts, additionalContex, prevAns)
	c.SystemPrompt = fmt.Sprintf(c.SystemPrompt, c.Thoughts, additionalContex, prevAns)

	// TODO(v) do something different for debate where we have question/idea and ask the hats after.
	return nil
}

// getThought is used to generate initial thoughts about a given question.
// This is supposed to create some guardrails for thought.
// This will not be good for slms but llms that are centered around reasoning
func (c *ContextBox) getThoughts(ctx context.Context) {

	prompt := `
    You are tasked with solving a problem. Start by carefully considering and listing all the known facts surrounding the scenario. What do you already know about the situation? What information is available to you?
    Next, identify the constraints based on these facts. What limitations or conditions must you take into account when approaching the problem? Consider factors like time, resources, and external influences that may affect the solution.
    Once you’ve fully considered the facts and constraints, generate potential solutions to the problem. Think creatively and strategically, taking into account the constraints you’ve identified. Focus on generating ideas that are practical, feasible, and innovative. Provide a rationale for each idea, considering how well it aligns with the constraints and solves the problem at hand.
    `

	param := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(prompt),
			openai.UserMessage(c.Prompt),
			//openai.UserMessage(s.FutureQuestions),
		}),
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

	result, err := GenerateCompletion(ctx, param, "", *vars.OpenaiClient)
	log.Println(result)
	if err != nil {
		c.Thoughts = "None"
	}

	slog.Debug("Debug", "DebugValue", result)

	c.Thoughts = result
}
