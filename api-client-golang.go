package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

// Configuration structures
type Config struct {
	MaxRetries    int         `json:"max_retries"`
	SystemPrompt  string      `json:"system_prompt"`
	DynamicPrompt string      `json:"dynamic_prompt"`
	Encrypted     bool        `json:"encrypted"`
	Providers     Providers   `json:"providers"`
	PromptTemplate interface{} `json:"prompt_template"`
}

type Providers struct {
	Active   string                 `json:"active"`
	Models   map[string]ModelConfig `json:"models"`
}

type ModelConfig struct {
	Model         string `json:"model"`
	APIKey        string `json:"api_key"`
	ModelEndpoint string `json:"model_endpoint"`
}

// Response structures
type APIResponse struct {
	Choices []struct {
		Delta struct {
			Content    string     `json:"content"`
			ToolCalls  []ToolCall `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
}

type ToolCall struct {
	Function struct {
		Name      string `json:"name"`
		Arguments struct {
			Command   string `json:"command"`
			ToolReason string `json:"tool_reason"`
		} `json:"arguments"`
	} `json:"function"`
	ID string `json:"id"`
}

// Message types for saving conversation history
type MessageType string

const (
	UserMessage      MessageType = "user"
	AssistantMessage MessageType = "assistant"
	SystemMessage    MessageType = "system"
	ToolMessage      MessageType = "tool"
)

// Global variables
var (
	config        Config
	termBuffer    string
	spinnerActive bool
	historyCommands string
	cmdHistoryLength int
	captureHistoryLength int
)

// Utility functions
func debugLog(message string) {
	// Can be enabled/disabled based on environment variable
	if os.Getenv("DEBUG") == "true" {
		fmt.Fprintf(os.Stderr, "DEBUG: %s\n", message)
	}
}

func errorLog(message string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
}

func loadConfig() error {
	// TODO: Implement config loading from file
	// For now, we'll use a hardcoded config for demonstration
	config = Config{
		MaxRetries: 3,
		SystemPrompt: "You are a helpful assistant.",
		Encrypted: false,
		Providers: Providers{
			Active: "openrouter",
			Models: map[string]ModelConfig{
				"openrouter": {
					Model:         "openai/gpt-4-turbo",
					APIKey:        os.Getenv("OPENROUTER_API_KEY"),
					ModelEndpoint: "https://openrouter.ai/api/v1/chat/completions",
				},
			},
		},
	}
	return nil
}

func saveMessage(convID, messageType, content string) error {
	// TODO: Implement saving messages to permanent storage
	debugLog(fmt.Sprintf("Saving message to conversation %s: Type: %s, Content length: %d", 
		convID, messageType, len(content)))
	return nil
}

func createNewConversation(model, systemPrompt string) (string, error) {
	// TODO: Implement creating a new conversation in storage
	// For simplicity, we'll generate a random conversation ID
	convID := fmt.Sprintf("conv_%d", time.Now().Unix())
	debugLog(fmt.Sprintf("Created new conversation: %s with model %s", convID, model))
	return convID, nil
}

func displaySpinner() {
	if !spinnerActive {
		return
	}
	
	// Simple spinner implementation
	spinChars := []string{"|", "/", "-", "\\"}
	fmt.Print("\r" + spinChars[int(time.Now().UnixNano())%len(spinChars)])
}

func toggleSpinner(state bool) {
	spinnerActive = state
	if !state {
		fmt.Print("\r")
	}
}

func displayInBox(text string) {
	if text == "" {
		fmt.Println("\n----- End of Response -----")
		return
	}
	
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fmt.Println(line)
	}
}

func pushToStdin(command string) {
	// TODO: Implement pushing commands to stdin
	// This would typically involve using a separate process or pipe
	fmt.Println("Executing command:", command)
}

// HTTP handling
func handleHTTPResponse(statusCode int, body string, retryCount, maxRetries int, convID string) (int, error) {
	switch statusCode {
	case 200:
		debugLog("HTTP 200 OK")
		return 0, nil
		
	case 401, 145:
		message := fmt.Sprintf("Error (%d): Invalid API key", statusCode)
		saveMessage(convID, string(SystemMessage), message)
		return 1, errors.New(message)
		
	case 429:
		if retryCount < maxRetries-1 {
			message := "Rate limit exceeded, retrying..."
			saveMessage(convID, string(SystemMessage), message)
			errorLog(message)
			return 2, errors.New(message)
		} else {
			message := fmt.Sprintf("Rate limit exceeded after %d retries", maxRetries)
			saveMessage(convID, string(SystemMessage), "Error: Rate limit exceeded")
			return 1, errors.New(message)
		}
		
	case 500, 502, 503, 504:
		if retryCount < maxRetries-1 {
			message := fmt.Sprintf("Server error (%d), retrying...", statusCode)
			saveMessage(convID, string(SystemMessage), message)
			errorLog(message)
			return 2, errors.New(message)
		} else {
			message := fmt.Sprintf("Server error after %d retries", maxRetries)
			saveMessage(convID, string(SystemMessage), "Error: Server error after multiple retries")
			return 1, errors.New(message)
		}
		
	default:
		message := fmt.Sprintf("Unexpected response (HTTP %d): %s", statusCode, body)
		saveMessage(convID, string(SystemMessage), fmt.Sprintf("Error: Unexpected response (HTTP %d)", statusCode))
		return 1, errors.New(message)
	}
}

func sendJSONToAPI(ctx context.Context, modelEndpoint, apiKey, content, convID string) error {
	var fullResponse string
	buffer := ""
	retryCount := 0
	
	for retryCount < config.MaxRetries {
		toggleSpinner(true)
		
		req, err := http.NewRequestWithContext(ctx, "POST", modelEndpoint, bytes.NewBuffer([]byte(content)))
		if err != nil {
			errorLog(fmt.Sprintf("Error creating request: %v", err))
			return err
		}
		
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorLog(fmt.Sprintf("Error sending request: %v", err))
			time.Sleep(time.Duration(math.Pow(2, float64(retryCount))) * time.Second)
			retryCount++
			continue
		}
		
		code, err := handleHTTPResponse(resp.StatusCode, "", retryCount, config.MaxRetries, convID)
		if code == 1 {
			resp.Body.Close()
			return err
		} else if code == 2 {
			resp.Body.Close()
			retryCount++
			continue
		}
		
		// Process streaming response
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errorLog(fmt.Sprintf("Error reading response: %v", err))
				}
				break
			}
			
			line = strings.TrimSpace(line)
			if line == "" || line == ": OPENROUTER PROCESSING" {
				continue
			}
			
			if strings.HasPrefix(line, "CODE: ") {
				toggleSpinner(false)
				resp.Body.Close()
				return nil
			}
			
			if line == "[DONE]" {
				if buffer != "" {
					errorLog(fmt.Sprintf("ERROR, Buffer is %s", buffer))
					buffer = ""
				}
				continue
			}
			
			if strings.HasPrefix(line, "data: ") {
				if buffer != "" {
					errorLog(fmt.Sprintf("ERROR, Buffer is %s", buffer))
					buffer = ""
				}
				buffer += line[6:]
			}
			
			if strings.HasSuffix(line, "}") {
				var apiResp APIResponse
				err := json.Unmarshal([]byte(buffer), &apiResp)
				if err != nil {
					saveMessage(convID, string(SystemMessage), "Error: Invalid JSON response")
					errorLog(fmt.Sprintf("Invalid JSON received: %s", buffer))
					buffer = ""
					continue
				}
				
				toggleSpinner(false)
				
				if len(apiResp.Choices) > 0 {
					// Handle tool calls
					if len(apiResp.Choices[0].Delta.ToolCalls) > 0 {
						toolCallsJSON, _ := json.Marshal(apiResp.Choices[0].Delta.ToolCalls)
						saveMessage(convID, string(ToolMessage), string(toolCallsJSON))
						for _, toolCall := range apiResp.Choices[0].Delta.ToolCalls {
							debugLog(fmt.Sprintf("Tool call: %s", toolCall.Function.Name))
						}
						buffer = ""
						continue
					}
					
					// Handle content
					content := apiResp.Choices[0].Delta.Content
					if content != "" {
						if strings.Contains(content, "<cmd>") && strings.Contains(content, "</cmd>") {
							// Extract code between <cmd> tags
							start := strings.Index(content, "<cmd>") + 5
							end := strings.Index(content, "</cmd>")
							if start > 4 && end > start {
								code := content[start:end]
								rest := strings.Replace(content, content[start-5:end+6], "", 1)
								
								if rest != "" {
									displayInBox(rest)
								}
								
								if code != "" {
									pushToStdin(code)
								}
							}
						} else {
							displayInBox(content)
						}
						
						fullResponse += content
					} else {
						// End of response
						displayInBox("")
						saveMessage(convID, string(AssistantMessage), fullResponse)
					}
				}
				
				buffer = ""
			}
		}
		
		resp.Body.Close()
		break
	}
	
	return nil
}

func sendMessage(ctx context.Context, prompt, pipedInput, convID string) (string, error) {
	messageToSend := prompt
	
	if pipedInput != "" {
		messageToSend += "\nReference: " + pipedInput
	}
	
	// Create new conversation if no ID provided
	if convID == "" {
		activeProvider := config.Providers.Active
		modelConfig := config.Providers.Models[activeProvider]
		
		systemPrompt := config.SystemPrompt
		if config.Encrypted {
			// TODO: Implement decryption
			systemPrompt = "Decrypted: " + systemPrompt
		}
		
		// Apply dynamic prompt if configured
		dynamicPrompt := os.ExpandEnv(config.DynamicPrompt)
		systemPrompt += " " + dynamicPrompt
		
		// Add terminal buffer if available
		if termBuffer != "" {
			systemPrompt += fmt.Sprintf("\nTerminal Buffer (last %d lines): %s", 
				captureHistoryLength, termBuffer)
		}
		
		var err error
		convID, err = createNewConversation(modelConfig.Model, systemPrompt)
		if err != nil {
			return "", fmt.Errorf("failed to create conversation: %w", err)
		}
		
		// TODO: Apply proper template formatting
		// This is a simplified version for demonstration
		content := fmt.Sprintf(`{
			"model": "%s",
			"messages": [
				{"role": "system", "content": "%s"},
				{"role": "user", "content": "%s"}
			],
			"stream": true
		}`, modelConfig.Model, systemPrompt, messageToSend)
		
		// Save user message and send request
		saveMessage(convID, string(UserMessage), messageToSend)
		err = sendJSONToAPI(ctx, modelConfig.ModelEndpoint, modelConfig.APIKey, content, convID)
		if err != nil {
			return "", fmt.Errorf("API request failed: %w", err)
		}
	}
	
	return convID, nil
}

func main() {
	if err := loadConfig(); err != nil {
		errorLog(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}
	
	// Example usage
	ctx := context.Background()
	_, err := sendMessage(ctx, "Hello, how can I help you?", "", "")
	if err != nil {
		errorLog(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}
