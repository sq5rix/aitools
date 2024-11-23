package models

const (
    DefaultTextModel   = "llama2:3.1"
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
    Models []string `json:"models"`
}

