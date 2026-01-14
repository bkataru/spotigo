// Package main provides the entry point for the Spotigo CLI application.
// Spotigo is a command-line interface for interacting with Spotify and RAG-based music analysis.
package main

import (
	"os"

	"github.com/bkataru/spotigo/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
