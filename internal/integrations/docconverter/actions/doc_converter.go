package actions

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/juicycleff/smartform/v1"
	"github.com/ledongthuc/pdf"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"golang.org/x/net/html"
)

type FileInput struct {
	ID          string      `json:"id,omitempty"`
	Ext         string      `json:"ext"`
	FileName    string      `json:"fileName"`
	MimeType    string      `json:"mimeType"`
	Path        string      `json:"path"`
	URL         string      `json:"url,omitempty"`
	DownloadURL string      `json:"downloadUrl,omitempty"`
	Size        interface{} `json:"size"`
	SizeBytes   int64       `json:"sizeBytes,omitempty"`
	Src         string      `json:"src,omitempty"`
	UploadedAt  string      `json:"uploadedAt"`
	StorageKey  string      `json:"storageKey,omitempty"`
	IsPublic    bool        `json:"isPublic,omitempty"`
}

func (f *FileInput) GetSize() (int64, error) {
	if f.SizeBytes > 0 {
		return f.SizeBytes, nil
	}

	switch v := f.Size.(type) {
	case float64:
		return int64(v), nil
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		var size int64
		_, err := fmt.Sscanf(v, "%d", &size)
		return size, err
	default:
		return 0, fmt.Errorf("unexpected size type: %T", v)
	}
}

func (f *FileInput) GetDownloadURL() string {
	// Priority: DownloadURL > URL > Path
	if f.DownloadURL != "" {
		return f.DownloadURL
	}
	if f.URL != "" {
		return f.URL
	}
	return f.Path
}

type docConverterActionProps struct {
	InputFile       *FileInput `json:"inputFile"`
	CleanupText     bool       `json:"cleanupText"`
	MaxTextLength   int        `json:"maxTextLength"`
	ExtractMetadata bool       `json:"extractMetadata"`
}

type DocConverterAction struct{}

func (a *DocConverterAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:          "document_converter",
		DisplayName: "Document Text Extractor",
		Description: "Extracts text from uploaded documents (PDF, DOCX, HTML, XML, TXT)",
		Type:        core.ActionTypeAction,
		Icon:        "file-text",
		SampleOutput: map[string]any{
			"text": "This is the extracted text from the document...",
			"metadata": map[string]interface{}{
				"format":   "pdf",
				"pages":    5,
				"filename": "document.pdf",
			},
			"wordCount":      150,
			"characterCount": 850,
			"format":         "pdf",
		},
		Settings: core.ActionSettings{},
	}
}

func (a *DocConverterAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("document_converter", "Document Text Extractor")

	form.FileField("inputFile", "Document File").
		// Required(true).
		DefaultValue(nil).
		HelpText("Upload a document file. Supports: PDF, DOCX, HTML, XML, TXT formats.")

	form.CheckboxField("extractMetadata", "Extract Metadata").
		Required(false).
		DefaultValue(true).
		HelpText("Extract document metadata such as page count, format info, etc.")

	form.CheckboxField("cleanupText", "Clean Up Text").
		Required(false).
		DefaultValue(true).
		HelpText("Remove extra whitespace and clean up the extracted text.")

	form.NumberField("maxTextLength", "Maximum Text Length").
		Required(false).
		DefaultValue(0).
		HelpText("Maximum length of extracted text (0 for unlimited). Text will be truncated if longer.")

	return form.Build()
}

func (a *DocConverterAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *DocConverterAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[docConverterActionProps](ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	if input.InputFile == nil {
		return nil, fmt.Errorf("document file is required")
	}

	if len(input.InputFile.Src) > 100 {
		fmt.Printf("  Src preview: %s...\n", input.InputFile.Src[:100])
	}

	fileSize, err := input.InputFile.GetSize()
	if err != nil {
		fmt.Printf("WARNING - Could not parse file size: %v\n", err)
		fileSize = 0
	}

	ext := strings.ToLower(filepath.Ext(input.InputFile.FileName))

	var fileContent []byte

	if input.InputFile.Src != "" && strings.HasPrefix(input.InputFile.Src, "http") {
		fmt.Printf("DEBUG - Checking if Src contains a public document URL: %s\n", input.InputFile.Src)

		content, extractedExt, err := parsePublicDocumentURL(input.InputFile.Src)
		if err == nil && len(content) > 0 {
			fileContent = content
			if extractedExt != "" && ext == "" {
				ext = extractedExt
			}
			fmt.Printf("DEBUG - Successfully extracted from public document URL, size: %d bytes\n", len(fileContent))
		} else {
			fmt.Printf("DEBUG - Not a supported public document URL or extraction failed: %v\n", err)
		}
	}

	// Method 2: Check if Src contains base64 data
	if len(fileContent) == 0 && input.InputFile.Src != "" {
		if strings.HasPrefix(input.InputFile.Src, "data:") {
			fmt.Println("DEBUG - Extracting from data URL in Src")
			parts := strings.Split(input.InputFile.Src, ",")
			if len(parts) == 2 {
				decoded, err := base64.StdEncoding.DecodeString(parts[1])
				if err == nil {
					fileContent = decoded
					fmt.Println("DEBUG - Successfully extracted from data URL")
				} else {
					fmt.Printf("DEBUG - Failed to decode data URL: %v\n", err)
				}
			}
		} else if !strings.HasPrefix(input.InputFile.Src, "http") {
			// Try as raw base64
			fmt.Println("DEBUG - Attempting to decode Src as base64")
			decoded, err := base64.StdEncoding.DecodeString(input.InputFile.Src)
			if err == nil {
				fileContent = decoded
				fmt.Println("DEBUG - Successfully decoded Src as base64")
			}
		}
	}

	// Method 3: Download from URL (Path, URL, or DownloadURL)
	if len(fileContent) == 0 {
		downloadURL := input.InputFile.GetDownloadURL()
		if downloadURL != "" {
			// Check if it's a public document URL
			if strings.Contains(downloadURL, "docs.google.com") ||
				strings.Contains(downloadURL, "drive.google.com") ||
				strings.Contains(downloadURL, "dropbox.com") {
				fmt.Printf("DEBUG - Found public document URL: %s\n", downloadURL)

				content, extractedExt, err := parsePublicDocumentURL(downloadURL)
				if err == nil && len(content) > 0 {
					fileContent = content
					if extractedExt != "" && ext == "" {
						ext = extractedExt
					}
					fmt.Printf("DEBUG - Successfully extracted from public document URL, size: %d bytes\n", len(fileContent))
				}
			}

			// If still no content, try regular download
			if len(fileContent) == 0 {
				fmt.Printf("DEBUG - Attempting to download from: %s\n", downloadURL)

				// Since these are Wakflo URLs that require authentication,
				// we need to use the context's HTTP client if available
				fileContent, err = downloadFileWithContext(ctx, downloadURL)
				if err != nil {
					// Log the error but try to provide helpful information
					fmt.Printf("DEBUG - Download failed: %v\n", err)

					// If it's a 401 error, the file exists but needs auth
					if strings.Contains(err.Error(), "401") {
						return nil, fmt.Errorf("file requires authentication. The file appears to be stored in Wakflo but requires proper authentication to access. File ID: %s, URL: %s", input.InputFile.ID, downloadURL)
					}

					return nil, fmt.Errorf("failed to retrieve file content: %v", err)
				}
				fmt.Printf("DEBUG - Successfully downloaded file, size: %d bytes\n", len(fileContent))
			}
		} else {
			return nil, fmt.Errorf("no valid file source found. Src is empty and no download URL provided")
		}
	}

	// Validate we have content
	if len(fileContent) == 0 {
		return nil, fmt.Errorf("no file content could be retrieved")
	}

	// Extract text based on file type
	var extractedText string
	var metadata map[string]interface{}

	switch ext {
	case ".pdf":
		extractedText, metadata, err = extractFromPDF(fileContent)
	case ".docx":
		extractedText, metadata, err = extractFromDOCX(fileContent)
	case ".html", ".htm":
		extractedText, metadata, err = extractFromHTML(fileContent)
	case ".xml":
		extractedText, metadata, err = extractFromXML(fileContent)
	case ".txt":
		extractedText, metadata, err = extractFromTXT(fileContent)
	default:
		return nil, fmt.Errorf("unsupported file format: %s. Supported formats: PDF, DOCX, HTML, XML, TXT", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to extract text from %s file: %v", ext, err)
	}

	// Clean up text if requested
	if input.CleanupText {
		extractedText = cleanupText(extractedText)
	}

	// Apply text length limit if specified
	if input.MaxTextLength > 0 && len(extractedText) > input.MaxTextLength {
		extractedText = extractedText[:input.MaxTextLength] + "..."
	}

	// Prepare result
	result := map[string]interface{}{
		"text":           extractedText,
		"wordCount":      countWords(extractedText),
		"characterCount": len(extractedText),
		"format":         strings.TrimPrefix(ext, "."),
		"filename":       input.InputFile.FileName,
		"fileSize":       fileSize,
	}

	if input.ExtractMetadata && metadata != nil {
		result["metadata"] = metadata
	}

	return result, nil
}

func parsePublicDocumentURL(url string) ([]byte, string, error) {
	if strings.Contains(url, "docs.google.com") {
		return extractFromGoogleDocs(url)
	}

	if strings.Contains(url, "drive.google.com") {
		return extractFromGoogleDrive(url)
	}

	if strings.Contains(url, "dropbox.com") {
		return extractFromDropbox(url)
	}

	if strings.Contains(url, "onedrive.live.com") || strings.Contains(url, "1drv.ms") {
		return extractFromOneDrive(url)
	}

	if isDirectFileURL(url) {
		content, err := downloadFile(url)
		if err != nil {
			return nil, "", err
		}
		ext := filepath.Ext(url)
		return content, ext, nil
	}

	return nil, "", fmt.Errorf("unsupported document URL type")
}

func extractFromGoogleDocs(url string) ([]byte, string, error) {
	docID := extractGoogleDocID(url)
	if docID == "" {
		return nil, "", fmt.Errorf("invalid Google Docs URL")
	}

	exportFormats := []struct {
		format string
		ext    string
	}{
		{"pdf", ".pdf"},
		{"docx", ".docx"},
		{"txt", ".txt"},
		{"html", ".html"},
	}

	for _, format := range exportFormats {
		exportURL := fmt.Sprintf("https://docs.google.com/document/d/%s/export?format=%s", docID, format.format)

		content, err := downloadFile(exportURL)
		if err == nil && len(content) > 0 {
			return content, format.ext, nil
		}

		exportURL = fmt.Sprintf("https://docs.google.com/document/u/0/d/%s/export?format=%s", docID, format.format)
		content, err = downloadFile(exportURL)
		if err == nil && len(content) > 0 {
			return content, format.ext, nil
		}
	}

	return nil, "", fmt.Errorf("unable to export Google Doc - it may not be publicly accessible")
}

func extractFromGoogleDrive(url string) ([]byte, string, error) {
	fileID := extractGoogleDriveFileID(url)
	if fileID == "" {
		return nil, "", fmt.Errorf("invalid Google Drive URL")
	}

	downloadURL := fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)

	content, err := downloadFileWithRedirects(downloadURL)
	if err != nil {
		downloadURL = fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", fileID)
		content, err = downloadFile(downloadURL)
		if err != nil {
			return nil, "", fmt.Errorf("unable to download from Google Drive: %v", err)
		}
	}

	ext := detectFileExtension(content)

	return content, ext, nil
}

func extractFromDropbox(url string) ([]byte, string, error) {
	downloadURL := strings.Replace(url, "www.dropbox.com", "dl.dropboxusercontent.com", 1)
	downloadURL = strings.Replace(downloadURL, "?dl=0", "?raw=1", 1)

	if !strings.Contains(downloadURL, "?raw=1") && !strings.Contains(downloadURL, "?dl=1") {
		if strings.Contains(downloadURL, "?") {
			downloadURL += "&raw=1"
		} else {
			downloadURL += "?raw=1"
		}
	}

	content, err := downloadFile(downloadURL)
	if err != nil {
		return nil, "", fmt.Errorf("unable to download from Dropbox: %v", err)
	}

	ext := detectFileExtension(content)
	return content, ext, nil
}

func extractFromOneDrive(url string) ([]byte, string, error) {
	content, err := downloadFile(url)
	if err != nil {
		return nil, "", fmt.Errorf("unable to download from OneDrive - may require authentication: %v", err)
	}

	ext := detectFileExtension(content)
	return content, ext, nil
}

func extractGoogleDocID(url string) string {
	patterns := []string{
		`/document/d/([a-zA-Z0-9-_]+)`,
		`/d/([a-zA-Z0-9-_]+)`,
		`id=([a-zA-Z0-9-_]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

func extractGoogleDriveFileID(url string) string {
	patterns := []string{
		`/file/d/([a-zA-Z0-9-_]+)`,
		`/d/([a-zA-Z0-9-_]+)`,
		`id=([a-zA-Z0-9-_]+)`,
		`/open\?id=([a-zA-Z0-9-_]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// isDirectFileURL checks if URL points to a direct file
func isDirectFileURL(url string) bool {
	// Check if URL ends with common file extensions
	extensions := []string{
		".pdf", ".docx", ".doc", ".txt", ".html", ".htm",
		".xml", ".csv", ".xlsx", ".xls", ".pptx", ".ppt",
		".odt", ".rtf", ".md",
	}

	lowerURL := strings.ToLower(url)
	for _, ext := range extensions {
		if strings.HasSuffix(lowerURL, ext) {
			return true
		}
	}

	return false
}

// detectFileExtension detects file type from content
func detectFileExtension(content []byte) string {
	// Check magic bytes for common file formats
	if len(content) < 4 {
		return ".txt"
	}

	// PDF
	if bytes.HasPrefix(content, []byte("%PDF")) {
		return ".pdf"
	}

	// DOCX (ZIP-based format)
	if bytes.HasPrefix(content, []byte("PK\x03\x04")) {
		// Could be DOCX, XLSX, or PPTX - need more checking
		// For simplicity, assume DOCX
		return ".docx"
	}

	// HTML
	if bytes.Contains(content[:min(1000, len(content))], []byte("<html")) ||
		bytes.Contains(content[:min(1000, len(content))], []byte("<!DOCTYPE html")) {
		return ".html"
	}

	// XML
	if bytes.HasPrefix(content, []byte("<?xml")) {
		return ".xml"
	}

	// Plain text as fallback
	return ".txt"
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// downloadFile downloads a file from URL
func downloadFile(url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid URL: %s", url)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// downloadFileWithRedirects handles HTTP redirects properly
func downloadFileWithRedirects(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// END OF NEW FUNCTIONS

// Helper to identify which source was used
func getSourceUsed(file *FileInput, hasContent bool) string {
	if hasContent {
		if file.Src != "" && strings.HasPrefix(file.Src, "data:") {
			return "data_url_from_src"
		}
		if file.Src != "" && !strings.HasPrefix(file.Src, "http") {
			return "base64_from_src"
		}
		if file.DownloadURL != "" {
			return "download_url"
		}
		if file.URL != "" {
			return "url"
		}
		if file.Path != "" {
			return "path"
		}
	}
	return "none"
}

// Download file with context (may have auth)
func downloadFileWithContext(ctx sdkcontext.PerformContext, url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid URL: %s", url)
	}

	// Try to get HTTP client from context if available
	// Some SDKs provide authenticated clients through context
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// The Wakflo SDK might provide auth headers through context
	// Check if context has methods like GetAuthHeaders() or GetHTTPClient()
	// For now, we'll try without auth and see what happens

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// cleanPDFString removes non-printable characters and fixes encoding issues
func cleanPDFString(s string) string {
	var result strings.Builder
	for _, r := range s {
		// Only keep printable characters and common whitespace
		if unicode.IsPrint(r) || r == '\n' || r == '\t' || r == '\r' {
			// Skip replacement character and other problematic Unicode
			if r != '\ufffd' && r != '\u0000' && r != '\u00a0' {
				if r == '\u00a0' {
					// Replace non-breaking space with regular space
					result.WriteRune(' ')
				} else {
					result.WriteRune(r)
				}
			}
		}
	}
	return result.String()
}

// isPrintableChar checks if a string contains a printable character
func isPrintableChar(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := rune(s[0])
	return unicode.IsPrint(r) && r != '\ufffd'
}

// reconstructTextFromChars reconstructs text from character-by-character extraction
func reconstructTextFromChars(rawText string) string {
	var result strings.Builder
	var currentWord strings.Builder

	// Clean the raw text first
	cleanedText := cleanPDFString(rawText)
	runes := []rune(cleanedText)

	// Track last character to avoid duplicates
	var lastChar rune = 0

	for i, r := range runes {
		// Skip non-printable characters
		if !unicode.IsPrint(r) && r != '\n' && r != '\t' {
			continue
		}

		// Skip duplicate consecutive characters that are likely artifacts
		// But allow legitimate doubles like "ll", "ee", "ss" in words
		if r == lastChar && i > 0 {
			// Check if this might be a legitimate double letter
			nextIsLetter := i+1 < len(runes) && unicode.IsLetter(runes[i+1])
			prevIsLetter := i-1 >= 0 && unicode.IsLetter(runes[i-1])
			currentIsLetter := unicode.IsLetter(r)

			// Skip if it's a capital letter duplicate at word boundary
			if unicode.IsUpper(r) && (!nextIsLetter || !prevIsLetter) {
				lastChar = r
				continue
			}

			// Skip if it's not part of a word (spaces, punctuation duplicates)
			if !currentIsLetter {
				lastChar = r
				continue
			}
		}

		char := string(r)

		// Check for word boundaries
		if char == " " || char == "\n" || char == "\t" {
			if currentWord.Len() > 0 {
				result.WriteString(currentWord.String())
				currentWord.Reset()
			}
			if char == "\n" {
				result.WriteString("\n")
			} else if result.Len() > 0 && !strings.HasSuffix(result.String(), " ") && !strings.HasSuffix(result.String(), "\n") {
				result.WriteString(" ")
			}
			lastChar = r
		} else if strings.ContainsAny(char, ".,;:!?()[]{}\"'") {
			// Punctuation - add current word and punctuation
			if currentWord.Len() > 0 {
				result.WriteString(currentWord.String())
				currentWord.Reset()
			}
			result.WriteString(char)
			if strings.ContainsAny(char, ".!?") && i+1 < len(runes) {
				result.WriteString(" ")
			}
			lastChar = r
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			// Regular alphanumeric character
			currentWord.WriteString(char)
			lastChar = r
		} else if char == "$" || char == "€" || char == "£" || char == "¥" {
			// Currency symbols
			if currentWord.Len() > 0 {
				result.WriteString(currentWord.String())
				result.WriteString(" ")
				currentWord.Reset()
			}
			result.WriteString(char)
			lastChar = r
		} else if char == "&" || char == "@" || char == "#" || char == "%" || char == "/" || char == "-" {
			// Special symbols that might appear in text
			if currentWord.Len() > 0 {
				result.WriteString(currentWord.String())
				currentWord.Reset()
			}
			if result.Len() > 0 && !strings.HasSuffix(result.String(), " ") {
				result.WriteString(" ")
			}
			result.WriteString(char)
			result.WriteString(" ")
			lastChar = r
		}
	}

	// Add any remaining word
	if currentWord.Len() > 0 {
		result.WriteString(currentWord.String())
	}

	return result.String()
}

// extractFromPDF extracts text from PDF using pure Go
func extractFromPDF(data []byte) (string, map[string]interface{}, error) {
	reader, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", nil, fmt.Errorf("failed to create PDF reader: %v", err)
	}

	var allTextElements []string
	pageCount := reader.NumPage()

	for pageNum := 1; pageNum <= pageCount; pageNum++ {
		page := reader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		texts := page.Content().Text
		for _, text := range texts {
			if text.S != "" {
				// Clean the text before adding it
				cleaned := cleanPDFString(text.S)
				if cleaned != "" {
					allTextElements = append(allTextElements, cleaned)
				}
			}
		}
	}

	// Join all text elements
	rawText := strings.Join(allTextElements, "")

	// First, try to detect if we have character-by-character extraction
	// Count single characters vs multi-character strings
	singleChars := 0
	multiChars := 0
	for _, elem := range allTextElements {
		trimmed := strings.TrimSpace(elem)
		if len(trimmed) == 1 && isPrintableChar(trimmed) {
			singleChars++
		} else if len(trimmed) > 1 {
			multiChars++
		}
	}

	var extractedText string

	// If mostly single characters, we need to reconstruct words
	if singleChars > multiChars*2 {
		// Character-by-character extraction detected
		extractedText = reconstructTextFromChars(rawText)
	} else {
		// Text blocks are already properly extracted
		// Just join with spaces where appropriate
		var result strings.Builder
		for i, elem := range allTextElements {
			elem = strings.TrimSpace(elem)
			if elem == "" {
				continue
			}

			if i > 0 && result.Len() > 0 {
				lastChar := result.String()[result.Len()-1]
				// Add space if needed
				if !strings.ContainsAny(string(lastChar), " \n\t") && !strings.HasPrefix(elem, " ") {
					// Check if we need space between elements
					if !strings.ContainsAny(string(lastChar), ".,;:!?") || unicode.IsLetter(rune(elem[0])) {
						result.WriteString(" ")
					}
				}
			}
			result.WriteString(elem)
		}
		extractedText = result.String()
	}

	// Final cleanup
	extractedText = strings.TrimSpace(extractedText)
	extractedText = regexp.MustCompile(`\s+`).ReplaceAllString(extractedText, " ")
	extractedText = regexp.MustCompile(`\s+([.,;:!?])`).ReplaceAllString(extractedText, "$1")
	extractedText = regexp.MustCompile(`([.,!?])\s*([A-Z])`).ReplaceAllString(extractedText, "$1 $2")
	extractedText = regexp.MustCompile(`\n\s*\n+`).ReplaceAllString(extractedText, "\n\n")

	metadata := map[string]interface{}{
		"pages":  pageCount,
		"format": "pdf",
	}

	return extractedText, metadata, nil
}

// extractFromDOCX extracts text from DOCX files
func extractFromDOCX(data []byte) (string, map[string]interface{}, error) {
	// Open DOCX as ZIP archive
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", nil, fmt.Errorf("failed to open DOCX file: %v", err)
	}

	// Find document.xml file
	var docFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "word/document.xml" {
			docFile = file
			break
		}
	}

	if docFile == nil {
		return "", nil, fmt.Errorf("document.xml not found in DOCX file")
	}

	// Read document.xml
	rc, err := docFile.Open()
	if err != nil {
		return "", nil, fmt.Errorf("failed to open document.xml: %v", err)
	}
	defer rc.Close()

	xmlData, err := io.ReadAll(rc)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read document.xml: %v", err)
	}

	// Parse XML and extract text
	type WordDocument struct {
		XMLName xml.Name `xml:"document"`
		Body    struct {
			Paragraphs []struct {
				Runs []struct {
					Text []struct {
						Content string `xml:",chardata"`
					} `xml:"t"`
				} `xml:"r"`
			} `xml:"p"`
		} `xml:"body"`
	}

	var doc WordDocument
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return "", nil, fmt.Errorf("failed to parse DOCX XML: %v", err)
	}

	var textBuilder strings.Builder
	paragraphCount := 0

	for _, paragraph := range doc.Body.Paragraphs {
		paragraphText := ""
		for _, run := range paragraph.Runs {
			for _, text := range run.Text {
				paragraphText += text.Content
			}
		}
		if strings.TrimSpace(paragraphText) != "" {
			textBuilder.WriteString(paragraphText)
			textBuilder.WriteString("\n")
			paragraphCount++
		}
	}

	metadata := map[string]interface{}{
		"format":     "docx",
		"paragraphs": paragraphCount,
	}

	return strings.TrimSpace(textBuilder.String()), metadata, nil
}

// extractFromHTML extracts text from HTML files
func extractFromHTML(data []byte) (string, map[string]interface{}, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var textBuilder strings.Builder
	var title string

	// Extract text from HTML nodes
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}

		// Extract title
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				title = strings.TrimSpace(n.FirstChild.Data)
			}
		}

		// Skip script and style tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	metadata := map[string]interface{}{
		"format": "html",
	}
	if title != "" {
		metadata["title"] = title
	}

	return strings.TrimSpace(textBuilder.String()), metadata, nil
}

// extractFromXML extracts text content from XML files
func extractFromXML(data []byte) (string, map[string]interface{}, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var textBuilder strings.Builder
	elementCount := 0

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse XML: %v", err)
		}

		switch t := token.(type) {
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		case xml.StartElement:
			elementCount++
		}
	}

	metadata := map[string]interface{}{
		"format":   "xml",
		"elements": elementCount,
	}

	return strings.TrimSpace(textBuilder.String()), metadata, nil
}

// extractFromTXT handles plain text files
func extractFromTXT(data []byte) (string, map[string]interface{}, error) {
	text := string(data)
	lines := strings.Split(text, "\n")

	metadata := map[string]interface{}{
		"format": "txt",
		"lines":  len(lines),
	}

	return text, metadata, nil
}

// cleanupText removes excessive whitespace and normalizes the text
func cleanupText(text string) string {
	// Replace multiple consecutive whitespaces with single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Replace multiple consecutive newlines with double newlines
	re = regexp.MustCompile(`\n\s*\n\s*\n+`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

// countWords counts the number of words in the text
func countWords(text string) int {
	if text == "" {
		return 0
	}
	return len(strings.Fields(text))
}

func NewDocConverterAction() sdk.Action {
	return &DocConverterAction{}
}
