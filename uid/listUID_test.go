package uid

import (
	"fmt"
	"testing"
)

func TestUIDGen(t *testing.T) {
	fmt.Println(GetRUID(8))
	fmt.Println(GetRUID(8))
	fmt.Println(GetRUID(10))
}

func BenchmarkUIDGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetRUID(8)
	}
}
