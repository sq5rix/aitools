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
        assistPrompt = flag.String("assistant", "", "System prompt")
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

    // Initialize chat session
    chatSession := aiClient.NewChatSession(selectedModel, *systemPrompt, *assistPrompt)

    // Process initial prompt
    var response string
    var err error

    if *imageFile != "" {
        fmt.Println("Processing image prompt...")
        response, err = chatSession.SendWithImage(userPrompt, *imageFile)
    } else {
        fmt.Println("Processing text prompt...")
        response, err = chatSession.Send(userPrompt)
    }

    if err != nil {
        fmt.Printf("Error in chat: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("\n%s: %s", selectedModel, response)

    // Enter chat loop if not in pipe mode
    fmt.Println("\nEnter your message (type 'exit' to quit):")
    scanner := bufio.NewScanner(os.Stdin)
    
    for {
        fmt.Print("\nYou: ")
        if !scanner.Scan() {
            break
        }
        
        input := strings.TrimSpace(scanner.Text())
        if input == "exit" {
            fmt.Println("Ending chat session.")
            break
        }
        
        if input == "" {
            continue
        }

        if *debug {
            fmt.Printf("Debug: Sending message: %s\n", input)
        }

        response, err = chatSession.Send(input)
        if err != nil {
            fmt.Printf("Error in chat: %v\n", err)
            continue
        }

        fmt.Printf("\n%s: %s\n", selectedModel, response)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading input: %v\n", err)
        os.Exit(1)
    }
}

