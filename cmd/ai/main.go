package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/sq5rix/aitools/internal/ai/client"
    "github.com/sq5rix/aitools/internal/ai/help"
)

func main() {
    var (
        systemPrompt = flag.String("system", "", "System prompt")
        imageFile    = flag.String("image", "", "Path to image file")
        modelName    = flag.String("model", "llama2", "Model name to use")
        showHelp     = flag.Bool("help", false, "Show help")
        listModels   = flag.Bool("list", false, "List available models")
    )
    flag.Parse()

    if *showHelp {
        fmt.Println(help.GetHelp())
        return
    }

    aiClient := client.NewOllamaClient("http://localhost:11434")

    if *listModels {
        models, err := aiClient.ListModels()
        if err != nil {
            fmt.Printf("Error listing models: %v\n", err)
            os.Exit(1)
        }
        fmt.Println("Available models:")
        for _, model := range models {
            fmt.Printf("- %s\n", model)
        }
        return
    }

    // Read user input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your prompt: ")
    userPrompt, _ := reader.ReadString('\n')
    userPrompt = strings.TrimSpace(userPrompt)

    if userPrompt == "" {
        fmt.Println(help.GetErrorNoPrompt())
        return
    }

    var response string
    var err error

    if *imageFile != "" {
        response, err = aiClient.GenerateWithImage(userPrompt, *systemPrompt, *modelName, *imageFile)
    } else {
        response, err = aiClient.Generate(userPrompt, *systemPrompt, *modelName)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(response)
}

