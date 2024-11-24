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

func (c *OllamaClient) Generate(prompt, systemPrompt, model string) (string, error) {
    requestBody := models.GenerateRequest{
        Model:    model,
        Prompt:   prompt,
        System:   systemPrompt,
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

func (c *OllamaClient) GenerateWithImage(prompt, systemPrompt, model, imagePath string) (string, error) {
    imageData, err := os.ReadFile(imagePath)
    if err != nil {
        return "", fmt.Errorf("error reading image: %v", err)
    }

    requestBody := models.GenerateRequest{
        Model:    model,
        Prompt:   prompt,
        System:   systemPrompt,
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

