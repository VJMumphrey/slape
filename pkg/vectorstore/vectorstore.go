/*
Package vectorstore is used to create independant context information for each disscussion.
Vectorstores allow us to index information and keep relevant bits for generation.
*/
package vectorstore

import (
	"fmt"

	"github.com/coder/hnsw"
)

// Utilize HNSW algorithim to create a vector for Internet searching
type VectoreStore struct {
}

// This is a simple function to get us going and will need more work in future to implement
// correctly. Apparently this has some openai compat
func (v *VectoreStore) Startup() {
	g := hnsw.NewGraph[int]()
	g.Add(
		hnsw.MakeNode(1, []float32{1, 1, 1}),
		hnsw.MakeNode(2, []float32{1, -1, 0.999}),
		hnsw.MakeNode(3, []float32{1, 0, -0.5}),
	)

	// to get this vector we need to tokenize some text
	// then embed it
	neighbors := g.Search(
		[]float32{0.5, 0.5, 0.5},
		1,
	)
	fmt.Printf("best friend: %v\n", neighbors[0].Value)
	// Output: best friend: [1 1 1]
}
