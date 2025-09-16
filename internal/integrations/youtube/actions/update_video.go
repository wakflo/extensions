package actions

import (
	"context"
	"errors"
	"fmt"
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
	VideoID       string   `json:"video_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Tags          string   `json:"tags"`
	CategoryID    string   `json:"category_id"`
	PrivacyStatus string   `json:"privacy_status"`
	PlaylistIDs   []string `json:"playlist_ids"`
}

type UpdateVideoAction struct{}

func (a *UpdateVideoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_update_video",
		DisplayName:   "Update YouTube Video",
		Description:   "Update metadata for a YouTube video including title, description, tags, privacy settings, and manage playlists.",
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
				"privacyStatus": "public",
			},
			"playlists": map[string]any{
				"added":   []string{"PLxxxxx1"},
				"current": []string{"PLxxxxx1", "PLxxxxx2"},
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

	form.TextField("title", "Title").
		Placeholder("Enter video title").
		HelpText("The video's title (max 100 characters)").
		Required(false)

	form.TextareaField("description", "Description").
		Placeholder("Enter video description").
		HelpText("The video's description (max 5000 characters)").
		Required(false)

	shared.RegisterMultiPlaylistProps(form, "playlist_ids", "Add to Playlists", false).
		HelpText("Select playlists to add this video to")

	form.TextareaField("tags", "Tags").
		Placeholder("tag1, tag2, tag3").
		HelpText("Comma-separated list of tags associated with the video").
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
		HelpText("Select a category for the video").
		Required(false)

	form.SelectField("privacy_status", "Privacy Status").
		AddOption("private", "Private").
		AddOption("unlisted", "Unlisted").
		AddOption("public", "Public").
		HelpText("Set the video's privacy status").
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

	// Initialize result
	result := map[string]any{
		"id": input.VideoID,
	}

	// Check if we need to update video metadata
	needsUpdate := false
	var video *youtube.Video

	// First, get the current video if we need to update it
	if input.Title != "" || input.Description != "" || input.Tags != "" || input.CategoryID != "" || input.PrivacyStatus != "" {
		// Get current video
		videoListCall := youtubeService.Videos.List([]string{"snippet", "status"})
		videoListCall = videoListCall.Id(input.VideoID)

		currentVideo, err := videoListCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch video: %w", err)
		}

		if len(currentVideo.Items) == 0 {
			return nil, errors.New("video not found")
		}

		video = currentVideo.Items[0]
		needsUpdate = true

		// Update fields if provided
		if input.Title != "" {
			video.Snippet.Title = input.Title
		}
		if input.Description != "" {
			video.Snippet.Description = input.Description
		}
		if input.Tags != "" {
			tags := strings.Split(input.Tags, ",")
			var cleanTags []string
			for _, tag := range tags {
				trimmed := strings.TrimSpace(tag)
				if trimmed != "" {
					cleanTags = append(cleanTags, trimmed)
				}
			}
			video.Snippet.Tags = cleanTags
		}
		if input.CategoryID != "" {
			video.Snippet.CategoryId = input.CategoryID
		}
		if input.PrivacyStatus != "" {
			video.Status.PrivacyStatus = input.PrivacyStatus
		}
	}

	// Update video if needed
	if needsUpdate {
		parts := []string{"snippet", "status"}
		updateCall := youtubeService.Videos.Update(parts, video)

		updatedVideo, err := updateCall.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to update video: %w", err)
		}

		// Add updated fields to result
		result["title"] = updatedVideo.Snippet.Title
		result["description"] = updatedVideo.Snippet.Description
		result["tags"] = updatedVideo.Snippet.Tags
		result["categoryId"] = updatedVideo.Snippet.CategoryId
		result["status"] = map[string]any{
			"privacyStatus": updatedVideo.Status.PrivacyStatus,
		}
	}

	// Handle playlist updates
	playlistResults := map[string]any{
		"added":   []string{},
		"current": []string{},
	}

	if len(input.PlaylistIDs) > 0 {
		// Add video to playlists
		for _, playlistID := range input.PlaylistIDs {
			playlistItem := &youtube.PlaylistItem{
				Snippet: &youtube.PlaylistItemSnippet{
					PlaylistId: playlistID,
					ResourceId: &youtube.ResourceId{
						Kind:    "youtube#video",
						VideoId: input.VideoID,
					},
				},
			}

			insertCall := youtubeService.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
			_, err := insertCall.Do()
			if err != nil {
				// Log error but continue with other playlists
				ctx.Logger().Warn("failed to add video to playlist", "playlist_id", playlistID, "error", err)
			} else {
				playlistResults["added"] = append(playlistResults["added"].([]string), playlistID)
			}
		}

		playlistResults["current"] = input.PlaylistIDs
	}

	result["playlists"] = playlistResults

	return result, nil
}

func NewUpdateVideoAction() sdk.Action {
	return &UpdateVideoAction{}
}
