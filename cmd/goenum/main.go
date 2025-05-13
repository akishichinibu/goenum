package main

import (
	"os"

	"github.com/akishichinibu/goenum/internal/gen"
)

func main() {
	workdir := os.Args[1]
	if err := gen.Gen(workdir); err != nil {
		gen.Logger.Error(err.Error())
		os.Exit(1)
	}
}
