package actions

import (
	"context"
	"errors"
	"strings"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/youtube/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type updateVideoActionProps struct {
	VideoID                 string `json:"video_id"`
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
	DefaultLanguage         string `json:"default_language"`
	DefaultAudioLanguage    string `json:"default_audio_language"`
}

type UpdateVideoAction struct{}

func (a *UpdateVideoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_update_video",
		DisplayName:   "Update YouTube Video",
		Description:   "Update metadata for a YouTube video including title, description, tags, privacy settings, and more.",
		Type:          core.ActionTypeAction,
		Documentation: updateVideoDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"id":          "dQw4w9WgXcQ",
			"title":       "Updated Video Title",
			"description": "Updated video description",
			"tags":        []string{"tag1", "tag2", "tag3"},
			"categoryId":  "22",
			"status": map[string]any{
				"privacyStatus":       "public",
				"embeddable":          true,
				"publicStatsViewable": true,
				"madeForKids":         false,
				"license":             "youtube",
			},
			"snippet": map[string]any{
				"publishedAt":          "2023-01-01T00:00:00Z",
				"channelId":            "UCuAXFkgsw1L7xaCfnd5JJOw",
				"channelTitle":         "Channel Name",
				"defaultLanguage":      "en",
				"defaultAudioLanguage": "en",
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *UpdateVideoAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("youtube_update_video", "Update YouTube Video")

	// Channel selection
	shared.RegisterChannelProps(form, "channel_id", "Channel", true).
		HelpText("Select the channel containing the video you want to update")

	// Video selection - refreshes when channel changes
	shared.RegisterVideoProps(form, "video_id", "Video to Update", true).
		HelpText("Select the video you want to update")

	// Basic Information Section
	form.SectionField("basicInfo", "Basic Information")

	form.TextField("title", "Title").
		Placeholder("Enter video title").
		HelpText("The video's title (max 100 characters)").
		Required(false)

	form.TextareaField("description", "Description").
		Placeholder("Enter video description").
		HelpText("The video's description (max 5000 characters)").
		Required(false)

	form.TextareaField("tags", "Tags").
		Placeholder("tag1, tag2, tag3").
		HelpText("Comma-separated list of tags associated with the video").
		Required(false)

	// Category selection with common YouTube categories
	form.SelectField("category_id", "Category").
		// AddOption("", "Keep current category").
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
		HelpText("Select a category for the video").
		Required(false)

	// Privacy and Status Section
	form.SectionField("privacyStatus", "Privacy and Status Settings")

	form.SelectField("privacy_status", "Privacy Status").
		// AddOption("", "Keep current privacy").
		AddOption("private", "Private").
		AddOption("unlisted", "Unlisted").
		AddOption("public", "Public").
		HelpText("Set the video's privacy status").
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
		HelpText("Indicates whether the video is designated as child-directed").
		Required(false)

	form.CheckboxField("self_declared_made_for_kids", "Self-Declared Made for Kids").
		DefaultValue(false).
		HelpText("Channel owner's designation of whether the video is child-directed").
		Required(false)

	form.SelectField("license", "License").
		// AddOption("", "Keep current license").
		AddOption("youtube", "Standard YouTube License").
		AddOption("creativeCommon", "Creative Commons - Attribution").
		HelpText("The video's license").
		Required(false)

	// Additional Settings Section
	form.SectionField("additionalSettings", "Additional Settings")

	// Remove recording date field as it's not updatable via API
	// form.DateTimeField("recording_date", "Recording Date").
	//	Placeholder("2023-01-01T00:00:00Z").
	//	HelpText("The date and time when the video was recorded (RFC 3339 format)").
	//	Required(false)

	form.TextField("default_language", "Default Language").
		Placeholder("en").
		HelpText("The language of the video's default audio track (ISO 639-1 two-letter code)").
		Required(false)

	form.TextField("default_audio_language", "Default Audio Language").
		Placeholder("en").
		HelpText("The language of the video's default audio track (ISO 639-1 two-letter code)").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *UpdateVideoAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *UpdateVideoAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[updateVideoActionProps](ctx)
	if err != nil {
		return nil, err
	}

	if input.VideoID == "" {
		return nil, errors.New("video ID is required")
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// First, get the current video to preserve existing data
	// Remove recordingDetails from the list call since we won't be updating it
	currentVideo, err := youtubeService.Videos.List([]string{"snippet", "status"}).
		Id(input.VideoID).Do()
	if err != nil {
		return nil, err
	}

	if len(currentVideo.Items) == 0 {
		return nil, errors.New("video not found")
	}

	video := currentVideo.Items[0]

	// Update snippet if any snippet fields are provided
	updateSnippet := false
	if input.Title != "" {
		video.Snippet.Title = input.Title
		updateSnippet = true
	}
	if input.Description != "" {
		video.Snippet.Description = input.Description
		updateSnippet = true
	}
	if input.Tags != "" {
		// Split tags by comma and trim whitespace
		tags := strings.Split(input.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		video.Snippet.Tags = tags
		updateSnippet = true
	}
	if input.CategoryID != "" {
		video.Snippet.CategoryId = input.CategoryID
		updateSnippet = true
	}
	if input.DefaultLanguage != "" {
		video.Snippet.DefaultLanguage = input.DefaultLanguage
		updateSnippet = true
	}
	if input.DefaultAudioLanguage != "" {
		video.Snippet.DefaultAudioLanguage = input.DefaultAudioLanguage
		updateSnippet = true
	}

	// Update status if any status fields are provided
	updateStatus := false
	if input.PrivacyStatus != "" {
		video.Status.PrivacyStatus = input.PrivacyStatus
		updateStatus = true
	}
	// For boolean fields, we always update them since false is a valid value
	video.Status.Embeddable = input.Embeddable
	video.Status.PublicStatsViewable = input.PublicStatsViewable
	video.Status.MadeForKids = input.MadeForKids
	video.Status.SelfDeclaredMadeForKids = input.SelfDeclaredMadeForKids
	updateStatus = true

	if input.License != "" {
		video.Status.License = input.License
		updateStatus = true
	}

	// Remove recording details update logic entirely
	// Recording details cannot be updated via the YouTube API

	// Build the parts list based on what we're updating
	var parts []string
	if updateSnippet {
		parts = append(parts, "snippet")
	}
	if updateStatus {
		parts = append(parts, "status")
	}

	if len(parts) == 0 {
		return nil, errors.New("no fields to update")
	}

	// Make the update call
	updateCall := youtubeService.Videos.Update(parts, video)
	updatedVideo, err := updateCall.Do()
	if err != nil {
		return nil, err
	}

	// Build the response
	result := map[string]any{
		"id":          updatedVideo.Id,
		"title":       updatedVideo.Snippet.Title,
		"description": updatedVideo.Snippet.Description,
		"tags":        updatedVideo.Snippet.Tags,
		"categoryId":  updatedVideo.Snippet.CategoryId,
		"snippet": map[string]any{
			"publishedAt":  updatedVideo.Snippet.PublishedAt,
			"channelId":    updatedVideo.Snippet.ChannelId,
			"channelTitle": updatedVideo.Snippet.ChannelTitle,
		},
		"status": map[string]any{
			"uploadStatus":            updatedVideo.Status.UploadStatus,
			"privacyStatus":           updatedVideo.Status.PrivacyStatus,
			"license":                 updatedVideo.Status.License,
			"embeddable":              updatedVideo.Status.Embeddable,
			"publicStatsViewable":     updatedVideo.Status.PublicStatsViewable,
			"madeForKids":             updatedVideo.Status.MadeForKids,
			"selfDeclaredMadeForKids": updatedVideo.Status.SelfDeclaredMadeForKids,
		},
	}

	// Add language info if available
	if updatedVideo.Snippet.DefaultLanguage != "" {
		result["defaultLanguage"] = updatedVideo.Snippet.DefaultLanguage
	}
	if updatedVideo.Snippet.DefaultAudioLanguage != "" {
		result["defaultAudioLanguage"] = updatedVideo.Snippet.DefaultAudioLanguage
	}

	// Remove recording date from response since it's not updatable
	// Recording details would need to be fetched separately if needed

	return result, nil
}

func NewUpdateVideoAction() sdk.Action {
	return &UpdateVideoAction{}
}
