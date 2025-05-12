package main

import "github.com/akishichinibu/goenum/internal/gen"

func main() {
	if err := gen.Gen("."); err != nil {
		panic(err)
	}
}
