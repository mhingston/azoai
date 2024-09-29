package azoai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAIRequest struct {
	SystemPrompt string
	Message      string
	ApiBaseUrl   string
	APIKey       string
	APIVersion   string
	Deployment   string
	Temperature  *float64
	TopP         *float64
}

type OpenAIChatCompletionRequest struct {
	Messages    []OpenAIChatCompletionMessage `json:"messages"`
	Temperature *float64                      `json:"temperature,omitempty"`
	TopP        *float64                      `json:"top_p,omitempty"`
}

type OpenAIChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatCompletionResponse struct {
	Choices []OpenAIChatCompletionChoice `json:"choices"`
}

type OpenAIChatCompletionChoice struct {
	Message OpenAIChatCompletionMessage `json:"message"`
}

func InvokeOpenAIRequest(request OpenAIRequest) (string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"api-key":      request.APIKey,
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		request.ApiBaseUrl, request.Deployment, request.APIVersion)

	payload := OpenAIChatCompletionRequest{
		Messages: []OpenAIChatCompletionMessage{
			{Role: "system", Content: request.SystemPrompt},
			{Role: "user", Content: request.Message},
		},
		Temperature: request.Temperature,
		TopP:        request.TopP,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response status: %s", resp.Status)
	}

	var response OpenAIChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("invalid response or path does not exist")
}
