package shared

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gookit/goutil/arrutil"
	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"

	sdkcore "github.com/wakflo/go-sdk/v2/core"
)

var (
	form = smartform.NewAuthForm("youtube-auth", "Youtube OAuth", smartform.AuthStrategyOAuth2)
	_    = form.
		OAuthField("oauth", "Youtube OAuth").
		AuthorizationURL("https://accounts.google.com/o/oauth2/auth").
		TokenURL("https://oauth2.googleapis.com/token").
		Scopes([]string{
			"https://www.googleapis.com/auth/youtube",
			"https://www.googleapis.com/auth/youtube.channel-memberships.creator",
			"https://www.googleapis.com/auth/youtube.force-ssl",
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/youtube.upload",
			"https://www.googleapis.com/auth/youtubepartner",
			"https://www.googleapis.com/auth/youtubepartner-channel-audit",
		}).
		Build()
)

var SharedYoutubeAuuth = form.Build()

// RegisterChannelProps adds a dynamic channel selector field to the form
func RegisterChannelProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getChannelID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Create HTTP client
		client := &http.Client{}

		url := "https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&mine=true&maxResults=50"

		// Create request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Wakflo-YouTube-Integration/1.0")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to read response body: %v", err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[YouTube] ERROR response body: %s", string(body))

			// Try to parse Google API error format
			var apiError struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Errors  []struct {
						Domain  string `json:"domain"`
						Reason  string `json:"reason"`
						Message string `json:"message"`
					} `json:"errors"`
				} `json:"error"`
			}

			if jsonErr := json.Unmarshal(body, &apiError); jsonErr == nil {
				log.Printf("[YouTube] Google API Error Details:")
				log.Printf("[YouTube]   Code: %d", apiError.Error.Code)
				log.Printf("[YouTube]   Message: %s", apiError.Error.Message)
				for i, e := range apiError.Error.Errors {
					log.Printf("[YouTube]   Error[%d] - Domain: %s, Reason: %s, Message: %s",
						i, e.Domain, e.Reason, e.Message)
				}
			}

			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, resp.Status)
		}

		log.Printf("[YouTube] Response body preview: %s", string(body[:min(len(body), 500)]))

		var channelList YouTubeChannelList
		err = json.Unmarshal(body, &channelList)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
			log.Printf("[YouTube] Raw response that failed to unmarshal: %s", string(body))
			return nil, err
		}

		items := arrutil.Map[YouTubeChannel, map[string]any](channelList.Items, func(input YouTubeChannel) (target map[string]any, find bool) {
			channelInfo := map[string]any{
				"id":   input.ID,
				"name": input.Snippet.Title,
			}

			// Add additional info if available
			if input.Snippet.CustomURL != "" {
				channelInfo["customUrl"] = input.Snippet.CustomURL
				log.Printf("[YouTube]   Custom URL: %s", input.Snippet.CustomURL)
			}

			if input.Statistics.SubscriberCount != "" {
				channelInfo["subscribers"] = input.Statistics.SubscriberCount
				log.Printf("[YouTube]   Subscribers: %s", input.Statistics.SubscriberCount)
			}

			return channelInfo, true
		})

		return ctx.Respond(items, len(items))
	}

	helpText := "Select a YouTube channel"
	if fieldName == "channel_id" {
		helpText = "Select the YouTube channel to use"
	}

	return form.SelectField(fieldName, label).
		Placeholder("Select a channel").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getChannelID)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText(helpText)
}

// RegisterPlaylistProps adds a dynamic playlist selector field to the form
func RegisterPlaylistProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getPlaylistID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		// Create HTTP client
		client := &http.Client{}

		// Build URL with query parameters
		baseURL := "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50"

		// Check if we should filter by channel
		if input.ChannelID != "" {
			baseURL += "&channelId=" + input.ChannelID
		} else {
			// Default to user's own playlists
			baseURL += "&mine=true"
		}

		// Create request
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to create request: %v", err)
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Wakflo-YouTube-Integration/1.0")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to read response body: %v", err)
			return nil, err
		}

		log.Printf("[YouTube] Response body length: %d bytes", len(body))

		if resp.StatusCode != http.StatusOK {
			log.Printf("[YouTube] ERROR response body: %s", string(body))

			// Try to parse Google API error format
			var apiError struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Errors  []struct {
						Domain       string `json:"domain"`
						Reason       string `json:"reason"`
						Message      string `json:"message"`
						LocationType string `json:"locationType"`
						Location     string `json:"location"`
					} `json:"errors"`
				} `json:"error"`
			}

			if jsonErr := json.Unmarshal(body, &apiError); jsonErr == nil {
				log.Printf("[YouTube] Google API Error Details:")
				log.Printf("[YouTube]   Code: %d", apiError.Error.Code)
				log.Printf("[YouTube]   Message: %s", apiError.Error.Message)
				for i, e := range apiError.Error.Errors {
					log.Printf("[YouTube]   Error[%d]:", i)
					log.Printf("[YouTube]     Domain: %s", e.Domain)
					log.Printf("[YouTube]     Reason: %s", e.Reason)
					log.Printf("[YouTube]     Message: %s", e.Message)
					if e.LocationType != "" {
						log.Printf("[YouTube]     LocationType: %s", e.LocationType)
					}
					if e.Location != "" {
						log.Printf("[YouTube]     Location: %s", e.Location)
					}
				}
			}

			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, resp.Status)
		}

		var playlistList YouTubePlaylistList
		err = json.Unmarshal(body, &playlistList)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
			log.Printf("[YouTube] Raw response that failed to unmarshal: %s", string(body))
			return nil, err
		}

		items := arrutil.Map[YouTubePlaylist, map[string]any](playlistList.Items, func(input YouTubePlaylist) (target map[string]any, find bool) {
			playlistInfo := map[string]any{
				"id":   input.ID,
				"name": input.Snippet.Title,
			}

			// Add video count if available
			if input.ContentDetails.ItemCount > 0 {
				playlistInfo["videoCount"] = input.ContentDetails.ItemCount
				log.Printf("[YouTube]   Video count: %d", input.ContentDetails.ItemCount)
			}

			// Add channel title for context
			if input.Snippet.ChannelTitle != "" {
				playlistInfo["channel"] = input.Snippet.ChannelTitle
				log.Printf("[YouTube]   Channel: %s", input.Snippet.ChannelTitle)
			}

			return playlistInfo, true
		})

		return ctx.Respond(items, len(items))
	}

	builder := form.SelectField(fieldName, label).
		Placeholder("Select a playlist").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPlaylistID)).
				WithFieldReference("channel_id", "channel_id").
				WithSearchSupport().
				End().
				RefreshOn("channel_id").
				GetDynamicSource(),
		).
		Required(required).
		HelpText("Select a YouTube playlist")

	return builder
}

// RegisterVideoProps adds a dynamic video selector field to the form
func RegisterVideoProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getVideoID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		// Create HTTP client
		client := &http.Client{}

		// Build URL with query parameters
		baseURL := "https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&maxResults=50&order=date"

		// Filter by channel if provided
		if input.ChannelID != "" {
			baseURL += "&channelId=" + input.ChannelID
		} else {
			// Default to user's own videos
			baseURL += "&forMine=true"
		}

		// Create request
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to create request: %v", err)
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Wakflo-YouTube-Integration/1.0")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to read response body: %v", err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[YouTube] ERROR response body: %s", string(body))
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, resp.Status)
		}

		var videoList YouTubeVideoList
		err = json.Unmarshal(body, &videoList)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
			return nil, err
		}

		items := arrutil.Map[YouTubeVideo, map[string]any](videoList.Items, func(input YouTubeVideo) (target map[string]any, find bool) {
			videoInfo := map[string]any{
				"id":   input.ID.VideoID,
				"name": input.Snippet.Title,
			}

			// Add published date if available
			if input.Snippet.PublishedAt != "" {
				videoInfo["publishedAt"] = input.Snippet.PublishedAt
			}

			// Add channel title for context
			if input.Snippet.ChannelTitle != "" {
				videoInfo["channel"] = input.Snippet.ChannelTitle
			}

			return videoInfo, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.SelectField(fieldName, label).
		Placeholder("Select a video").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getVideoID)).
				WithFieldReference("channel_id", "channel_id").
				WithSearchSupport().
				End().
				RefreshOn("channel_id").
				GetDynamicSource(),
		).
		HelpText("Select a YouTube video")
}

// RegisterMultiPlaylistProps adds a multi-select dynamic playlist field to the form
func RegisterMultiPlaylistProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getPlaylistID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		// Create HTTP client
		client := &http.Client{}

		// Build URL with query parameters
		baseURL := "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50"

		// Check if we should filter by channel
		if input.ChannelID != "" {
			baseURL += "&channelId=" + input.ChannelID
		} else {
			// Default to user's own playlists
			baseURL += "&mine=true"
		}

		// Create request
		req, err := http.NewRequest("GET", baseURL, nil)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to create request: %v", err)
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Wakflo-YouTube-Integration/1.0")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to read response body: %v", err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("[YouTube] ERROR response body: %s", string(body))
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, resp.Status)
		}

		var playlistList YouTubePlaylistList
		err = json.Unmarshal(body, &playlistList)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
			return nil, err
		}

		items := arrutil.Map[YouTubePlaylist, map[string]any](playlistList.Items, func(input YouTubePlaylist) (target map[string]any, find bool) {
			playlistInfo := map[string]any{
				"id":   input.ID,
				"name": input.Snippet.Title,
			}

			// Add video count if available
			if input.ContentDetails.ItemCount > 0 {
				playlistInfo["videoCount"] = input.ContentDetails.ItemCount
			}

			// Add channel title for context
			if input.Snippet.ChannelTitle != "" {
				playlistInfo["channel"] = input.Snippet.ChannelTitle
			}

			return playlistInfo, true
		})

		return ctx.Respond(items, len(items))
	}

	return form.MultiSelectField(fieldName, label).
		Placeholder("Select playlists").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getPlaylistID)).
				WithFieldReference("channel_id", "channel_id").
				WithSearchSupport().
				End().
				RefreshOn("channel_id").
				GetDynamicSource(),
		).
		HelpText("Select one or more YouTube playlists")
}

// Helper function for min since Go doesn't have a built-in for ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
