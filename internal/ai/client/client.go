package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "io"

    "github.com/sq5rix/aitools/internal/ai/models"
)

type OllamaClient struct {
    BaseURL string
}

func NewOllamaClient(baseURL string) *OllamaClient {
    return &OllamaClient{
        BaseURL: baseURL,
    }
}

func New(baseURL string) *OllamaClient {
    return &OllamaClient{
        BaseURL: baseURL,
    }
}

func (c *OllamaClient) Generate(prompt, systemPrompt, assistPrompt, model string) (string, error) {
    requestBody := models.GenerateRequest{
        Model:    model,
        Prompt:   prompt,
        System:   systemPrompt,
        Assist:   assistPrompt,
        Stream:   false,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

    fmt.Printf("Debug: Sending request to %s/api/generate\n", c.BaseURL)
    resp, err := http.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
    }

    var response models.GenerateResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    return response.Response, nil
}

func (c *OllamaClient) GenerateWithImage(prompt, systemPrompt, assistPrompt, model, imagePath string) (string, error) {
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
        return "", fmt.Errorf("error reading image: %v", err)
    }

    requestBody := models.GenerateRequest{
        Model:    model,
        Prompt:   prompt,
        System:   systemPrompt,
        Assist:   assistPrompt,
        Stream:   false,
        Images:   []string{string(imageData)},
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

    resp, err := http.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    var response models.GenerateResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    return response.Response, nil
}

// ChatSession represents an ongoing conversation with the model
type ChatSession struct {
    client       *OllamaClient
    model        string
    systemPrompt string
    assistPrompt string
    messages     []models.ChatMessage
}

// NewChatSession creates a new chat session with the specified model and prompts
func (c *OllamaClient) NewChatSession(model, systemPrompt, assistPrompt string) *ChatSession {
    session := &ChatSession{
        client:       c,
        model:        model,
        systemPrompt: systemPrompt,
        assistPrompt: assistPrompt,
        messages:     make([]models.ChatMessage, 0),
    }

    // Add system message if provided
    if systemPrompt != "" {
        session.messages = append(session.messages, models.ChatMessage{
            Role:    "system",
            Content: systemPrompt,
        })
    }

    // Add assistant context if provided
    if assistPrompt != "" {
        session.messages = append(session.messages, models.ChatMessage{
            Role:    "assistant",
            Content: assistPrompt,
        })
    }

    return session
}

// Send sends a message to the chat session and returns the model's response
func (s *ChatSession) Send(message string) (string, error) {
    // Add user message to history
    s.messages = append(s.messages, models.ChatMessage{
        Role:    "user",
        Content: message,
    })

    requestBody := models.ChatRequest{
        Model:    s.model,
        Messages: s.messages,
        Stream:   false,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

    resp, err := http.Post(s.client.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
    }

    var response models.ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    // Add assistant's response to history
    s.messages = append(s.messages, models.ChatMessage{
        Role:    "assistant",
        Content: response.Message.Content,
    })

    return response.Message.Content, nil
}

// SendWithImage sends a message with an image to the chat session
func (s *ChatSession) SendWithImage(message, imagePath string) (string, error) {
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
        return "", fmt.Errorf("error reading image: %v", err)
    }

    // Add user message with image to history
    s.messages = append(s.messages, models.ChatMessage{
        Role:    "user",
        Content: message,
        Images:  []string{string(imageData)},
    })

    requestBody := models.ChatRequest{
        Model:    s.model,
        Messages: s.messages,
        Stream:   false,
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("error marshaling request: %v", err)
    }

    resp, err := http.Post(s.client.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
    }

    var response models.ChatResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    // Add assistant's response to history
    s.messages = append(s.messages, models.ChatMessage{
        Role:    "assistant",
        Content: response.Message.Content,
    })

    return response.Message.Content, nil
}

func (c *OllamaClient) ListModels() ([]string, error) {
    resp, err := http.Get(c.BaseURL + "/api/tags")
    if err != nil {
        return nil, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    var response models.ListModelsResponse
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("error decoding response: %v", err)
    }

    // Extract just the model names
    modelNames := make([]string, len(response.Models))
    for i, model := range response.Models {
        modelNames[i] = model.Name
    }

    return modelNames, nil
}

