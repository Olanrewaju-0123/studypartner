package services

// FlashcardData represents a flashcard structure
type FlashcardData struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// QuizData represents a quiz question structure
type QuizData struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"answer"` // Index of correct option
}
