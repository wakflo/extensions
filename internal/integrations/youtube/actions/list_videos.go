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

type listVideosActionProps struct {
	SearchQuery     string `json:"search_query"`
	ChannelID       string `json:"channel_id"`
	PlaylistID      string `json:"playlist_id"`
	VideoIDs        string `json:"video_ids"`
	MaxResults      int64  `json:"max_results"`
	Order           string `json:"order"`
	VideoDuration   string `json:"video_duration"`
	VideoType       string `json:"video_type"`
	PublishedAfter  string `json:"published_after"`
	PublishedBefore string `json:"published_before"`
}

type ListVideosAction struct{}

func (a *ListVideosAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "youtube_list_videos",
		DisplayName:   "List YouTube Videos",
		Description:   "List YouTube Videos: Search and retrieve videos from YouTube based on various criteria including search query, channel ID, playlist ID, or specific video IDs.",
		Type:          core.ActionTypeAction,
		Documentation: listVideosDocs,
		Icon:          "youtube",
		SampleOutput: map[string]any{
			"videos": []map[string]any{
				{
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
			},
			"nextPageToken": "CAoQAA",
			"totalResults":  100,
		},
		Settings: core.ActionSettings{},
	}
}

func (a *ListVideosAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("youtube_list_videos", "List YouTube Videos")

	// Search parameters
	form.TextField("search_query", "Search Query").
		Placeholder("Enter search keywords").
		HelpText("Keywords to search for videos. Leave empty to list without search.").
		Required(false)

	shared.RegisterChannelProps(form, "channel_id", "Channel ID", false)

	shared.RegisterPlaylistProps(form, "playlist_id", "Playlist ID", true)

	form.TextField("video_ids", "Video IDs").
		Placeholder("videoId1,videoId2,videoId3").
		HelpText("Comma-separated list of specific video IDs to retrieve").
		Required(false)

	// Result configuration
	form.NumberField("max_results", "Max Results").
		DefaultValue(25).
		HelpText("Maximum number of videos to return (1-50)").
		Required(false)

	form.SelectField("order", "Sort Order").
		AddOption("relevance", "Relevance").
		AddOption("date", "Date").
		AddOption("rating", "Rating").
		AddOption("title", "Title").
		AddOption("videoCount", "Video Count").
		AddOption("viewCount", "View Count").
		HelpText("Order to sort the results").
		Required(false)

	// Filters
	form.SelectField("video_duration", "Video Duration").
		AddOption("any", "Any").
		AddOption("short", "Short (< 4 minutes)").
		AddOption("medium", "Medium (4-20 minutes)").
		AddOption("long", "Long (> 20 minutes)").
		HelpText("Filter by video duration").
		Required(false)

	form.SelectField("video_type", "Video Type").
		AddOption("any", "Any").
		AddOption("episode", "Episode").
		AddOption("movie", "Movie").
		HelpText("Filter by video type").
		Required(false)

	form.DateTimeField("published_after", "Published After").
		Placeholder("2023-01-01T00:00:00Z").
		HelpText("Filter videos published after this date (RFC 3339 format)").
		Required(false)

	form.DateTimeField("published_before", "Published Before").
		Placeholder("2023-12-31T23:59:59Z").
		HelpText("Filter videos published before this date (RFC 3339 format)").
		Required(false)

	schema := form.Build()

	return schema
}

func (a *ListVideosAction) Auth() *core.AuthMetadata {
	return nil
}

func (a *ListVideosAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[listVideosActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(*authCtx.TokenSource))
	if err != nil {
		return nil, err
	}

	// Validate that at least one search criterion is provided
	if input.SearchQuery == "" && input.ChannelID == "" && input.PlaylistID == "" && input.VideoIDs == "" {
		return nil, errors.New("at least one search criterion must be provided (search query, channel ID, playlist ID, or video IDs)")
	}

	// Set default max results if not provided
	if input.MaxResults == 0 {
		input.MaxResults = 25
	}

	var videos []map[string]any
	var nextPageToken string
	var totalResults int64

	// Handle different listing scenarios
	if input.VideoIDs != "" {
		// Get specific videos by IDs
		videoIDs := strings.Split(strings.TrimSpace(input.VideoIDs), ",")
		videos, err = getVideosByIDs(youtubeService, videoIDs)
		if err != nil {
			return nil, err
		}
		totalResults = int64(len(videos))
	} else if input.PlaylistID != "" {
		// Get videos from playlist
		videos, nextPageToken, err = getPlaylistVideos(youtubeService, input.PlaylistID, input.MaxResults)
		if err != nil {
			return nil, err
		}
	} else {
		// Search for videos
		videos, nextPageToken, totalResults, err = searchVideos(youtubeService, input)
		if err != nil {
			return nil, err
		}
	}

	result := map[string]any{
		"videos":       videos,
		"totalResults": totalResults,
	}

	if nextPageToken != "" {
		result["nextPageToken"] = nextPageToken
	}

	return result, nil
}

func searchVideos(service *youtube.Service, input *listVideosActionProps) ([]map[string]any, string, int64, error) {
	searchCall := service.Search.List([]string{"id", "snippet"}).
		Q(input.SearchQuery).
		Type("video").
		MaxResults(input.MaxResults)

	// Apply filters
	if input.ChannelID != "" {
		searchCall = searchCall.ChannelId(input.ChannelID)
	}

	if input.Order != "" {
		searchCall = searchCall.Order(input.Order)
	}

	if input.VideoDuration != "" {
		searchCall = searchCall.VideoDuration(input.VideoDuration)
	}

	if input.VideoType != "" {
		searchCall = searchCall.VideoDefinition(input.VideoType)
	}

	if input.PublishedAfter != "" {
		searchCall = searchCall.PublishedAfter(input.PublishedAfter)
	}

	if input.PublishedBefore != "" {
		searchCall = searchCall.PublishedBefore(input.PublishedBefore)
	}

	searchResponse, err := searchCall.Do()
	if err != nil {
		return nil, "", 0, err
	}

	// Get video IDs from search results
	var videoIDs []string
	for _, item := range searchResponse.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
	}

	// Get detailed video information
	videos, err := getVideosByIDs(service, videoIDs)
	if err != nil {
		return nil, "", 0, err
	}

	return videos, searchResponse.NextPageToken, searchResponse.PageInfo.TotalResults, nil
}

func getPlaylistVideos(service *youtube.Service, playlistID string, maxResults int64) ([]map[string]any, string, error) {
	playlistCall := service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
		PlaylistId(playlistID).
		MaxResults(maxResults)

	playlistResponse, err := playlistCall.Do()
	if err != nil {
		return nil, "", err
	}

	// Get video IDs from playlist
	var videoIDs []string
	for _, item := range playlistResponse.Items {
		videoIDs = append(videoIDs, item.ContentDetails.VideoId)
	}

	// Get detailed video information
	videos, err := getVideosByIDs(service, videoIDs)
	if err != nil {
		return nil, "", err
	}

	return videos, playlistResponse.NextPageToken, nil
}

func getVideosByIDs(service *youtube.Service, videoIDs []string) ([]map[string]any, error) {
	if len(videoIDs) == 0 {
		return []map[string]any{}, nil
	}

	videosCall := service.Videos.List([]string{"snippet", "contentDetails", "statistics"}).
		Id(strings.Join(videoIDs, ","))

	videosResponse, err := videosCall.Do()
	if err != nil {
		return nil, err
	}

	var videos []map[string]any
	for _, video := range videosResponse.Items {
		videoData := map[string]any{
			"id":          video.Id,
			"title":       video.Snippet.Title,
			"description": video.Snippet.Description,
			"channel": map[string]any{
				"id":    video.Snippet.ChannelId,
				"title": video.Snippet.ChannelTitle,
			},
			"publishedAt": video.Snippet.PublishedAt,
			"duration":    video.ContentDetails.Duration,
			"tags":        video.Snippet.Tags,
		}

		// Add statistics if available
		if video.Statistics != nil {
			videoData["viewCount"] = video.Statistics.ViewCount
			videoData["likeCount"] = video.Statistics.LikeCount
			videoData["commentCount"] = video.Statistics.CommentCount
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
			videoData["thumbnails"] = thumbnails
		}

		videos = append(videos, videoData)
	}

	return videos, nil
}

func NewListVideosAction() sdk.Action {
	return &ListVideosAction{}
}
