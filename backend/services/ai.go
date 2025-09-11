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
	// Try AI services if available, otherwise create a simple summary
	summary, err := callHuggingFace("microsoft/DialoGPT-medium", content)
	if err != nil {
		// Fallback to simple text summarization
		return createSimpleSummary(content), nil
	}

	return strings.TrimSpace(summary), nil
}

// GenerateFlashcards creates flashcards from the given text
func GenerateFlashcards(content string) ([]FlashcardData, error) {
	// Try AI services if available, otherwise create simple flashcards
	prompt := fmt.Sprintf(`Create 5 flashcards from the following text. Each flashcard should have a clear question and a concise answer. Format the response as JSON with this structure:
[
  {"question": "Question 1", "answer": "Answer 1"},
  {"question": "Question 2", "answer": "Answer 2"}
]

Text:
%s`, content)

	response, err := callHuggingFace("microsoft/DialoGPT-medium", prompt)
	if err != nil {
		// Fallback to simple flashcards
		return createSimpleFlashcards(content), nil
	}

	// Parse JSON response
	var flashcards []FlashcardData
	if err := json.Unmarshal([]byte(response), &flashcards); err != nil {
		// If JSON parsing fails, create simple flashcards
		return createSimpleFlashcards(content), nil
	}

	return flashcards, nil
}

// GenerateQuiz creates quiz questions from the given text
func GenerateQuiz(content string) ([]QuizData, error) {
	// Try AI services if available, otherwise create simple quiz
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

	response, err := callHuggingFace("microsoft/DialoGPT-medium", prompt)
	if err != nil {
		// Fallback to simple quiz
		return createSimpleQuiz(content), nil
	}

	// Parse JSON response
	var quiz []QuizData
	if err := json.Unmarshal([]byte(response), &quiz); err != nil {
		// If JSON parsing fails, create simple quiz
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
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	
	// Add HuggingFace API key if available
	if apiKey := os.Getenv("HUGGINGFACE_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

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

// createSimpleSummary creates a basic summary when AI fails
func createSimpleSummary(content string) string {
	sentences := strings.Split(content, ".")
	if len(sentences) <= 3 {
		return content // If content is short, return as is
	}
	
	// Take first 3 sentences as summary
	summary := strings.Join(sentences[:3], ".") + "."
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
