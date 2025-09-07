package services

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractPDFText extracts text content from PDF bytes
func ExtractPDFText(pdfBytes []byte) (string, error) {
	reader := bytes.NewReader(pdfBytes)
	pdfReader, err := pdf.NewReader(reader, int64(len(pdfBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var text strings.Builder
	for i := 1; i <= pdfReader.NumPage(); i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}

		content, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip pages that can't be read
		}
		text.WriteString(content)
		text.WriteString("\n")
	}

	return strings.TrimSpace(text.String()), nil
}

// ExtractDOCXText extracts text content from DOCX bytes
func ExtractDOCXText(docxBytes []byte) (string, error) {
	reader := bytes.NewReader(docxBytes)
	zipReader, err := zip.NewReader(reader, int64(len(docxBytes)))
	if err != nil {
		return "", fmt.Errorf("failed to create ZIP reader: %w", err)
	}

	var text strings.Builder
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			rc, err := file.Open()
			if err != nil {
				continue
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				continue
			}

			// Parse XML to extract text
			docText, err := extractTextFromXML(content)
			if err != nil {
				continue
			}
			text.WriteString(docText)
		}
	}

	return strings.TrimSpace(text.String()), nil
}

// extractTextFromXML extracts text content from DOCX XML
func extractTextFromXML(xmlData []byte) (string, error) {
	type Document struct {
		XMLName xml.Name `xml:"document"`
		Body    struct {
			Paragraphs []struct {
				Runs []struct {
					Text string `xml:"t"`
				} `xml:"r"`
			} `xml:"p"`
		} `xml:"body"`
	}

	var doc Document
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return "", err
	}

	var text strings.Builder
	for _, paragraph := range doc.Body.Paragraphs {
		for _, run := range paragraph.Runs {
			text.WriteString(run.Text)
		}
		text.WriteString("\n")
	}

	return text.String(), nil
}
