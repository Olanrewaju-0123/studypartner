package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	prompt := fmt.Sprintf(`Please provide a concise summary of the following text. Focus on the main points and key concepts:

%s

Summary:`, content)

	// Try Ollama first, fallback to HuggingFace
	summary, err := callOllama("llama2", prompt)
	if err != nil {
		// Fallback to HuggingFace
		summary, err = callHuggingFace("microsoft/DialoGPT-medium", prompt)
		if err != nil {
			return "", fmt.Errorf("failed to generate summary: %w", err)
		}
	}

	return strings.TrimSpace(summary), nil
}

// GenerateFlashcards creates flashcards from the given text
func GenerateFlashcards(content string) ([]FlashcardData, error) {
	prompt := fmt.Sprintf(`Create 5 flashcards from the following text. Each flashcard should have a clear question and a concise answer. Format the response as JSON with this structure:
[
  {"question": "Question 1", "answer": "Answer 1"},
  {"question": "Question 2", "answer": "Answer 2"}
]

Text:
%s`, content)

	response, err := callOllama("llama2", prompt)
	if err != nil {
		// Fallback to HuggingFace
		response, err = callHuggingFace("microsoft/DialoGPT-medium", prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to generate flashcards: %w", err)
		}
	}

	// Parse JSON response
	var flashcards []FlashcardData
	if err := json.Unmarshal([]byte(response), &flashcards); err != nil {
		// If JSON parsing fails, create simple flashcards
		flashcards = createSimpleFlashcards(content)
	}

	return flashcards, nil
}

// GenerateQuiz creates quiz questions from the given text
func GenerateQuiz(content string) ([]QuizData, error) {
	prompt := fmt.Sprintf(`Create 5 multiple choice quiz questions from the following text. Each question should have 4 options (A, B, C, D) with one correct answer. Format the response as JSON with this structure:
[
  {
    "question": "Question 1?",
    "options": ["Option A", "Option B", "Option C", "Option D"],
    "answer": 0
  }
]

Text:
%s`, content)

	response, err := callOllama("llama2", prompt)
	if err != nil {
		// Fallback to HuggingFace
		response, err = callHuggingFace("microsoft/DialoGPT-medium", prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to generate quiz: %w", err)
		}
	}

	// Parse JSON response
	var quiz []QuizData
	if err := json.Unmarshal([]byte(response), &quiz); err != nil {
		// If JSON parsing fails, create simple quiz
		quiz = createSimpleQuiz(content)
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
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	// Add your HuggingFace API key here
	// req.Header.Set("Authorization", "Bearer YOUR_HUGGINGFACE_API_KEY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var hfResp []HuggingFaceResponse
	if err := json.Unmarshal(body, &hfResp); err != nil {
		return "", err
	}

	if len(hfResp) > 0 {
		return hfResp[0].GeneratedText, nil
	}

	return "", fmt.Errorf("no response from HuggingFace")
}

// createSimpleFlashcards creates basic flashcards when AI fails
func createSimpleFlashcards(content string) []FlashcardData {
	sentences := strings.Split(content, ".")
	var flashcards []FlashcardData

	for i, sentence := range sentences {
		if i >= 5 || len(strings.TrimSpace(sentence)) < 10 {
			break
		}
		
		flashcards = append(flashcards, FlashcardData{
			Question: fmt.Sprintf("What is mentioned about: %s?", strings.TrimSpace(sentence)[:min(50, len(strings.TrimSpace(sentence)))]),
			Answer:   strings.TrimSpace(sentence),
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

		quiz = append(quiz, QuizData{
			Question: fmt.Sprintf("Which of the following is true about: %s?", strings.TrimSpace(sentence)[:min(50, len(strings.TrimSpace(sentence)))]),
			Options:  []string{strings.TrimSpace(sentence), "Option B", "Option C", "Option D"},
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
