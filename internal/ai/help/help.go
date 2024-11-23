// internal/ai/help/help.go
package help

import (
    "fmt"
    
    "github.com/yourorg/yourtools/internal/ai/models"
)

const (
    usage = `AI CLI - Simple interface for Ollama models

Usage:
  ai [flags] < prompt`

    flags = `
Flags:
  -m, --model string     Model name to use (default "%s")
  -i, --image string     Path to image file (for vision models)
  -s, --system string    System prompt to set context/behavior
  -h, --help            Show help message`

    examples = `
Examples:
  # Basic usage
  echo "How are you?" | ai
  echo "How are you?" | ai -m llama2

  # With system prompt
  echo "Analyze this code" | ai -s "You are an expert Go programmer"
  
  # Image analysis (automatically uses %s)
  echo "Describe this image" | ai -i ./photo.jpg
  echo "Describe this image" | ai -m llava -i ./photo.jpg -s "Focus on technical details"`

    models = `
Available Models:
  Text: llama2, codellama, mistral (default: %s)
  Vision: llava (default: %s)`
)

// ... rest remains the same ...

