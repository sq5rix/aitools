# AI Tools

A command-line interface (CLI) tool for interacting with Ollama AI models, supporting both text and image-based prompts.

## Features

- Text generation with AI models
- Image-based prompts (vision models)
- System and assistant prompts support
- Pipe input support
- Model selection
- Debug mode
- Interactive and non-interactive modes

## Prerequisites

- Go 1.20 or higher
- [Ollama](https://ollama.ai/) running locally (default: http://localhost:11434)

## Installation

```bash
go install github.com/sq5rix/aitools@latest
```

## Usage

### Basic Usage

```bash
# Interactive mode
ai

# With a specific prompt
ai "What is the capital of France?"

# Pipe input
echo "What is the capital of France?" | ai
```

### Command Line Options

```bash
-model string     Model name to use (default: depends on input type)
-system string    System prompt for context
-assistant string Additional assistant instructions
-image string     Path to image file for vision tasks
-help            Show help information
-list            List available models
-debug           Enable debug output
```

### Examples

1. Basic text generation:

```bash
ai "Explain quantum computing"
```

2. Using a specific model:

```bash
ai -model llama2 "Explain quantum computing"
```

3. With system prompt:

```bash
ai -system "You are a helpful assistant" "How do I learn Go programming?"
```

4. Image analysis:

```bash
ai -image path/to/image.jpg "What's in this image?"
```

5. Debug mode:

```bash
ai -debug "Tell me a joke"
```

6. List available models:

```bash
ai -list
```

## Input Modes

### Interactive Mode

When run without piped input, the tool enters interactive mode and prompts for user input.

### Pipe Mode

The tool can accept input through Unix pipes:

```bash
cat prompt.txt | ai
```

## Default Models

- Text Generation: Uses the default text model configured in models package
- Vision Tasks: Uses the default vision model when processing images

## Error Handling

The tool provides clear error messages for common issues:

- Empty prompts
- Model availability
- API connection issues
- Image file access
- Invalid inputs

## Debug Mode

Enable debug mode with `-debug` to see:

- Selected model
- User prompt
- System prompt
- Image file path (if applicable)
- Other processing details

## Project Structure

```
aitools/
├── cmd/
│   └── ai/
│       └── main.go           # Main CLI application
├── internal/
│   └── ai/
│       ├── client/          # Ollama API client
│       ├── help/            # Help text and messages
│       └── models/          # Model configurations
└── README.md
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built with [Ollama](https://ollama.ai/)
- Inspired by the need for simple AI interaction tools

## Support

For support, please open an issue in the GitHub repository.
