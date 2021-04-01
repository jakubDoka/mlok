package memory

import (
	"math/rand"
	"testing"

	"github.com/jakubDoka/mlok/logic/memory/gen"
)

func Test(t *testing.T) {
	vec := make(gen.Int32Vec, 100000)
	for i := range vec {
		vec[i] = rand.Int31() % 10
	}
	vec.Sort(func(a, b int32) bool {
		return a > b
	})
	t.Error(vec)
}
