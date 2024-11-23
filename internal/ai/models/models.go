// internal/ai/help/help.go
package help

const (
    usage = `AI CLI - Simple interface for Ollama models

Usage:
  ai [flags] < prompt`

    flags = `
Flags:
  -m, --model string        Model name to use (default "%s")
  -i, --image string        Path to image file (for vision models)
  -s, --system string       System prompt to set context/behavior
  -a, --assist string       Assistant prompt to set role/personality
  -t, --temperature float   Temperature (0-1, default: 0.7, lower = more focused)
  -p, --top-p float        Top P sampling (0-1, default: 0.9)
  -h, --help               Show help message`

    examples = `
Examples:
  # Basic usage
  echo "How are you?" | ai
  
  # With system and assistant prompts
  echo "Review this code" | ai -s "You are an expert in Go" -a "I am a helpful code reviewer"
  
  # Adjust generation parameters
  echo "Generate ideas" | ai -t 0.9 -p 0.95    # More creative
  echo "Solve this math" | ai -t 0.2 -p 0.8    # More focused
  
  # Image analysis
  echo "Describe this image" | ai -i ./photo.jpg -a "I am an art critic"
  
  # Complex example
  echo "Analyze architecture" | ai -m llava -i diagram.png \
    -s "Focus on scalability" \
    -a "I am a senior cloud architect" \
    -t 0.4 -p 0.9`

    models = `
Available Models:
  Text: llama2, codellama, mistral (default: %s)
  Vision: llava (default: %s)

Parameters:
  Temperature: Controls randomness (0-1)
    0.2: More focused, deterministic
    0.7: Balanced (default)
    0.9: More creative, varied

  Top P: Controls sampling (0-1)
    0.5: More focused sampling
    0.9: Balanced (default)
    0.95: More varied sampling`
)

