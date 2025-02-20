package pipeline

// ContextBox is a struct that contains a
// vector store of information usable by models.
type ContextBox struct {
	ConversationHistory   *[]string
	FutureQuestions       *[]string
	InternetSearchResults *[]string
	//Embeddings            openai.VectorStore
}
