package pipeline

import (
	"fmt"
)

// ContextBox is a struct that contains a
// group of strings that contains context on a given problem.
// This is coupled with the system prompt chosen is what makes the models understand
// the gven situation more.

// This information should be kept within a pipeline for privacy and safety reasons.
type ContextBox struct {
	SystemPrompt          string
	Prompt                string
	ConversationHistory   *[]string
	FutureQuestions       string
	InternetSearchResults *[]string
	ToolResults           *[]string
	//VectorStore           vectorstore.VectorStore{}
}

// PromptBuilder takes the ContextBox and builds the system prompt
func (c *ContextBox) PromptBuilder(previousAnswer string) error {

	// TODO(v,t) Go and gather the additional context from

	// TODO(v) vector store
	context := ""

	// minirag
	rag := ""

	additionalContex := context + rag

	c.SystemPrompt = fmt.Sprintf(c.SystemPrompt, additionalContex, previousAnswer)

	// TODO(v) do something differnt for debate where we have question/idea and ask the hats after.
	return nil
}
