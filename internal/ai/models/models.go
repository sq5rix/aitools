package models

const (
    DefaultTextModel   = "llama3.1:latest"
    DefaultVisionModel = "llama3.2:latest"
)

type GenerateRequest struct {
    Model    string   `json:"model"`
    Prompt   string   `json:"prompt"`
    System   string   `json:"system,omitempty"`
    Assist   string   `json:"assistant,omitempty"`
    Stream   bool     `json:"stream"`
    Images   []string `json:"images,omitempty"`
}

type GenerateResponse struct {
    Model    string `json:"model"`
    Response string `json:"response"`
    Done     bool   `json:"done"`
}

type ListModelsResponse struct {
    Models []ModelInfo `json:"models"`
}

type ModelInfo struct {
    Name    string `json:"name"`
    ModifiedAt string `json:"modified_at,omitempty"`
    Size    uint64 `json:"size,omitempty"`
}
