package main

import (
	"os"

	"github.com/kmrok/cfn-diff/internal/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
