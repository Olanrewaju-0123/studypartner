package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pgvector/pgvector-go"
)

type EmbeddingRequest struct {
	Inputs string `json:"inputs"`
}

type EmbeddingResponse []float32

// GenerateEmbedding creates an embedding vector for the given text
func GenerateEmbedding(text string) (pgvector.Vector, error) {
	// Truncate text if too long (most models have limits)
	if len(text) > 512 {
		text = text[:512]
	}

	// Use HuggingFace's all-MiniLM-L6-v2 model for embeddings
	embedding, err := callHuggingFaceEmbedding(text)
	if err != nil {
		return pgvector.NewVector([]float32{}), fmt.Errorf("failed to generate embedding: %w", err)
	}

	return pgvector.NewVector(embedding), nil
}

// callHuggingFaceEmbedding makes a request to HuggingFace embedding API
func callHuggingFaceEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Inputs: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/sentence-transformers/all-MiniLM-L6-v2", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// Add your HuggingFace API key here
	// req.Header.Set("Authorization", "Bearer YOUR_HUGGINGFACE_API_KEY")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var embedding EmbeddingResponse
	if err := json.Unmarshal(body, &embedding); err != nil {
		return nil, err
	}

	return embedding, nil
}

// CalculateSimilarity calculates cosine similarity between two vectors
func CalculateSimilarity(vec1, vec2 []float32) float64 {
	if len(vec1) != len(vec2) {
		return 0
	}

	var dotProduct, norm1, norm2 float64
	for i := range vec1 {
		dotProduct += float64(vec1[i] * vec2[i])
		norm1 += float64(vec1[i] * vec1[i])
		norm2 += float64(vec2[i] * vec2[i])
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (norm1 * norm2)
}
