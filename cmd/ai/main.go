// cmd/ai/main.go
package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "os"

    "github.com/spf13/pflag"
    "github.com/yourorg/yourtools/internal/ai/client"
    "github.com/yourorg/yourtools/internal/ai/help"
)

func main() {
    // Using pflag to support both POSIX-style and GNU-style flags
    var (
        modelName string
        imagePath string
        showHelp  bool
    )

    // Define flags with both short and long versions
    pflag.StringVarP(&modelName, "model", "m", "", "Model name to use")
    pflag.StringVarP(&imagePath, "image", "i", "", "Path to image file (for vision models)")
    pflag.BoolVarP(&showHelp, "help", "h", false, "Show help message")

    // Override default usage message
    pflag.Usage = func() {
        fmt.Fprint(os.Stderr, help.GetHelp())
    }

    pflag.Parse()

    // Show help if requested or if no model specified
    if showHelp || modelName == "" {
        pflag.Usage()
        os.Exit(0)
    }

    // Check if there's input on stdin
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        fmt.Println(help.GetErrorNoPrompt())
        os.Exit(1)
    }

    // Read prompt from stdin
    reader := bufio.NewReader(os.Stdin)
    prompt, err := reader.ReadString('\n')
    if err != nil && err != io.EOF {
        log.Fatalf("Error reading prompt: %v", err)
    }

    aiClient := client.NewOllamaClient("http://localhost:11434")

    if imagePath != "" {
        response, err := aiClient.GenerateWithImage(modelName, prompt, imagePath)
        if err != nil {
            log.Fatalf("Error generating response: %v", err)
        }
        fmt.Print(response)
    } else {
        response, err := aiClient.Generate(modelName, prompt)
        if err != nil {
            log.Fatalf("Error generating response: %v", err)
        }
        fmt.Print(response)
    }
}

