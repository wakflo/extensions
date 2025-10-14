package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/youtube/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type downloadCaptionActionProps struct {
	ChannelID  string `json:"channel_id"`
	VideoID    string `json:"video_id"`
	CaptionID  string `json:"caption_id"`
	Format     string `json:"format"`
	Language   string `json:"language"`
	OutputType string `json:"output_type"`
}

type DownloadCaptionAction struct{}

func (a *DownloadCaptionAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_download_caption",
		DisplayName:   "Download YouTube Caption",
		Description:   "Download caption/subtitle tracks from your YouTube videos. You can optionally specify the format and translate to a different language.",
		Type:          core.ActionTypeAction,
		Documentation: downloadCaptionDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"videoId":          "dQw4w9WgXcQ",
			"captionId":        "caption_track_id",
			"format":           "srt",
			"language":         "en",
			"captionContent":   "1\n00:00:00,000 --> 00:00:05,000\nHello world, this is a sample caption\n\n2\n00:00:05,000 --> 00:00:10,000\nThis is the second line of the caption",
			"originalLanguage": "en",
			"isTranslated":     false,
			"filename":         "video_title_en_20240112_143022.srt",
			"mimeType":         "application/x-subrip",
			"fileSize":         1024,
			"downloadable":     true,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *DownloadCaptionAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("youtube_download_caption", "Download YouTube Caption")

	// Channel selector - always required since you can only download from your own videos
	shared.RegisterChannelProps(form, "channel_id", "Your Channel", true)

	// Dynamic video selector - filtered by channel
	shared.RegisterVideoProps(form, "video_id", "Select Video", true)

	// Caption track selector (would need to be populated dynamically based on video)
	form.TextField("caption_id", "Caption Track ID").
		Placeholder("BVieFbbEG1qqmMjyYrK1-joaYuT-4MxUSoNkxv8IGTst13xOHkN").
		HelpText("The caption track ID to download. Leave empty to automatically select the first available track.").
		Required(false)

	// Format selection
	form.SelectField("format", "Caption Format").
		AddOption("original", "Original Format").
		AddOption("srt", "SRT - SubRip Subtitle").
		AddOption("vtt", "VTT - Web Video Text Tracks").
		AddOption("sbv", "SBV - SubViewer Subtitle").
		AddOption("ttml", "TTML - Timed Text Markup Language").
		AddOption("scc", "SCC - Scenarist Closed Caption").
		DefaultValue("srt").
		HelpText("Select the format to download the captions in").
		Required(false)

	shared.RegisterLanguageProps(form, "language", "Translate to Language", false)

	// Add output type selection
	form.SelectField("output_type", "Output Type").
		AddOption("text", "Caption Text Only").
		AddOption("file", "Downloadable File").
		AddOption("both", "Both Text and File").
		DefaultValue("both").
		HelpText("Choose how you want to receive the caption data").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *DownloadCaptionAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *DownloadCaptionAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[downloadCaptionActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if input.ChannelID == "" {
		return nil, errors.New("channel selection is required")
	}
	if input.VideoID == "" {
		return nil, errors.New("video selection is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// First, verify the video belongs to the selected channel
	videoCall := youtubeService.Videos.List([]string{"snippet"}).Id(input.VideoID)
	videoResp, err := videoCall.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to verify video ownership: %w", err)
	}

	if len(videoResp.Items) == 0 {
		return nil, errors.New("video not found")
	}

	video := videoResp.Items[0]
	if video.Snippet.ChannelId != input.ChannelID {
		return nil, errors.New("you can only download captions from videos in your own channel")
	}

	// Store video title for filename generation
	videoTitle := sanitizeFilename(video.Snippet.Title)

	// If no caption ID is provided, list available captions and use the first one
	captionID := input.CaptionID
	var captionLanguage string

	if captionID == "" {
		// List captions for the video
		captionListCall := youtubeService.Captions.List([]string{"snippet"}, input.VideoID)
		captionListResp, err := captionListCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to list captions: %w", err)
		}

		if len(captionListResp.Items) == 0 {
			return nil, errors.New("no captions available for this video")
		}

		// Use the first available caption track
		captionID = captionListResp.Items[0].Id
		captionLanguage = captionListResp.Items[0].Snippet.Language

		// Try to find a caption in the requested language if translation is not needed
		if input.Language != "" {
			for _, caption := range captionListResp.Items {
				if caption.Snippet.Language == input.Language {
					captionID = caption.Id
					captionLanguage = caption.Snippet.Language
					break
				}
			}
		}
	}

	// Build the download call
	downloadCall := youtubeService.Captions.Download(captionID)

	// Set format if specified
	actualFormat := input.Format
	if actualFormat == "" {
		actualFormat = "srt"
	}

	if actualFormat != "original" {
		downloadCall = downloadCall.Tfmt(actualFormat)
	}

	// Set translation language if specified
	if input.Language != "" {
		downloadCall = downloadCall.Tlang(input.Language)
		captionLanguage = input.Language
	}

	// Execute the download
	resp, err := downloadCall.Download()
	if err != nil {
		return nil, fmt.Errorf("failed to download caption: %w", err)
	}
	defer resp.Body.Close()

	// Read the caption content
	captionContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read caption content: %w", err)
	}

	// Get caption metadata for additional info
	captionListCall := youtubeService.Captions.List([]string{"snippet"}, input.VideoID).Id(captionID)
	captionInfo, err := captionListCall.Do()
	if err != nil {
		// Don't fail if we can't get metadata, just log it
		ctx.Logger().Warn("Failed to get caption metadata", "error", err)
	}

	// Build response
	response := map[string]any{
		"videoId":        input.VideoID,
		"channelId":      input.ChannelID,
		"captionId":      captionID,
		"format":         determineFormat(actualFormat),
		"captionContent": string(captionContent),
		"contentLength":  len(captionContent),
	}

	// Add language info
	if input.Language != "" {
		response["language"] = input.Language
		response["isTranslated"] = true
	} else {
		response["isTranslated"] = false
		if captionLanguage != "" {
			response["language"] = captionLanguage
		}
	}

	// Add caption metadata if available
	if captionInfo != nil && len(captionInfo.Items) > 0 {
		caption := captionInfo.Items[0]
		response["originalLanguage"] = caption.Snippet.Language
		response["trackName"] = caption.Snippet.Name
		response["trackKind"] = caption.Snippet.TrackKind

		// Add status info if available
		if caption.Snippet.Status != "" {
			response["status"] = caption.Snippet.Status
		}
		if caption.Snippet.FailureReason != "" {
			response["failureReason"] = caption.Snippet.FailureReason
		}

		// Update language if not already set
		if captionLanguage == "" {
			captionLanguage = caption.Snippet.Language
		}
	}

	// Add content preview (first 500 characters)
	if len(captionContent) > 0 {
		preview := string(captionContent)
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		response["contentPreview"] = preview
	}

	// Handle output type
	outputType := input.OutputType
	if outputType == "" {
		outputType = "both"
	}

	// Include caption content if requested
	if outputType == "text" || outputType == "both" {
		response["captionContent"] = string(captionContent)
	}

	// Add file information for downloadable output
	if outputType == "file" || outputType == "both" {
		// Generate filename
		ext := getFileExtension(actualFormat)
		filename := generateFilename(videoTitle, captionLanguage, ext)
		mimeType := getMimeType(actualFormat)

		// Add file metadata to response
		response["filename"] = filename
		response["mimeType"] = mimeType
		response["fileSize"] = len(captionContent)

		// For Wakflo, the best approach is to create a temporary file
		// that can be downloaded. Check if the context provides file creation
		if fileCreator, ok := ctx.(interface {
			CreateDownloadableFile(filename string, content []byte, mimeType string) (string, error)
		}); ok {
			downloadURL, err := fileCreator.CreateDownloadableFile(filename, captionContent, mimeType)
			if err == nil {
				response["downloadURL"] = downloadURL
				response["downloadable"] = true
			} else {
				ctx.Logger().Warn("Failed to create downloadable file", "error", err)
			}
		}

		// Fallback: If no download URL service is available,
		// return base64 content for client-side handling
		if _, hasURL := response["downloadURL"]; !hasURL {
			response["fileContent"] = base64.StdEncoding.EncodeToString(captionContent)
			response["isBase64"] = true
			response["downloadable"] = true
		}
	}

	return response, nil
}

func determineFormat(format string) string {
	if format == "" || format == "original" {
		return "original"
	}
	return strings.ToUpper(format)
}

func getFileExtension(format string) string {
	switch strings.ToLower(format) {
	case "srt":
		return "srt"
	case "vtt":
		return "vtt"
	case "sbv":
		return "sbv"
	case "ttml":
		return "ttml"
	case "scc":
		return "scc"
	default:
		return "txt"
	}
}

func getMimeType(format string) string {
	switch strings.ToLower(format) {
	case "srt":
		return "application/x-subrip"
	case "vtt":
		return "text/vtt"
	case "sbv":
		return "text/plain"
	case "ttml":
		return "application/ttml+xml"
	case "scc":
		return "text/plain"
	default:
		return "text/plain"
	}
}

func generateFilename(videoTitle, language, extension string) string {
	// Create a timestamp
	timestamp := time.Now().Format("20060102_150405")

	// Build filename components
	parts := []string{videoTitle}

	if language != "" {
		parts = append(parts, language)
	}

	parts = append(parts, timestamp)

	// Join with underscores and add extension
	basename := strings.Join(parts, "_")
	return fmt.Sprintf("%s.%s", basename, extension)
}

func sanitizeFilename(filename string) string {
	// Remove or replace characters that are invalid in filenames
	replacements := map[string]string{
		"/":  "-",
		"\\": "-",
		":":  "-",
		"*":  "",
		"?":  "",
		"\"": "",
		"<":  "",
		">":  "",
		"|":  "",
		"\n": "_",
		"\r": "",
		"\t": "_",
	}

	result := filename
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// Trim spaces and limit length
	result = strings.TrimSpace(result)
	if len(result) > 100 {
		result = result[:100]
	}

	// Replace multiple spaces/underscores with single underscore
	result = strings.ReplaceAll(result, "  ", " ")
	result = strings.ReplaceAll(result, " ", "_")
	result = strings.ReplaceAll(result, "__", "_")

	return result
}

func NewDownloadCaptionAction() sdk.Action {
	return &DownloadCaptionAction{}
}
