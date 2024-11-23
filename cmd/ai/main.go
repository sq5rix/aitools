package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/sq5rix/aitools/internal/ai/client"
    "github.com/sq5rix/aitools/internal/ai/help"
    "github.com/sq5rix/aitools/internal/ai/models"
)

func main() {
    var (
        systemPrompt = flag.String("system", "", "System prompt")
        imageFile    = flag.String("image", "", "Path to image file")
        modelName    = flag.String("model", "", "Model name to use")
        showHelp     = flag.Bool("help", false, "Show help")
        listModels   = flag.Bool("list", false, "List available models")
        debug        = flag.Bool("debug", false, "Enable debug output")
    )
    flag.Parse()
    // Check if we're getting piped input
    stat, _ := os.Stdin.Stat()
    isPipe := (stat.Mode() & os.ModeCharDevice) == 0
    var userPrompt string


    if isPipe {
        // Read from pipe
        scanner := bufio.NewScanner(os.Stdin)
        var input []string
        for scanner.Scan() {
            input = append(input, scanner.Text())
        }
        userPrompt  = strings.Join(input, "\n")

    } else {
        // Interactive mode
        fmt.Print("Enter your prompt: ")
        scanner := bufio.NewScanner(os.Stdin)
        scanner.Scan()
        userPrompt   = scanner.Text()
    }

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

    userPrompt = strings.TrimSpace(userPrompt)

    if userPrompt == "" {
        fmt.Println(help.GetErrorNoPrompt())
        return
    }

    // Select appropriate model based on input
    selectedModel := *modelName
    if selectedModel == "" {
        if *imageFile != "" {
            selectedModel = models.DefaultVisionModel
        } else {
            selectedModel = models.DefaultTextModel
        }
    }

    if *debug {
        fmt.Printf("Debug: Using model: %s\n", selectedModel)
        fmt.Printf("Debug: User prompt: %s\n", userPrompt)
        fmt.Printf("Debug: System prompt: %s\n", *systemPrompt)
        if *imageFile != "" {
            fmt.Printf("Debug: Image file: %s\n", *imageFile)
        }
    }

    var response string
    var err error

    if *imageFile != "" {
        fmt.Println("Processing image prompt...")
        response, err = aiClient.GenerateWithImage(userPrompt, *systemPrompt, selectedModel, *imageFile)
    } else {
        fmt.Println("Processing text prompt...")
        response, err = aiClient.Generate(userPrompt, *systemPrompt, selectedModel)
    }

    if err != nil {
        fmt.Printf("Error generating response: %v\n", err)
        os.Exit(1)
    }

    if response == "" {
        fmt.Println("Warning: Received empty response from model")
    }

    fmt.Printf("\nResponse from %s:\n%s\n", selectedModel, response)
}

