/*
Package discussion is used to define the discussion component of slape. This is the v2 of pipelines.
It plays the same purpose, but focuses more on the conversation being had by the models and not the models themselves.

It would be good to consider pipeline as a deprecated package that is kept around for code and inspiration.
*/
package discussion

import (
	"github.com/VJMumphrey/slape/pkg/vectorstore"
)

type Discussion struct {
	ConversationHistory
	Tools
	vs VectorStore
}
