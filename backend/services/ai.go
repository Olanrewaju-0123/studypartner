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
	prompt := fmt.Sprintf(`Create 6-8 comprehensive educational flashcards from the following text. Each flashcard should have a clear, specific question and a detailed, accurate answer. Make questions diverse and cover different aspects of the content. Format the response as valid JSON with this exact structure:
[
  {"question": "What is the main topic discussed?", "answer": "The main topic is..."},
  {"question": "What are the key concepts?", "answer": "The key concepts include..."},
  {"question": "How does this work?", "answer": "This works by..."},
  {"question": "What are the implications?", "answer": "The implications are..."}
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
	prompt := fmt.Sprintf(`Create 6-8 comprehensive multiple choice quiz questions from the following text. Each question should have 4 unique, plausible options with one correct answer. Make questions diverse and cover different aspects of the content. Ensure all options are different and meaningful. Format the response as valid JSON with this exact structure:
[
  {
    "question": "What is the main topic discussed?",
    "options": ["The correct answer", "A plausible but wrong answer", "Another wrong option", "A third wrong option"],
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

	// Create flashcards from meaningful sentences with diverse question types
	questionTypes := []string{
		"What is the main concept discussed in: %s?",
		"Explain the key idea: %s?",
		"What does this statement mean: %s?",
		"How does this relate to the topic: %s?",
		"What is the significance of: %s?",
		"What are the implications of: %s?",
		"What can we learn from: %s?",
		"What is the purpose of: %s?",
	}

	for i, sentence := range processedSentences {
		if i >= 8 { // Create up to 8 flashcards
			break
		}
		
		// Use different question types for variety
		questionType := questionTypes[i%len(questionTypes)]
		question := fmt.Sprintf(questionType, sentence[:min(40, len(sentence))])
		
		flashcards = append(flashcards, FlashcardData{
			Question: question,
			Answer:   sentence,
		})
	}

	// If no flashcards were created, create comprehensive ones from content
	if len(flashcards) == 0 {
		words := strings.Fields(content)
		if len(words) > 10 {
			// Create flashcards from word chunks with varied questions
			chunkSize := len(words) / 5
			if chunkSize < 5 {
				chunkSize = 5
			}
			
			chunkQuestions := []string{
				"What is the main topic in this section: %s?",
				"What concept is explained here: %s?",
				"What information is provided about: %s?",
				"What does this part discuss: %s?",
				"What is the focus of this section: %s?",
			}
			
			for i := 0; i < 5 && i*chunkSize < len(words); i++ {
				start := i * chunkSize
				end := min((i+1)*chunkSize, len(words))
				chunk := strings.Join(words[start:end], " ")
				
				questionType := chunkQuestions[i%len(chunkQuestions)]
				question := fmt.Sprintf(questionType, chunk[:min(30, len(chunk))])
				
				flashcards = append(flashcards, FlashcardData{
					Question: question,
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

	// Create quiz questions from meaningful sentences with diverse options
	questionTemplates := []struct {
		question string
		options  []string
	}{
		{
			question: "What is the main concept discussed in: %s?",
			options:  []string{"The correct concept", "A related but different concept", "An unrelated concept", "A completely opposite concept"},
		},
		{
			question: "Which statement best describes: %s?",
			options:  []string{"The accurate description", "A partially correct description", "An incorrect description", "An irrelevant description"},
		},
		{
			question: "What does this statement mean: %s?",
			options:  []string{"The intended meaning", "A different interpretation", "A misunderstanding", "An unrelated meaning"},
		},
		{
			question: "How does this relate to the topic: %s?",
			options:  []string{"Directly related", "Indirectly related", "Not related", "Opposite to the topic"},
		},
		{
			question: "What is the significance of: %s?",
			options:  []string{"High significance", "Moderate significance", "Low significance", "No significance"},
		},
		{
			question: "What are the implications of: %s?",
			options:  []string{"Positive implications", "Negative implications", "Mixed implications", "No implications"},
		},
	}

	for i, sentence := range processedSentences {
		if i >= 6 { // Create up to 6 quiz questions
			break
		}
		
		// Use different question templates for variety
		template := questionTemplates[i%len(questionTemplates)]
		question := fmt.Sprintf(template.question, sentence[:min(40, len(sentence))])
		
		// Create unique options based on the sentence content
		options := make([]string, 4)
		options[0] = sentence // Correct answer
		
		// Generate plausible but incorrect alternatives
		words := strings.Fields(sentence)
		if len(words) > 3 {
			// Create alternatives by modifying the sentence
			options[1] = strings.Join(words[:len(words)/2], " ") + " with different context"
			options[2] = "This statement is partially accurate but incomplete"
			options[3] = "This statement is not mentioned in the original text"
		} else {
			// Use template options for short sentences
			options = template.options
			options[0] = sentence // Override first option with correct answer
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
			
			chunkQuestions := []string{
				"What is the main topic in this section: %s?",
				"What concept is explained here: %s?",
				"What information is provided about: %s?",
				"What does this part discuss: %s?",
				"What is the focus of this section: %s?",
			}
			
			for i := 0; i < 5 && i*chunkSize < len(words); i++ {
				start := i * chunkSize
				end := min((i+1)*chunkSize, len(words))
				chunk := strings.Join(words[start:end], " ")
				
				questionType := chunkQuestions[i%len(chunkQuestions)]
				question := fmt.Sprintf(questionType, chunk[:min(30, len(chunk))])
				
				// Create unique options for each chunk
				options := []string{
					chunk,
					fmt.Sprintf("This section discusses %s", words[min(i*2, len(words)-1)]),
					"This information is not present in the text",
					"This section covers a completely different topic",
				}
				
				quiz = append(quiz, QuizData{
					Question: question,
					Options:  options,
					Answer:   0,
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
