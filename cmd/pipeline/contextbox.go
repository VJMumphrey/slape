package pipeline

// ContextBox is a struct that contains a 
// group of strings that contains context on a given problem.
// This information should be kept within a pipeline for privacy and safety reasons.
type ContextBox struct {
	// This is maybe and probably is kept in the convo history anyways
	ConversationHistory   *[]string
	FutureQuestions       *[]string
	InternetSearchResults *[]string
	ToolResults           *[]string
	//Embeddings            openai.VectorStore
}
