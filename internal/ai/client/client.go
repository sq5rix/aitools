// internal/ai/client/client.go
package client

type ModelInfo struct {
    Name        string `json:"name"`
    Size        int64  `json:"size"`
    ModifiedAt  string `json:"modified_at"`
    PullCount   int    `json:"pull_count,omitempty"`
    Description string `json:"description,omitempty"`
}

func (c *OllamaClient) ListModels() ([]ModelInfo, error) {
    resp, err := http.Get(c.baseURL + "/api/tags")
    if err != nil {
        return nil, fmt.Errorf("failed to list models: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to list models: %s", string(body))
    }

    var result struct {
        Models []ModelInfo `json:"models"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return result.Models, nil
}


