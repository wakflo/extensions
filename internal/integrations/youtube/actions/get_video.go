package actions

import (
	"context"
	"errors"
	"log"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/extensions/internal/integrations/youtube/shared"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type getVideoActionProps struct {
	SearchOwnChannel bool   `json:"search_own_channel"`
	VideoID          string `json:"video_id"`
	OwnVideoID       string `json:"own_video_id"`
	Parts            string `json:"parts"`
}

type GetVideoAction struct{}

func (a *GetVideoAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_get_video",
		DisplayName:   "Get YouTube Video",
		Description:   "Get detailed information about a specific YouTube video by its ID. You can search for any public video or select from your own channel's videos.",
		Type:          core.ActionTypeAction,
		Documentation: getVideoDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"id":          "dQw4w9WgXcQ",
			"title":       "Sample Video Title",
			"description": "Sample video description",
			"channel": map[string]any{
				"id":    "UCuAXFkgsw1L7xaCfnd5JJOw",
				"title": "Sample Channel",
			},
			"publishedAt":  "2023-01-01T00:00:00Z",
			"duration":     "PT3M32S",
			"viewCount":    "1000000",
			"likeCount":    "50000",
			"commentCount": "5000",
			"thumbnails": map[string]any{
				"default": map[string]any{
					"url":    "https://i.ytimg.com/vi/dQw4w9WgXcQ/default.jpg",
					"width":  120,
					"height": 90,
				},
			},
		},
		Settings: core.ActionSettings{},
	}
}

func (a *GetVideoAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("youtube_get_video", "Get YouTube Video")

	// Add checkbox to toggle between own channel and public videos
	form.CheckboxField("search_own_channel", "Search Only My Channel").
		DefaultValue(false).
		HelpText("Check this to search only within your own channel's videos").
		Required(false)

	// Channel selector - visible when search_own_channel is true
	shared.RegisterChannelProps(form, "channel_id", "Channel", false).
		VisibleWhenEquals("search_own_channel", true)

	// Public video ID field - visible when search_own_channel is false
	form.TextField("video_id", "Video ID").
		Placeholder("dQw4w9WgXcQ").
		HelpText("Enter the YouTube video ID (e.g., from youtube.com/watch?v=VIDEO_ID)").
		Required(true).
		VisibleWhenEquals("search_own_channel", false)

	// Dynamic video selector for own channel - visible when search_own_channel is true
	shared.RegisterVideoProps(form, "own_video_id", "Select Video", true).
		VisibleWhenEquals("search_own_channel", true)

	// Data parts selection
	form.SelectField("parts", "Data Parts").
		AddOption("basic", "Basic (snippet, statistics)").
		AddOption("detailed", "Detailed (snippet, statistics, contentDetails, status)").
		AddOption("full", "Full (all available data)").
		DefaultValue("detailed").
		HelpText("Select which parts of the video data to retrieve").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *GetVideoAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *GetVideoAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[getVideoActionProps](ctx)
	if err != nil {
		return nil, err
	}

	// Determine which video ID to use
	var videoID string
	if input.SearchOwnChannel {
		if input.OwnVideoID == "" {
			return nil, errors.New("please select a video from your channel")
		}
		videoID = input.OwnVideoID
	} else {
		if input.VideoID == "" {
			return nil, errors.New("video ID is required")
		}
		videoID = input.VideoID
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Determine which parts to fetch
	parts := determineParts(input.Parts)

	// Make the API call
	call := youtubeService.Videos.List(parts).Id(videoID)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, errors.New("video not found")
	}

	video := response.Items[0]

	// Build the response - following the same pattern as list_videos.go
	videoData := map[string]any{
		"id": video.Id,
	}

	// Add snippet data
	if video.Snippet != nil {
		videoData["title"] = video.Snippet.Title
		videoData["description"] = video.Snippet.Description
		videoData["channel"] = map[string]any{
			"id":    video.Snippet.ChannelId,
			"title": video.Snippet.ChannelTitle,
		}
		videoData["publishedAt"] = video.Snippet.PublishedAt
		videoData["tags"] = video.Snippet.Tags
		videoData["categoryId"] = video.Snippet.CategoryId
		videoData["liveBroadcastContent"] = video.Snippet.LiveBroadcastContent

		if video.Snippet.DefaultLanguage != "" {
			videoData["defaultLanguage"] = video.Snippet.DefaultLanguage
		}
		if video.Snippet.DefaultAudioLanguage != "" {
			videoData["defaultAudioLanguage"] = video.Snippet.DefaultAudioLanguage
		}

		// Add thumbnails
		if video.Snippet.Thumbnails != nil {
			thumbnails := make(map[string]any)
			if video.Snippet.Thumbnails.Default != nil {
				thumbnails["default"] = map[string]any{
					"url":    video.Snippet.Thumbnails.Default.Url,
					"width":  video.Snippet.Thumbnails.Default.Width,
					"height": video.Snippet.Thumbnails.Default.Height,
				}
			}
			if video.Snippet.Thumbnails.Medium != nil {
				thumbnails["medium"] = map[string]any{
					"url":    video.Snippet.Thumbnails.Medium.Url,
					"width":  video.Snippet.Thumbnails.Medium.Width,
					"height": video.Snippet.Thumbnails.Medium.Height,
				}
			}
			if video.Snippet.Thumbnails.High != nil {
				thumbnails["high"] = map[string]any{
					"url":    video.Snippet.Thumbnails.High.Url,
					"width":  video.Snippet.Thumbnails.High.Width,
					"height": video.Snippet.Thumbnails.High.Height,
				}
			}
			if video.Snippet.Thumbnails.Standard != nil {
				thumbnails["standard"] = map[string]any{
					"url":    video.Snippet.Thumbnails.Standard.Url,
					"width":  video.Snippet.Thumbnails.Standard.Width,
					"height": video.Snippet.Thumbnails.Standard.Height,
				}
			}
			if video.Snippet.Thumbnails.Maxres != nil {
				thumbnails["maxres"] = map[string]any{
					"url":    video.Snippet.Thumbnails.Maxres.Url,
					"width":  video.Snippet.Thumbnails.Maxres.Width,
					"height": video.Snippet.Thumbnails.Maxres.Height,
				}
			}
			videoData["thumbnails"] = thumbnails
		}
	}

	// Add statistics if available
	if video.Statistics != nil {
		videoData["viewCount"] = video.Statistics.ViewCount
		videoData["likeCount"] = video.Statistics.LikeCount
		videoData["commentCount"] = video.Statistics.CommentCount
	}

	// Add content details
	if video.ContentDetails != nil {
		videoData["duration"] = video.ContentDetails.Duration
		videoData["definition"] = video.ContentDetails.Definition
		videoData["caption"] = video.ContentDetails.Caption
		videoData["licensedContent"] = video.ContentDetails.LicensedContent
		videoData["projection"] = video.ContentDetails.Projection

		if video.ContentDetails.Dimension != "" {
			videoData["dimension"] = video.ContentDetails.Dimension
		}
		if video.ContentDetails.ContentRating != nil {
			videoData["contentRating"] = video.ContentDetails.ContentRating
		}
		if video.ContentDetails.RegionRestriction != nil {
			regionRestriction := map[string]any{}
			if len(video.ContentDetails.RegionRestriction.Allowed) > 0 {
				regionRestriction["allowed"] = video.ContentDetails.RegionRestriction.Allowed
			}
			if len(video.ContentDetails.RegionRestriction.Blocked) > 0 {
				regionRestriction["blocked"] = video.ContentDetails.RegionRestriction.Blocked
			}
			videoData["regionRestriction"] = regionRestriction
		}
	}

	// Add status
	if video.Status != nil {
		status := map[string]any{
			"uploadStatus":            video.Status.UploadStatus,
			"privacyStatus":           video.Status.PrivacyStatus,
			"license":                 video.Status.License,
			"embeddable":              video.Status.Embeddable,
			"publicStatsViewable":     video.Status.PublicStatsViewable,
			"madeForKids":             video.Status.MadeForKids,
			"selfDeclaredMadeForKids": video.Status.SelfDeclaredMadeForKids,
		}
		videoData["status"] = status
	}

	// Add player embed HTML
	if video.Player != nil {
		videoData["player"] = map[string]any{
			"embedHtml": video.Player.EmbedHtml,
		}
	}

	// Add topic details
	if video.TopicDetails != nil && len(video.TopicDetails.TopicCategories) > 0 {
		videoData["topicDetails"] = map[string]any{
			"topicCategories": video.TopicDetails.TopicCategories,
		}
	}

	// Add recording details
	if video.RecordingDetails != nil {
		recordingDetails := map[string]any{}
		if video.RecordingDetails.RecordingDate != "" {
			recordingDetails["recordingDate"] = video.RecordingDetails.RecordingDate
		}
		if video.RecordingDetails.Location != nil {
			location := map[string]any{
				"latitude":  video.RecordingDetails.Location.Latitude,
				"longitude": video.RecordingDetails.Location.Longitude,
				"altitude":  video.RecordingDetails.Location.Altitude,
			}
			recordingDetails["location"] = location
		}
		if len(recordingDetails) > 0 {
			videoData["recordingDetails"] = recordingDetails
		}
	}

	// Add localization data
	if video.Localizations != nil {
		localizations := make(map[string]any)
		for lang, loc := range video.Localizations {
			localizations[lang] = map[string]any{
				"title":       loc.Title,
				"description": loc.Description,
			}
		}
		if len(localizations) > 0 {
			videoData["localizations"] = localizations
		}
	}

	// If searching own channel, verify the video belongs to the authenticated user
	if input.SearchOwnChannel {
		log.Printf("[YouTube] Retrieved video from own channel: %s", video.Snippet.Title)
	}

	return videoData, nil
}

func determineParts(partsOption string) []string {
	switch partsOption {
	case "basic":
		return []string{"snippet", "statistics"}
	case "detailed":
		return []string{"snippet", "statistics", "contentDetails", "status"}
	case "full":
		return []string{
			"snippet",
			"statistics",
			"contentDetails",
			"status",
			"player",
			"topicDetails",
			"recordingDetails",
			"localizations",
		}
	default:
		// Default to detailed
		return []string{"snippet", "statistics", "contentDetails", "status"}
	}
}

func NewGetVideoAction() sdk.Action {
	return &GetVideoAction{}
}
