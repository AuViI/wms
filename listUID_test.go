package main

import (
	"fmt"
	"testing"
)

func TestUIDGen(t *testing.T) {
	fmt.Println(getRUID(8))
	fmt.Println(getRUID(8))
	fmt.Println(getRUID(10))
}

func BenchmarkUIDGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getRUID(8)
	}
}
