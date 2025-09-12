package actions

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/youtube/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type uploadVideoActionProps struct {
	ChannelID               string `json:"channel_id"`
	VideoFile               string `json:"video_file"`
	Title                   string `json:"title"`
	Description             string `json:"description"`
	Tags                    string `json:"tags"`
	CategoryID              string `json:"category_id"`
	PrivacyStatus           string `json:"privacy_status"`
	Embeddable              bool   `json:"embeddable"`
	PublicStatsViewable     bool   `json:"public_stats_viewable"`
	MadeForKids             bool   `json:"made_for_kids"`
	SelfDeclaredMadeForKids bool   `json:"self_declared_made_for_kids"`
	License                 string `json:"license"`
	RecordingDate           string `json:"recording_date"`
	DefaultLanguage         string `json:"default_language"`
	DefaultAudioLanguage    string `json:"default_audio_language"`
	NotifySubscribers       bool   `json:"notify_subscribers"`
	AutoLevels              bool   `json:"auto_levels"`
	Stabilize               bool   `json:"stabilize"`
}

type UploadVideoAction struct{}

func (a *UploadVideoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_upload_video",
		DisplayName:   "Upload YouTube Video",
		Description:   "Upload a new video to YouTube with metadata including title, description, tags, privacy settings, and more.",
		Type:          core.ActionTypeAction,
		Documentation: uploadVideoDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"id":          "dQw4w9WgXcQ",
			"title":       "My Uploaded Video",
			"description": "This is my video description",
			"status": map[string]any{
				"uploadStatus":  "uploaded",
				"privacyStatus": "private",
				"license":       "youtube",
				"embeddable":    true,
				"madeForKids":   false,
			},
			"snippet": map[string]any{
				"publishedAt":  "2023-01-01T00:00:00Z",
				"channelId":    "UCuAXFkgsw1L7xaCfnd5JJOw",
				"channelTitle": "My Channel",
				"tags":         []string{"tag1", "tag2"},
				"categoryId":   "22",
			},
			"processingDetails": map[string]any{
				"processingStatus": "processing",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UploadVideoAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("youtube_upload_video", "Upload YouTube Video")

	// Channel selection
	shared.RegisterChannelProps(form, "channel_id", "Channel", true).
		HelpText("Select the channel where you want to upload the video")

	form.FileField("video_file", "Select Video File").
		Required(true).
		HelpText("Select the video file to upload. Supported formats: MP4, AVI, MOV, WMV, FLV, 3GPP, WebM").
		AddValidation(smartform.NewValidationBuilder().FileType(
			[]string{"mp4", "avi", "mov", "wmv", "flv", "3gpp", "webm"},
			"Please upload a supported video format",
		))

	// Basic Information Section
	form.SectionField("basicInfo", "Video Information")

	form.TextField("title", "Title").
		Placeholder("Enter video title").
		HelpText("The video's title (max 100 characters, required)").
		Required(true).
		AddValidation(smartform.NewValidationBuilder().MaxLength(100, "Title must be 100 characters or less"))

	form.TextareaField("description", "Description").
		Placeholder("Enter video description").
		HelpText("The video's description (max 5000 characters)").
		Required(false).
		AddValidation(smartform.NewValidationBuilder().MaxLength(5000, "Description must be 5000 characters or less"))

	form.TextareaField("tags", "Tags").
		Placeholder("tag1, tag2, tag3").
		HelpText("Comma-separated list of tags. Tags help viewers find your video.").
		Required(false)

	// Category selection with common YouTube categories
	form.SelectField("category_id", "Category").
		AddOption("1", "Film & Animation").
		AddOption("2", "Autos & Vehicles").
		AddOption("10", "Music").
		AddOption("15", "Pets & Animals").
		AddOption("17", "Sports").
		AddOption("19", "Travel & Events").
		AddOption("20", "Gaming").
		AddOption("22", "People & Blogs").
		AddOption("23", "Comedy").
		AddOption("24", "Entertainment").
		AddOption("25", "News & Politics").
		AddOption("26", "Howto & Style").
		AddOption("27", "Education").
		AddOption("28", "Science & Technology").
		DefaultValue("22").
		HelpText("Select a category for your video").
		Required(false)

	// Privacy and Publishing Section
	form.SectionField("privacyPublishing", "Privacy and Publishing Settings")

	form.SelectField("privacy_status", "Privacy Status").
		AddOption("private", "Private").
		AddOption("unlisted", "Unlisted").
		AddOption("public", "Public").
		DefaultValue("private").
		HelpText("Set the initial privacy status. You can change this later.").
		Required(true)

	form.CheckboxField("notify_subscribers", "Notify Subscribers").
		DefaultValue(true).
		HelpText("Send notification to channel subscribers (only applies to public videos)").
		Required(false)

	form.CheckboxField("embeddable", "Allow Embedding").
		DefaultValue(true).
		HelpText("Allow the video to be embedded on other websites").
		Required(false)

	form.CheckboxField("public_stats_viewable", "Public Stats Viewable").
		DefaultValue(true).
		HelpText("Allow the video's statistics to be publicly viewable").
		Required(false)

	// Content Settings Section
	form.SectionField("contentSettings", "Content Settings")

	form.CheckboxField("made_for_kids", "Made for Kids").
		DefaultValue(false).
		HelpText("IMPORTANT: Is this video made for kids? This affects available features and COPPA compliance.").
		Required(true)

	form.CheckboxField("self_declared_made_for_kids", "Self-Declared Made for Kids").
		DefaultValue(false).
		HelpText("Channel owner's designation of whether the video is child-directed").
		Required(false)

	form.SelectField("license", "License").
		AddOption("youtube", "Standard YouTube License").
		AddOption("creativeCommon", "Creative Commons - Attribution").
		DefaultValue("youtube").
		HelpText("Choose the license for your video").
		Required(false)

	// Advanced Settings Section
	form.SectionField("advancedSettings", "Advanced Settings")

	form.DateTimeField("recording_date", "Recording Date").
		Placeholder("2023-01-01T00:00:00Z").
		HelpText("The date and time when the video was recorded (RFC 3339 format)").
		Required(false)

	form.TextField("default_language", "Default Language").
		Placeholder("en").
		HelpText("The language of the video's default audio track (ISO 639-1 two-letter code, e.g., 'en' for English)").
		Required(false)

	form.TextField("default_audio_language", "Default Audio Language").
		Placeholder("en").
		HelpText("The language of the video's default audio track (ISO 639-1 two-letter code)").
		Required(false)

	// Processing Options Section
	form.SectionField("processingOptions", "Processing Options")

	form.CheckboxField("auto_levels", "Auto-levels").
		DefaultValue(false).
		HelpText("Automatically adjust lighting and color levels").
		Required(false)

	form.CheckboxField("stabilize", "Stabilize").
		DefaultValue(false).
		HelpText("Apply video stabilization to reduce camera shake").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *UploadVideoAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *UploadVideoAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[uploadVideoActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.VideoFile == "" {
		return nil, errors.New("video file is required")
	}

	if input.Title == "" {
		return nil, errors.New("video title is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Create the video resource
	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       input.Title,
			Description: input.Description,
			CategoryId:  input.CategoryID,
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus:           input.PrivacyStatus,
			Embeddable:              input.Embeddable,
			PublicStatsViewable:     input.PublicStatsViewable,
			MadeForKids:             input.MadeForKids,
			SelfDeclaredMadeForKids: input.SelfDeclaredMadeForKids,
			License:                 input.License,
		},
	}

	// Add tags if provided
	if input.Tags != "" {
		// Split tags by comma and trim whitespace
		tags := strings.Split(input.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		upload.Snippet.Tags = tags
	}

	// Add language settings if provided
	if input.DefaultLanguage != "" {
		upload.Snippet.DefaultLanguage = input.DefaultLanguage
	}
	if input.DefaultAudioLanguage != "" {
		upload.Snippet.DefaultAudioLanguage = input.DefaultAudioLanguage
	}

	// Add recording details if provided
	if input.RecordingDate != "" {
		upload.RecordingDetails = &youtube.VideoRecordingDetails{
			RecordingDate: input.RecordingDate,
		}
	}

	// Open the video file
	file, err := os.Open(input.VideoFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open video file: %v", err)
	}
	defer file.Close()

	// Get file info for progress tracking
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// Create the API call
	call := youtubeService.Videos.Insert([]string{"snippet", "status", "recordingDetails"}, upload)

	// Set upload parameters
	if input.NotifySubscribers {
		call = call.NotifySubscribers(input.NotifySubscribers)
	}
	if input.AutoLevels {
		call = call.AutoLevels(input.AutoLevels)
	}
	if input.Stabilize {
		call = call.Stabilize(input.Stabilize)
	}

	// Set channel ID if provided
	if input.ChannelID != "" {
		call = call.OnBehalfOfContentOwnerChannel(input.ChannelID)
	}

	// Create a progress reporter
	progressReader := &progressReader{
		Reader: file,
		Total:  fileInfo.Size(),
		ctx:    ctx,
	}

	// Upload the video
	response, err := call.Media(progressReader, googleapi.ChunkSize(1024*1024)).Do() // 1MB chunks
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %v", err)
	}

	// Build the response
	result := map[string]any{
		"id":          response.Id,
		"title":       response.Snippet.Title,
		"description": response.Snippet.Description,
		"status": map[string]any{
			"uploadStatus":            response.Status.UploadStatus,
			"privacyStatus":           response.Status.PrivacyStatus,
			"license":                 response.Status.License,
			"embeddable":              response.Status.Embeddable,
			"publicStatsViewable":     response.Status.PublicStatsViewable,
			"madeForKids":             response.Status.MadeForKids,
			"selfDeclaredMadeForKids": response.Status.SelfDeclaredMadeForKids,
		},
		"snippet": map[string]any{
			"publishedAt":  response.Snippet.PublishedAt,
			"channelId":    response.Snippet.ChannelId,
			"channelTitle": response.Snippet.ChannelTitle,
			"tags":         response.Snippet.Tags,
			"categoryId":   response.Snippet.CategoryId,
		},
	}

	// Add processing details if available
	if response.ProcessingDetails != nil {
		result["processingDetails"] = map[string]any{
			"processingStatus": response.ProcessingDetails.ProcessingStatus,
		}
		if response.ProcessingDetails.ProcessingProgress != nil {
			result["processingProgress"] = map[string]any{
				"partsTotal":     response.ProcessingDetails.ProcessingProgress.PartsTotal,
				"partsProcessed": response.ProcessingDetails.ProcessingProgress.PartsProcessed,
				"timeLeftMs":     response.ProcessingDetails.ProcessingProgress.TimeLeftMs,
			}
		}
	}

	return result, nil
}

// progressReader wraps an io.Reader to report upload progress
type progressReader struct {
	io.Reader
	Total   int64
	Current int64
	ctx     sdkcontext.PerformContext
	lastPct int
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)

	// Calculate percentage
	if pr.Total > 0 {
		pct := int(float64(pr.Current) * 100 / float64(pr.Total))
		// Only log every 10% to avoid spam
		if pct > pr.lastPct && pct%10 == 0 {
			pr.lastPct = pct
			// You could emit progress events here if the SDK supports it
			// For now, we'll just log it internally
		}
	}

	return n, err
}

// downloadFileWithContext downloads a file with context (may have auth)
func downloadFileWithContext(ctx sdkcontext.PerformContext, url string) ([]byte, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, fmt.Errorf("invalid URL: %s", url)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

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

func NewUploadVideoAction() sdk.Action {
	return &UploadVideoAction{}
}
