package pipeline

// ContextBox is a struct that contains a
// group of strings that contains context on a given problem.
// This is coupled with the system prompt chosen is what makes the models understand
// the gven situation more.

// This information should be kept within a pipeline for privacy and safety reasons.
type ContextBox struct {
	ConversationHistory   *[]string
	FutureQuestions       *[]string
	InternetSearchResults *[]string
	ToolResults           *[]string
	//VectorStore           vectorstore.VectorStore{}
}

// PromptBuilder takes the ContextBox and builds the system prompt
func (c *ContextBox) PromptBuilder() (string, error) { return "", nil }
