package models

const (
    DefaultTextModel   = "llama3.1:latest"
    DefaultVisionModel = "llava:13b"
)

type GenerateRequest struct {
    Model    string   `json:"model"`
    Prompt   string   `json:"prompt"`
    System   string   `json:"system,omitempty"`
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
