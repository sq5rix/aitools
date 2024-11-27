package models

const (
    DefaultTextModel   = "llama3.1:latest"
    DefaultVisionModel = "llama3.2:latest"
)

type ModelInfo struct {
    Name        string `json:"name"`
    ModifiedAt  string `json:"modified_at"`
    Size        int64  `json:"size"`
    Digest      string `json:"digest"`
    Details     struct {
        Format      string   `json:"format"`
        Family      string   `json:"family"`
        Families    []string `json:"families"`
        ParameterSize string `json:"parameter_size"`
        QuantizationLevel string `json:"quantization_level"`
    } `json:"details"`
}

type ListModelsResponse struct {
    Models []ModelInfo `json:"models"`
}

type GenerateRequest struct {
    Model    string   `json:"model"`
    Prompt   string   `json:"prompt"`
    System   string   `json:"system,omitempty"`
    Assist   string   `json:"assist,omitempty"`
    Stream   bool     `json:"stream"`
    Images   []string `json:"images,omitempty"`
}

type GenerateResponse struct {
    Model    string `json:"model"`
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
    Role    string   `json:"role"`              // "system", "user", or "assistant"
    Content string   `json:"content"`           // The message content
    Images  []string `json:"images,omitempty"`  // Optional base64-encoded images
}

// ChatRequest represents a request to the chat API
type ChatRequest struct {
    Model    string        `json:"model"`
    Messages []ChatMessage `json:"messages"`
    Stream   bool         `json:"stream"`
    Format   string       `json:"format,omitempty"`
}

// ChatResponse represents a response from the chat API
type ChatResponse struct {
    Model          string      `json:"model"`
    Created        string      `json:"created_at"`    // Changed from int64 to string
    Message        ChatMessage `json:"message"`
    Done           bool        `json:"done"`
    TotalDuration  int64       `json:"total_duration"`
}
