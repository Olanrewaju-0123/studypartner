package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type HuggingFaceRequest struct {
	Inputs string `json:"inputs"`
}

type HuggingFaceResponse struct {
	GeneratedText string `json:"generated_text"`
}

// GenerateSummary creates a summary of the given text using AI
func GenerateSummary(content string) (string, error) {
	// Validate input content
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("content cannot be empty")
	}

	// Create a proper summarization prompt
	prompt := fmt.Sprintf(`Please provide a comprehensive summary of the following text. The summary should be clear, well-structured, and capture the main points and key concepts:

%s

Summary:`, content)

	// Try AI services if available, otherwise create a simple summary
	summary, err := callHuggingFace("facebook/bart-large-cnn", prompt)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("HuggingFace API failed: %v, falling back to simple summary\n", err)
		// Fallback to simple text summarization
		return createSimpleSummary(content), nil
	}

	// Validate the AI response
	summary = strings.TrimSpace(summary)
	if summary == "" || len(summary) < 50 {
		fmt.Printf("AI returned poor summary, falling back to simple summary\n")
		return createSimpleSummary(content), nil
	}

	return summary, nil
}

// GenerateFlashcards creates flashcards from the given text
func GenerateFlashcards(content string) ([]FlashcardData, error) {
	// Try AI services if available, otherwise create simple flashcards
	prompt := fmt.Sprintf(`Create 5 educational flashcards from the following text. Each flashcard should have a clear, specific question and a concise, accurate answer. Format the response as valid JSON with this exact structure:
[
  {"question": "What is the main topic discussed?", "answer": "The main topic is..."},
  {"question": "What are the key concepts?", "answer": "The key concepts include..."}
]

Text to create flashcards from:
%s

Return only the JSON array, no additional text:`, content)

	response, err := callHuggingFace("microsoft/DialoGPT-medium", prompt)
	if err != nil {
		fmt.Printf("HuggingFace API failed for flashcards: %v, falling back to simple flashcards\n", err)
		// Fallback to simple flashcards
		return createSimpleFlashcards(content), nil
	}

	// Clean the response - remove any non-JSON text
	response = strings.TrimSpace(response)
	if strings.Contains(response, "```json") {
		start := strings.Index(response, "```json") + 7
		end := strings.Index(response[start:], "```")
		if end != -1 {
			response = response[start : start+end]
		}
	} else if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		end := strings.Index(response[start:], "```")
		if end != -1 {
			response = response[start : start+end]
		}
	}

	// Parse JSON response
	var flashcards []FlashcardData
	if err := json.Unmarshal([]byte(response), &flashcards); err != nil {
		fmt.Printf("JSON parsing failed for flashcards: %v, falling back to simple flashcards\n", err)
		// If JSON parsing fails, create simple flashcards
		return createSimpleFlashcards(content), nil
	}

	// Validate flashcards
	if len(flashcards) == 0 {
		fmt.Printf("No flashcards generated, falling back to simple flashcards\n")
		return createSimpleFlashcards(content), nil
	}

	return flashcards, nil
}

// GenerateQuiz creates quiz questions from the given text
func GenerateQuiz(content string) ([]QuizData, error) {
	// Try AI services if available, otherwise create simple quiz
	prompt := fmt.Sprintf(`Create 5 multiple choice quiz questions from the following text. Each question should have 4 options with one correct answer. Format the response as valid JSON with this exact structure:
[
  {
    "question": "What is the main topic discussed?",
    "options": ["Option A", "Option B", "Option C", "Option D"],
    "answer": 0
  }
]

Text to create quiz from:
%s

Return only the JSON array, no additional text:`, content)

	response, err := callHuggingFace("microsoft/DialoGPT-medium", prompt)
	if err != nil {
		fmt.Printf("HuggingFace API failed for quiz: %v, falling back to simple quiz\n", err)
		// Fallback to simple quiz
		return createSimpleQuiz(content), nil
	}

	// Clean the response - remove any non-JSON text
	response = strings.TrimSpace(response)
	if strings.Contains(response, "```json") {
		start := strings.Index(response, "```json") + 7
		end := strings.Index(response[start:], "```")
		if end != -1 {
			response = response[start : start+end]
		}
	} else if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		end := strings.Index(response[start:], "```")
		if end != -1 {
			response = response[start : start+end]
		}
	}

	// Parse JSON response
	var quiz []QuizData
	if err := json.Unmarshal([]byte(response), &quiz); err != nil {
		fmt.Printf("JSON parsing failed for quiz: %v, falling back to simple quiz\n", err)
		// If JSON parsing fails, create simple quiz
		return createSimpleQuiz(content), nil
	}

	// Validate quiz
	if len(quiz) == 0 {
		fmt.Printf("No quiz questions generated, falling back to simple quiz\n")
		return createSimpleQuiz(content), nil
	}

	return quiz, nil
}

// callOllama makes a request to the local Ollama API
func callOllama(model, prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}

// callHuggingFace makes a request to HuggingFace API
func callHuggingFace(model, prompt string) (string, error) {
	reqBody := HuggingFaceRequest{
		Inputs: prompt,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	// Add HuggingFace API key if available
	if apiKey := os.Getenv("HUGGINGFACE_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HuggingFace API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var hfResp []HuggingFaceResponse
	if err := json.Unmarshal(body, &hfResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(hfResp) > 0 {
		return hfResp[0].GeneratedText, nil
	}

	return "", fmt.Errorf("no response from HuggingFace")
}

// createSimpleSummary creates a basic summary when AI fails
func createSimpleSummary(content string) string {
	// Clean the content
	content = strings.TrimSpace(content)
	if content == "" {
		return "No content available for summary."
	}

	// Split into sentences
	sentences := strings.Split(content, ".")
	var cleanSentences []string
	
	// Clean and filter sentences
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 10 { // Only include meaningful sentences
			cleanSentences = append(cleanSentences, sentence)
		}
	}

	if len(cleanSentences) == 0 {
		return content // If no good sentences, return original
	}

	if len(cleanSentences) <= 3 {
		// If content is short, return as is with proper formatting
		return strings.Join(cleanSentences, ". ") + "."
	}
	
	// Take first 3-4 meaningful sentences as summary
	summarySentences := cleanSentences
	if len(cleanSentences) > 4 {
		summarySentences = cleanSentences[:4]
	}
	
	summary := strings.Join(summarySentences, ". ") + "."
	
	// Ensure summary is not too short
	if len(summary) < 50 {
		// Add more sentences if summary is too short
		if len(cleanSentences) > 4 {
			summary = strings.Join(cleanSentences[:min(6, len(cleanSentences))], ". ") + "."
		}
	}
	
	return summary
}

// createSimpleFlashcards creates basic flashcards when AI fails
func createSimpleFlashcards(content string) []FlashcardData {
	sentences := strings.Split(content, ".")
	var flashcards []FlashcardData

	for i, sentence := range sentences {
		if i >= 5 || len(strings.TrimSpace(sentence)) < 10 {
			break
		}
		
		cleanSentence := strings.TrimSpace(sentence)
		if len(cleanSentence) > 0 {
			question := "What is mentioned about: "
			if len(cleanSentence) > 50 {
				question += cleanSentence[:50] + "...?"
			} else {
				question += cleanSentence + "?"
			}
			
			flashcards = append(flashcards, FlashcardData{
				Question: question,
				Answer:   cleanSentence,
			})
		}
	}

	// If no flashcards were created, create a default one
	if len(flashcards) == 0 {
		flashcards = append(flashcards, FlashcardData{
			Question: "What is the main content of this note?",
			Answer:   content[:min(100, len(content))] + "...",
		})
	}

	return flashcards
}

// createSimpleQuiz creates basic quiz when AI fails
func createSimpleQuiz(content string) []QuizData {
	sentences := strings.Split(content, ".")
	var quiz []QuizData

	for i, sentence := range sentences {
		if i >= 5 || len(strings.TrimSpace(sentence)) < 10 {
			break
		}

		cleanSentence := strings.TrimSpace(sentence)
		if len(cleanSentence) > 0 {
			question := "Which statement is correct?"
			if len(cleanSentence) > 60 {
				question = fmt.Sprintf("Which statement is correct about: %s?", cleanSentence[:60]+"...")
			}
			
			quiz = append(quiz, QuizData{
				Question: question,
				Options:  []string{cleanSentence, "This is incorrect", "This is partially correct", "This is not mentioned"},
				Answer:   0,
			})
		}
	}

	// If no quiz questions were created, create a default one
	if len(quiz) == 0 {
		quiz = append(quiz, QuizData{
			Question: "What is the main topic of this note?",
			Options:  []string{content[:min(50, len(content))] + "...", "Topic B", "Topic C", "Topic D"},
			Answer:   0,
		})
	}

	return quiz
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
