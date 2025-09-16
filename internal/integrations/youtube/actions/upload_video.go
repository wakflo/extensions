package actions

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	ChannelID               string      `json:"channel_id"`
	VideoFile               interface{} `json:"video_file"` // Can be string (URL) or map
	Title                   string      `json:"title"`
	Description             string      `json:"description"`
	Tags                    string      `json:"tags"`
	CategoryID              string      `json:"category_id"`
	PrivacyStatus           string      `json:"privacy_status"`
	Embeddable              bool        `json:"embeddable"`
	PublicStatsViewable     bool        `json:"public_stats_viewable"`
	MadeForKids             bool        `json:"made_for_kids"`
	SelfDeclaredMadeForKids bool        `json:"self_declared_made_for_kids"`
	License                 string      `json:"license"`
	RecordingDate           string      `json:"recording_date"`
	DefaultLanguage         string      `json:"default_language"`
	DefaultAudioLanguage    string      `json:"default_audio_language"`
	NotifySubscribers       bool        `json:"notify_subscribers"`
	AutoLevels              bool        `json:"auto_levels"`
	Stabilize               bool        `json:"stabilize"`
	PlaylistID              string      `json:"playlist_id"` // New field for playlist
}

// FileInput structure from Wakflo
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

type UploadVideoAction struct{}

func (a *UploadVideoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_upload_video",
		DisplayName:   "Upload YouTube Video",
		Description:   "Upload a new video to YouTube with metadata including title, description, tags, privacy settings, and optionally add it to a playlist.",
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
			"playlistItem": map[string]any{
				"id":         "UExMOGg5Ym5heGdVVjNzTGFjVUF3aVBxLkNKcGJOdDl6cWJpTQ",
				"playlistId": "PLMOHg5YnaxgUV3sLacUAwiPq",
				"position":   0,
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

	// Add playlist selection using the shared helper
	shared.RegisterPlaylistProps(form, "playlist_id", "Add to Playlist", false).
		HelpText("Select a playlist to add this video to after upload (optional)")

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

	if input.VideoFile == nil {
		return nil, errors.New("video file is required")
	}

	if input.Title == "" {
		return nil, errors.New("video title is required")
	}

	// Convert map to FileInput or handle string URL
	fileInput := &FileInput{}

	switch v := input.VideoFile.(type) {
	case string:
		// Direct URL passed
		fileInput.URL = v
		fileInput.DownloadURL = v
	case map[string]interface{}:
		// Map structure passed
		if id, ok := v["id"].(string); ok {
			fileInput.ID = id
		}
		if ext, ok := v["ext"].(string); ok {
			fileInput.Ext = ext
		}
		if fileName, ok := v["fileName"].(string); ok {
			fileInput.FileName = fileName
		}
		if mimeType, ok := v["mimeType"].(string); ok {
			fileInput.MimeType = mimeType
		}
		if path, ok := v["path"].(string); ok {
			fileInput.Path = path
		}
		if url, ok := v["url"].(string); ok {
			fileInput.URL = url
		}
		if downloadURL, ok := v["downloadUrl"].(string); ok {
			fileInput.DownloadURL = downloadURL
		}
		if src, ok := v["src"].(string); ok {
			fileInput.Src = src
		}
		if size, ok := v["size"]; ok {
			fileInput.Size = size
		}
		if sizeBytes, ok := v["sizeBytes"].(float64); ok {
			fileInput.SizeBytes = int64(sizeBytes)
		}
	default:
		return nil, fmt.Errorf("video_file must be either a URL string or file object")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Get the video content
	var videoContent []byte

	// Method 1: Check if Src contains base64 data
	if fileInput.Src != "" {
		if strings.HasPrefix(fileInput.Src, "data:") {
			// Extract from data URL
			parts := strings.Split(fileInput.Src, ",")
			if len(parts) == 2 {
				decoded, err := base64.StdEncoding.DecodeString(parts[1])
				if err == nil {
					videoContent = decoded
				}
			}
		} else if !strings.HasPrefix(fileInput.Src, "http") {
			// Try as raw base64
			decoded, err := base64.StdEncoding.DecodeString(fileInput.Src)
			if err == nil {
				videoContent = decoded
			}
		}
	}

	// Method 2: Download from URL if no content yet
	if len(videoContent) == 0 {
		downloadURL := fileInput.GetDownloadURL()
		if downloadURL != "" {
			videoContent, err = downloadFileWithContext(ctx, downloadURL)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve video file: %v", err)
			}
		}
	}

	if len(videoContent) == 0 {
		return nil, errors.New("no video content could be retrieved")
	}

	// Get file size
	fileSize := int64(len(videoContent))
	if fileInput.SizeBytes > 0 {
		fileSize = fileInput.SizeBytes
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

	// Create a reader from the video content
	videoReader := bytes.NewReader(videoContent)

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

	// Create a progress reporter
	progressReader := &progressReader{
		Reader: videoReader,
		Total:  fileSize,
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

	// Add video to playlist if specified
	if input.PlaylistID != "" && response.Id != "" {
		playlistItem, err := addVideoToPlaylist(youtubeService, input.PlaylistID, response.Id)
		if err != nil {
			// Log the error but don't fail the entire operation since the video was uploaded successfully
			result["playlistError"] = fmt.Sprintf("Video uploaded successfully but failed to add to playlist: %v", err)
		} else {
			result["playlistItem"] = map[string]any{
				"id":         playlistItem.Id,
				"playlistId": playlistItem.Snippet.PlaylistId,
				"position":   playlistItem.Snippet.Position,
			}
		}
	}

	return result, nil
}

// addVideoToPlaylist adds an uploaded video to a playlist
func addVideoToPlaylist(service *youtube.Service, playlistID, videoID string) (*youtube.PlaylistItem, error) {
	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}

	call := service.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
	return call.Do()
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

	// Handle Google Drive URLs
	if strings.Contains(url, "drive.google.com") {
		// Convert sharing URL to direct download URL
		if strings.Contains(url, "/file/d/") {
			// Extract file ID from URL like: https://drive.google.com/file/d/FILE_ID/view
			parts := strings.Split(url, "/")
			for i, part := range parts {
				if part == "d" && i+1 < len(parts) {
					fileID := parts[i+1]
					url = fmt.Sprintf("https://drive.google.com/uc?export=download&id=%s", fileID)
					break
				}
			}
		}
	}

	client := &http.Client{
		Timeout: 5 * time.Minute, // Increased timeout for large video files
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects (default is 10)
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			// Copy headers to redirected request
			if len(via) > 0 {
				req.Header = via[len(via)-1].Header
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add common headers to handle various file hosting services
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d - %s", resp.StatusCode, string(body))
	}

	// Check if we got HTML instead of video content
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		return nil, fmt.Errorf("received HTML instead of video content - the URL may require authentication or is not a direct download link")
	}

	// Read the video content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if len(content) == 0 {
		return nil, fmt.Errorf("downloaded file is empty")
	}

	return content, nil
}

func NewUploadVideoAction() sdk.Action {
	return &UploadVideoAction{}
}
