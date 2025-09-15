package actions

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/juicycleff/smartform/v1"
	"github.com/ledongthuc/pdf"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"golang.org/x/net/html"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type extractDocumentsFromEmailActionProps struct {
	MailID          string   `json:"mail_id"`
	DocumentTypes   []string `json:"document_types"`
	ExtractAll      bool     `json:"extract_all"`
	CleanupText     bool     `json:"cleanup_text"`
	IncludeMetadata bool     `json:"include_metadata"`
	MaxTextLength   int      `json:"max_text_length"`
}

type ExtractDocumentsFromEmailAction struct{}

func (a *ExtractDocumentsFromEmailAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "extract_documents_from_email",
		DisplayName:   "Extract Documents from Email",
		Description:   "Extract text content from document attachments (PDF, DOCX, TXT, HTML, XML) in Gmail emails",
		Type:          core.ActionTypeAction,
		Documentation: extractDocumentsFromEmailDocs,
		SampleOutput: map[string]any{
			"success": true,
			"mail_id": "18d4a5b6c7e8f9g0",
			"subject": "Monthly Report",
			"documents": []map[string]any{
				{
					"filename":      "report.pdf",
					"type":          "pdf",
					"size":          45678,
					"text":          "Monthly Report\nJanuary 2024\n...",
					"word_count":    250,
					"char_count":    1500,
					"attachment_id": "ANGjdJ_8x...",
					"metadata": map[string]any{
						"pages":  5,
						"format": "pdf",
					},
				},
				{
					"filename":      "data.docx",
					"type":          "docx",
					"size":          23456,
					"text":          "Data Analysis\n...",
					"word_count":    180,
					"char_count":    950,
					"attachment_id": "BGHkdK_9y...",
					"metadata": map[string]any{
						"paragraphs": 12,
						"format":     "docx",
					},
				},
			},
			"total_documents_found": 2,
			"document_types_found":  []string{"pdf", "docx"},
			"email_metadata": map[string]any{
				"from": "reports@company.com",
				"date": "2024-01-15T10:30:00Z",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ExtractDocumentsFromEmailAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("extract_documents_from_email", "Extract Documents from Email")

	form.TextField("mail_id", "Mail ID").
		Placeholder("Enter message ID").
		Required(true).
		HelpText("The message ID of the email containing document attachments")

	form.MultiSelectField("document_types", "Document Types").
		Required(false).
		AddOptions([]*smartform.Option{
			{Value: "pdf", Label: "PDF (.pdf)"},
			{Value: "docx", Label: "Word (.docx)"},
			{Value: "txt", Label: "Text (.txt)"},
			{Value: "html", Label: "HTML (.html, .htm)"},
			{Value: "xml", Label: "XML (.xml)"},
			{Value: "csv", Label: "CSV (.csv)"},
			{Value: "md", Label: "Markdown (.md)"},
			{Value: "rtf", Label: "Rich Text (.rtf)"},
		}...).
		DefaultValue([]string{"pdf", "docx", "txt"}).
		HelpText("Select which document types to extract. Leave empty to extract all supported types.")

	form.CheckboxField("extract_all", "Extract All Documents").
		DefaultValue(true).
		HelpText("Extract text from all matching document attachments")

	form.CheckboxField("cleanup_text", "Clean Up Text").
		DefaultValue(true).
		HelpText("Remove extra whitespace and clean up the extracted text")

	form.CheckboxField("include_metadata", "Include Metadata").
		DefaultValue(true).
		HelpText("Include document metadata and email details")

	form.NumberField("max_text_length", "Max Text Length per Document").
		DefaultValue(0).
		HelpText("Maximum characters to extract per document (0 for unlimited)")

	return form.Build()
}

func (a *ExtractDocumentsFromEmailAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ExtractDocumentsFromEmailAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[extractDocumentsFromEmailActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.MailID == "" {
		return nil, errors.New("mail ID is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Fetch the email with full format
	message, err := gmailService.Users.Messages.Get("me", input.MailID).
		Format("full").
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch email: %v", err)
	}

	result := map[string]any{
		"success": true,
		"mail_id": input.MailID,
	}

	// Get email subject and metadata
	headers := extractGmailHeaders(message.Payload)
	result["subject"] = headers["Subject"]

	if input.IncludeMetadata {
		result["email_metadata"] = map[string]any{
			"from":       headers["From"],
			"to":         headers["To"],
			"date":       headers["Date"],
			"message_id": headers["Message-ID"],
		}
	}

	// Determine which document types to extract
	supportedTypes := getSupportedTypes(input.DocumentTypes)

	// Find and extract document attachments
	documentResults := []map[string]any{}
	documentCount := 0
	documentTypesFound := make(map[string]bool)

	// Process attachments
	err = processMessagePartForDocuments(gmailService, message.Id, message.Payload, input, supportedTypes, &documentResults, &documentCount, &documentTypesFound)
	if err != nil {
		return nil, fmt.Errorf("error processing attachments: %v", err)
	}

	if documentCount == 0 {
		return nil, fmt.Errorf("no supported document attachments found in email")
	}

	// Convert map to slice for document types found
	typesFound := []string{}
	for docType := range documentTypesFound {
		typesFound = append(typesFound, docType)
	}

	result["documents"] = documentResults
	result["total_documents_found"] = documentCount
	result["document_types_found"] = typesFound

	return result, nil
}

func getSupportedTypes(requestedTypes []string) map[string]bool {
	supportedTypes := make(map[string]bool)

	// If no types specified, support all
	if len(requestedTypes) == 0 {
		supportedTypes["pdf"] = true
		supportedTypes["docx"] = true
		supportedTypes["txt"] = true
		supportedTypes["html"] = true
		supportedTypes["xml"] = true
		supportedTypes["csv"] = true
		supportedTypes["md"] = true
		supportedTypes["rtf"] = true
		return supportedTypes
	}

	// Otherwise, only support requested types
	for _, docType := range requestedTypes {
		supportedTypes[strings.ToLower(docType)] = true
	}
	return supportedTypes
}

func processMessagePartForDocuments(service *gmail.Service, messageID string, part *gmail.MessagePart, input *extractDocumentsFromEmailActionProps, supportedTypes map[string]bool, documentResults *[]map[string]any, documentCount *int, documentTypesFound *map[string]bool) error {
	// Check if this part is a document attachment
	if part.Filename != "" {
		docType := getDocumentType(part.Filename)

		if docType != "" && supportedTypes[docType] {
			*documentCount++
			(*documentTypesFound)[docType] = true

			if !input.ExtractAll && len(*documentResults) > 0 {
				return nil // Only extract first document if ExtractAll is false
			}

			// Get attachment data
			var attachmentData []byte
			var err error

			if part.Body != nil && part.Body.AttachmentId != "" {
				// Fetch attachment
				attachment, err := service.Users.Messages.Attachments.Get("me", messageID, part.Body.AttachmentId).Do()
				if err != nil {
					return fmt.Errorf("failed to fetch attachment %s: %v", part.Filename, err)
				}

				// Decode base64 data
				attachmentData, err = base64.URLEncoding.DecodeString(attachment.Data)
				if err != nil {
					return fmt.Errorf("failed to decode attachment %s: %v", part.Filename, err)
				}
			} else if part.Body != nil && part.Body.Data != "" {
				// Inline attachment
				attachmentData, err = base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					return fmt.Errorf("failed to decode inline attachment %s: %v", part.Filename, err)
				}
			}

			if len(attachmentData) > 0 {
				// Extract text based on document type
				var extractedText string
				var metadata map[string]interface{}

				switch docType {
				case "pdf":
					extractedText, metadata, err = extractFromPDF(attachmentData)
				case "docx":
					extractedText, metadata, err = extractFromDOCX(attachmentData)
				case "txt", "csv", "md", "rtf":
					extractedText, metadata, err = extractFromText(attachmentData, docType)
				case "html":
					extractedText, metadata, err = extractFromHTML(attachmentData)
				case "xml":
					extractedText, metadata, err = extractFromXML(attachmentData)
				default:
					err = fmt.Errorf("unsupported document type: %s", docType)
				}

				if err != nil {
					// Log error but don't stop processing other documents
					// Since we're not in a loop, we just skip adding this document to results
					fmt.Printf("Warning: failed to extract text from %s: %v\n", part.Filename, err)
				} else {
					// Only process if extraction was successful
					// Clean up text if requested
					if input.CleanupText {
						extractedText = cleanupText(extractedText)
					}

					// Apply max length if specified
					if input.MaxTextLength > 0 && len(extractedText) > input.MaxTextLength {
						extractedText = extractedText[:input.MaxTextLength] + "..."
					}

					documentResult := map[string]any{
						"filename":   part.Filename,
						"type":       docType,
						"size":       part.Body.Size,
						"text":       extractedText,
						"word_count": countWords(extractedText),
						"char_count": len(extractedText),
					}

					if part.Body.AttachmentId != "" {
						documentResult["attachment_id"] = part.Body.AttachmentId
					}

					if input.IncludeMetadata && metadata != nil {
						documentResult["metadata"] = metadata
					}

					*documentResults = append(*documentResults, documentResult)
				}
			}
		}
	}

	// Recursively process nested parts
	for _, subPart := range part.Parts {
		if err := processMessagePartForDocuments(service, messageID, subPart, input, supportedTypes, documentResults, documentCount, documentTypesFound); err != nil {
			return err
		}
	}

	return nil
}

func getDocumentType(filename string) string {
	lowerFilename := strings.ToLower(filename)

	switch {
	case strings.HasSuffix(lowerFilename, ".pdf"):
		return "pdf"
	case strings.HasSuffix(lowerFilename, ".docx"):
		return "docx"
	case strings.HasSuffix(lowerFilename, ".txt"):
		return "txt"
	case strings.HasSuffix(lowerFilename, ".html") || strings.HasSuffix(lowerFilename, ".htm"):
		return "html"
	case strings.HasSuffix(lowerFilename, ".xml"):
		return "xml"
	case strings.HasSuffix(lowerFilename, ".csv"):
		return "csv"
	case strings.HasSuffix(lowerFilename, ".md"):
		return "md"
	case strings.HasSuffix(lowerFilename, ".rtf"):
		return "rtf"
	default:
		return ""
	}
}

func cleanupText(text string) string {
	// Replace multiple consecutive whitespaces with single space
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	// Replace multiple consecutive newlines with double newlines
	re = regexp.MustCompile(`\n\s*\n\s*\n+`)
	text = re.ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

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

func cleanPDFString(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' || r == '\r' {
			if r != '\ufffd' && r != '\u0000' && r != '\u00a0' {
				if r == '\u00a0' {
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

	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}

		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				title = strings.TrimSpace(n.FirstChild.Data)
			}
		}

		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		if n.Type == html.ElementNode {
			switch n.Data {
			case "br", "p", "div", "li":
				textBuilder.WriteString("\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				textBuilder.WriteString("\n\n")
			}
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

func extractFromText(data []byte, fileType string) (string, map[string]interface{}, error) {
	text := string(data)
	lines := strings.Split(text, "\n")

	metadata := map[string]interface{}{
		"format": fileType,
		"lines":  len(lines),
	}

	if fileType == "csv" && len(lines) > 0 {
		metadata["columns"] = len(strings.Split(lines[0], ","))
	}

	return text, metadata, nil
}

func NewExtractDocumentsFromEmailAction() sdk.Action {
	return &ExtractDocumentsFromEmailAction{}
}
