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

type config struct {
    systemPrompt  string
    assistPrompt  string
    imageFile     string
    modelName     string
    showHelp      bool
    listModels    bool
    debug         bool
    isPipe        bool
}

func parseFlags() *config {
    cfg := &config{}
    flag.StringVar(&cfg.systemPrompt, "system", "", "System prompt")
    flag.StringVar(&cfg.assistPrompt, "assistant", "", "Assistant prompt")
    flag.StringVar(&cfg.imageFile, "image", "", "Path to image file")
    flag.StringVar(&cfg.modelName, "model", "", "Model name to use")
    flag.BoolVar(&cfg.showHelp, "help", false, "Show help")
    flag.BoolVar(&cfg.listModels, "list", false, "List available models")
    flag.BoolVar(&cfg.debug, "debug", false, "Enable debug output")
    flag.Parse()

    // Check if we're getting piped input
    stat, _ := os.Stdin.Stat()
    cfg.isPipe = (stat.Mode() & os.ModeCharDevice) == 0
    return cfg
}

func main() {
    cfg := parseFlags()
    if cfg.showHelp {
        fmt.Println(help.GetHelp())
        return
    }

    aiClient := client.NewOllamaClient("http://localhost:11434")

    if cfg.listModels {
        handleListModels(aiClient)
        return
    }

    userPrompt := readInitialPrompt(cfg.isPipe)

    userPrompt = strings.TrimSpace(userPrompt)
    if userPrompt == "" {
        fmt.Println(help.GetErrorNoPrompt())
        return
    }

    selectedModel := selectModel(cfg)
    if cfg.debug {
        printDebugInfo(cfg, selectedModel, userPrompt)
    }

    chatSession := aiClient.NewChatSession(selectedModel, cfg.systemPrompt, cfg.assistPrompt)
    
    if err := handleInitialPrompt(chatSession, cfg, selectedModel, userPrompt); err != nil {
        fmt.Printf("Error in chat: %v\n", err)
        os.Exit(1)
    }

    if !cfg.isPipe {
        handleInteractiveChat(chatSession, cfg, selectedModel)
    }
}

func readInitialPrompt(isPipe bool) string {
    if isPipe {
        scanner := bufio.NewScanner(os.Stdin)
        var input []string
        for scanner.Scan() {
            input = append(input, scanner.Text())
        }
        return strings.Join(input, "\n")
    }
    
    fmt.Print("Enter your prompt: ")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    return scanner.Text()
}

func handleListModels(aiClient *client.OllamaClient) {
    models, err := aiClient.ListModels()
    if err != nil {
        fmt.Printf("Error listing models: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("Available models:")
    for _, model := range models {
        fmt.Printf("- %s\n", model)
    }
}

func selectModel(cfg *config) string {
    if cfg.modelName != "" {
        return cfg.modelName
    }
    if cfg.imageFile != "" {
        return models.DefaultVisionModel
    }
    return models.DefaultTextModel
}

func printDebugInfo(cfg *config, selectedModel, userPrompt string) {
    fmt.Printf("Debug: Using model: %s\n", selectedModel)
    fmt.Printf("Debug: User prompt: %s\n", userPrompt)
    fmt.Printf("Debug: System prompt: %s\n", cfg.systemPrompt)
    if cfg.imageFile != "" {
        fmt.Printf("Debug: Image file: %s\n", cfg.imageFile)
    }
}

func handleInitialPrompt(chatSession *client.ChatSession, cfg *config, selectedModel, userPrompt string) error {
    var response string
    var err error

    if cfg.imageFile != "" {
        fmt.Println("Processing image prompt...")
        response, err = chatSession.SendWithImage(userPrompt, cfg.imageFile)
    } else {
        fmt.Println("Processing text prompt...")
        response, err = chatSession.Send(userPrompt)
    }

    if err != nil {
        return err
    }

    fmt.Printf("\n%s: %s\n", selectedModel, response)
    return nil
}

func handleInteractiveChat(chatSession *client.ChatSession, cfg *config, selectedModel string) {
    fmt.Println("\nEnter your messages (type 'exit' to quit):")
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

        if cfg.debug {
            fmt.Printf("Debug: Sending message: %s\n", input)
        }

        response, err := chatSession.Send(input)
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

