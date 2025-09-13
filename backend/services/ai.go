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

	// Try multiple AI models for better results
	models := []string{
		"facebook/bart-large-cnn",
		"google/pegasus-xsum",
		"microsoft/DialoGPT-medium",
	}

	for _, model := range models {
		summary, err := callHuggingFace(model, prompt)
		if err != nil {
			fmt.Printf("HuggingFace API failed for model %s: %v\n", model, err)
			continue
		}

		// Validate the AI response
		summary = strings.TrimSpace(summary)
		if summary != "" && len(summary) >= 50 {
			fmt.Printf("Successfully generated summary using model %s\n", model)
			return summary, nil
		}
	}

	// If all AI models fail, use enhanced fallback
	fmt.Printf("All AI models failed, using enhanced fallback summary\n")
	return createSimpleSummary(content), nil
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

	// Try multiple AI models for better results
	models := []string{
		"microsoft/DialoGPT-medium",
		"facebook/bart-large-cnn",
		"google/pegasus-xsum",
	}

	for _, model := range models {
		response, err := callHuggingFace(model, prompt)
		if err != nil {
			fmt.Printf("HuggingFace API failed for flashcards with model %s: %v\n", model, err)
			continue
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
			fmt.Printf("JSON parsing failed for flashcards with model %s: %v\n", model, err)
			continue
		}

		// Validate flashcards
		if len(flashcards) > 0 {
			fmt.Printf("Successfully generated flashcards using model %s\n", model)
			return flashcards, nil
		}
	}

	// If all AI models fail, use enhanced fallback
	fmt.Printf("All AI models failed for flashcards, using enhanced fallback\n")
	return createSimpleFlashcards(content), nil
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

	// Try multiple AI models for better results
	models := []string{
		"microsoft/DialoGPT-medium",
		"facebook/bart-large-cnn",
		"google/pegasus-xsum",
	}

	for _, model := range models {
		response, err := callHuggingFace(model, prompt)
		if err != nil {
			fmt.Printf("HuggingFace API failed for quiz with model %s: %v\n", model, err)
			continue
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
			fmt.Printf("JSON parsing failed for quiz with model %s: %v\n", model, err)
			continue
		}

		// Validate quiz
		if len(quiz) > 0 {
			fmt.Printf("Successfully generated quiz using model %s\n", model)
			return quiz, nil
		}
	}

	// If all AI models fail, use enhanced fallback
	fmt.Printf("All AI models failed for quiz, using enhanced fallback\n")
	return createSimpleQuiz(content), nil
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

// createSimpleSummary creates a comprehensive summary when AI fails
func createSimpleSummary(content string) string {
	// Clean the content
	content = strings.TrimSpace(content)
	if content == "" {
		return "No content available for summary."
	}

	// Split into sentences using multiple delimiters
	sentences := strings.FieldsFunc(content, func(c rune) bool {
		return c == '.' || c == '!' || c == '?'
	})
	
	var cleanSentences []string
	
	// Clean and filter sentences
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 15 { // Only include meaningful sentences
			cleanSentences = append(cleanSentences, sentence)
		}
	}

	if len(cleanSentences) == 0 {
		// If no good sentences, create a basic summary from the content
		words := strings.Fields(content)
		if len(words) > 20 {
			// Take first 20 words as summary
			summary := strings.Join(words[:20], " ")
			return summary + "..."
		}
		return content
	}

	// Create a comprehensive summary
	var summarySentences []string
	
	// Take first 5-7 meaningful sentences for a good summary
	maxSentences := min(7, len(cleanSentences))
	for i := 0; i < maxSentences; i++ {
		summarySentences = append(summarySentences, cleanSentences[i])
	}
	
	summary := strings.Join(summarySentences, ". ") + "."
	
	// Ensure summary is substantial (at least 100 characters)
	if len(summary) < 100 && len(cleanSentences) > maxSentences {
		// Add more sentences if summary is too short
		additionalSentences := min(3, len(cleanSentences)-maxSentences)
		for i := maxSentences; i < maxSentences+additionalSentences; i++ {
			summarySentences = append(summarySentences, cleanSentences[i])
		}
		summary = strings.Join(summarySentences, ". ") + "."
	}
	
	return summary
}

// createSimpleFlashcards creates comprehensive flashcards when AI fails
func createSimpleFlashcards(content string) []FlashcardData {
	// Split into sentences using multiple delimiters
	sentences := strings.FieldsFunc(content, func(c rune) bool {
		return c == '.' || c == '!' || c == '?'
	})
	
	var flashcards []FlashcardData
	var processedSentences []string

	// Clean and filter sentences
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 20 { // Only include meaningful sentences
			processedSentences = append(processedSentences, sentence)
		}
	}

	// Create flashcards from meaningful sentences
	for i, sentence := range processedSentences {
		if i >= 8 { // Create up to 8 flashcards
			break
		}
		
		// Create different types of questions
		var question string
		if i%3 == 0 {
			question = "What is the main point about: " + sentence[:min(40, len(sentence))] + "?"
		} else if i%3 == 1 {
			question = "Explain the concept: " + sentence[:min(35, len(sentence))] + "?"
		} else {
			question = "What does this statement mean: " + sentence[:min(45, len(sentence))] + "?"
		}
		
		flashcards = append(flashcards, FlashcardData{
			Question: question,
			Answer:   sentence,
		})
	}

	// If no flashcards were created, create comprehensive ones from content
	if len(flashcards) == 0 {
		words := strings.Fields(content)
		if len(words) > 10 {
			// Create flashcards from word chunks
			chunkSize := len(words) / 5
			if chunkSize < 5 {
				chunkSize = 5
			}
			
			for i := 0; i < 5 && i*chunkSize < len(words); i++ {
				start := i * chunkSize
				end := min((i+1)*chunkSize, len(words))
				chunk := strings.Join(words[start:end], " ")
				
				flashcards = append(flashcards, FlashcardData{
					Question: fmt.Sprintf("What is discussed in this section: %s?", chunk[:min(30, len(chunk))]),
					Answer:   chunk,
				})
			}
		} else {
			// Fallback for very short content
			flashcards = append(flashcards, FlashcardData{
				Question: "What is the main content of this note?",
				Answer:   content,
			})
		}
	}

	return flashcards
}

// createSimpleQuiz creates comprehensive quiz when AI fails
func createSimpleQuiz(content string) []QuizData {
	// Split into sentences using multiple delimiters
	sentences := strings.FieldsFunc(content, func(c rune) bool {
		return c == '.' || c == '!' || c == '?'
	})
	
	var quiz []QuizData
	var processedSentences []string

	// Clean and filter sentences
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 25 { // Only include meaningful sentences
			processedSentences = append(processedSentences, sentence)
		}
	}

	// Create quiz questions from meaningful sentences
	for i, sentence := range processedSentences {
		if i >= 6 { // Create up to 6 quiz questions
			break
		}
		
		// Create different types of questions
		var question string
		var options []string
		
		if i%3 == 0 {
			question = fmt.Sprintf("What is the main point about: %s?", sentence[:min(40, len(sentence))])
			options = []string{
				sentence,
				"This statement is incorrect",
				"This is only partially true",
				"This is not mentioned in the text",
			}
		} else if i%3 == 1 {
			question = fmt.Sprintf("Which statement best describes: %s?", sentence[:min(35, len(sentence))])
			options = []string{
				sentence,
				"The opposite of this statement",
				"This is a false statement",
				"This is not relevant to the topic",
			}
		} else {
			question = fmt.Sprintf("What does this mean: %s?", sentence[:min(45, len(sentence))])
			options = []string{
				sentence,
				"This is a misunderstanding",
				"This is partially correct",
				"This is completely wrong",
			}
		}
		
		quiz = append(quiz, QuizData{
			Question: question,
			Options:  options,
			Answer:   0, // First option is always correct
		})
	}

	// If no quiz questions were created, create comprehensive ones from content
	if len(quiz) == 0 {
		words := strings.Fields(content)
		if len(words) > 15 {
			// Create quiz questions from word chunks
			chunkSize := len(words) / 5
			if chunkSize < 8 {
				chunkSize = 8
			}
			
			for i := 0; i < 5 && i*chunkSize < len(words); i++ {
				start := i * chunkSize
				end := min((i+1)*chunkSize, len(words))
				chunk := strings.Join(words[start:end], " ")
				
				quiz = append(quiz, QuizData{
					Question: fmt.Sprintf("What is discussed in this section: %s?", chunk[:min(30, len(chunk))]),
					Options: []string{
						chunk,
						"This section discusses something else",
						"This is not mentioned in the text",
						"This is a different topic entirely",
					},
					Answer: 0,
				})
			}
		} else {
			// Fallback for very short content
			quiz = append(quiz, QuizData{
				Question: "What is the main topic of this note?",
				Options: []string{
					content,
					"This is about a different topic",
					"This is not relevant",
					"This is incorrect information",
				},
				Answer: 0,
			})
		}
	}

	return quiz
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
