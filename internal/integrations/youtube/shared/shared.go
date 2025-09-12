package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
		log.Println("[YouTube] Starting getChannelID function")

		authCtx, err := ctx.AuthContext()
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to get auth context: %v", err)
			return nil, err
		}

		log.Printf("[YouTube] Auth context obtained successfully")
		log.Printf("[YouTube] Access token length: %d", len(authCtx.AccessToken))
		log.Printf("[YouTube] Access token preview: %s...", authCtx.AccessToken[:20])

		// Create HTTP client
		client := &http.Client{}

		url := "https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&mine=true&maxResults=50"
		log.Printf("[YouTube] Making GET request to: %s", url)

		// Create request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to create request: %v", err)
			return nil, err
		}

		// Set headers
		req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "Wakflo-YouTube-Integration/1.0")

		// Log request headers
		log.Println("[YouTube] Request headers:")
		for key, values := range req.Header {
			for _, value := range values {
				if key == "Authorization" {
					log.Printf("[YouTube]   %s: Bearer %s...", key, value[7:27]) // Log partial token
				} else {
					log.Printf("[YouTube]   %s: %s", key, value)
				}
			}
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		log.Printf("[YouTube] Response status code: %d", resp.StatusCode)
		log.Printf("[YouTube] Response status: %s", resp.Status)

		// Log response headers
		log.Println("[YouTube] Response headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				log.Printf("[YouTube]   %s: %s", key, value)
			}
		}

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

		log.Printf("[YouTube] Successfully parsed %d channels", len(channelList.Items))

		items := arrutil.Map[YouTubeChannel, map[string]any](channelList.Items, func(input YouTubeChannel) (target map[string]any, find bool) {
			channelInfo := map[string]any{
				"id":   input.ID,
				"name": input.Snippet.Title,
			}

			log.Printf("[YouTube] Processing channel: ID=%s, Name=%s", input.ID, input.Snippet.Title)

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

		log.Printf("[YouTube] Returning %d channel items", len(items))
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
		log.Println("[YouTube] Starting getPlaylistID function")

		authCtx, err := ctx.AuthContext()
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to get auth context: %v", err)
			return nil, err
		}

		log.Printf("[YouTube] Auth context obtained successfully")
		log.Printf("[YouTube] Access token length: %d", len(authCtx.AccessToken))
		log.Printf("[YouTube] Access token preview: %s...", authCtx.AccessToken[:20])

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		log.Printf("[YouTube] Input channel ID: '%s'", input.ChannelID)

		// Create HTTP client
		client := &http.Client{}

		// Build URL with query parameters
		baseURL := "https://www.googleapis.com/youtube/v3/playlists?part=snippet,contentDetails&maxResults=50"

		// Check if we should filter by channel
		if input.ChannelID != "" {
			baseURL += "&channelId=" + input.ChannelID
			log.Printf("[YouTube] Filtering playlists by channel ID: %s", input.ChannelID)
		} else {
			// Default to user's own playlists
			baseURL += "&mine=true"
			log.Printf("[YouTube] Fetching user's own playlists (mine=true)")
		}

		log.Printf("[YouTube] Making GET request to: %s", baseURL)

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

		// Log request headers
		log.Println("[YouTube] Request headers:")
		for key, values := range req.Header {
			for _, value := range values {
				if key == "Authorization" {
					log.Printf("[YouTube]   %s: Bearer %s...", key, value[7:27]) // Log partial token
				} else {
					log.Printf("[YouTube]   %s: %s", key, value)
				}
			}
		}

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		log.Printf("[YouTube] Response status code: %d", resp.StatusCode)
		log.Printf("[YouTube] Response status: %s", resp.Status)

		// Log response headers
		log.Println("[YouTube] Response headers:")
		for key, values := range resp.Header {
			for _, value := range values {
				log.Printf("[YouTube]   %s: %s", key, value)
			}
		}

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

		log.Printf("[YouTube] Response body preview: %s", string(body[:min(len(body), 500)]))

		var playlistList YouTubePlaylistList
		err = json.Unmarshal(body, &playlistList)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
			log.Printf("[YouTube] Raw response that failed to unmarshal: %s", string(body))
			return nil, err
		}

		log.Printf("[YouTube] Successfully parsed %d playlists", len(playlistList.Items))

		items := arrutil.Map[YouTubePlaylist, map[string]any](playlistList.Items, func(input YouTubePlaylist) (target map[string]any, find bool) {
			playlistInfo := map[string]any{
				"id":   input.ID,
				"name": input.Snippet.Title,
			}

			log.Printf("[YouTube] Processing playlist: ID=%s, Name=%s", input.ID, input.Snippet.Title)

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

		log.Printf("[YouTube] Returning %d playlist items", len(items))
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

// Helper function for min since Go doesn't have a built-in for ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RegisterVideoProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getVideoList := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		log.Println("[YouTube] Starting getVideoList function")

		authCtx, err := ctx.AuthContext()
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to get auth context: %v", err)
			return nil, err
		}

		log.Printf("[YouTube] Auth context obtained successfully")

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		log.Printf("[YouTube] Input channel ID: '%s'", input.ChannelID)

		// Create HTTP client
		client := &http.Client{}

		var uploadsPlaylistID string

		// If channel ID is provided, use it; otherwise get the authenticated user's channel
		if input.ChannelID != "" {
			// Get the specified channel's uploads playlist
			channelURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=contentDetails&id=%s", input.ChannelID)
			log.Printf("[YouTube] Getting channel by ID: %s", channelURL)

			channelReq, err := http.NewRequest("GET", channelURL, nil)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to create channel request: %v", err)
				return nil, err
			}

			channelReq.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
			channelReq.Header.Set("Accept", "application/json")

			channelResp, err := client.Do(channelReq)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to get channel: %v", err)
				return nil, err
			}
			defer channelResp.Body.Close()

			if channelResp.StatusCode != http.StatusOK {
				log.Printf("[YouTube] ERROR: Channel request failed with status: %d", channelResp.StatusCode)
				return nil, fmt.Errorf("failed to get channel: %s", channelResp.Status)
			}

			channelBody, err := io.ReadAll(channelResp.Body)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to read channel response: %v", err)
				return nil, err
			}

			var channelData struct {
				Items []struct {
					ContentDetails struct {
						RelatedPlaylists struct {
							Uploads string `json:"uploads"`
						} `json:"relatedPlaylists"`
					} `json:"contentDetails"`
				} `json:"items"`
			}

			err = json.Unmarshal(channelBody, &channelData)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to unmarshal channel data: %v", err)
				return nil, err
			}

			if len(channelData.Items) == 0 {
				log.Printf("[YouTube] ERROR: Channel not found with ID: %s", input.ChannelID)
				return nil, fmt.Errorf("channel not found")
			}

			uploadsPlaylistID = channelData.Items[0].ContentDetails.RelatedPlaylists.Uploads
		} else {
			// Get the authenticated user's channel
			channelURL := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&mine=true"
			log.Printf("[YouTube] Getting user's own channel: %s", channelURL)

			channelReq, err := http.NewRequest("GET", channelURL, nil)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to create channel request: %v", err)
				return nil, err
			}

			channelReq.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
			channelReq.Header.Set("Accept", "application/json")

			channelResp, err := client.Do(channelReq)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to get channel: %v", err)
				return nil, err
			}
			defer channelResp.Body.Close()

			if channelResp.StatusCode != http.StatusOK {
				log.Printf("[YouTube] ERROR: Channel request failed with status: %d", channelResp.StatusCode)
				return nil, fmt.Errorf("failed to get channel: %s", channelResp.Status)
			}

			channelBody, err := io.ReadAll(channelResp.Body)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to read channel response: %v", err)
				return nil, err
			}

			var channelData struct {
				Items []struct {
					ContentDetails struct {
						RelatedPlaylists struct {
							Uploads string `json:"uploads"`
						} `json:"relatedPlaylists"`
					} `json:"contentDetails"`
				} `json:"items"`
			}

			err = json.Unmarshal(channelBody, &channelData)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to unmarshal channel data: %v", err)
				return nil, err
			}

			if len(channelData.Items) == 0 {
				log.Printf("[YouTube] ERROR: No channel found for authenticated user")
				return nil, errors.New("no channel found")
			}

			uploadsPlaylistID = channelData.Items[0].ContentDetails.RelatedPlaylists.Uploads
		}

		log.Printf("[YouTube] Found uploads playlist ID: %s", uploadsPlaylistID)

		// Now get the videos from the uploads playlist
		videosURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=%s&maxResults=50", uploadsPlaylistID)
		log.Printf("[YouTube] Getting videos from playlist: %s", videosURL)

		videosReq, err := http.NewRequest("GET", videosURL, nil)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to create videos request: %v", err)
			return nil, err
		}

		videosReq.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
		videosReq.Header.Set("Accept", "application/json")

		videosResp, err := client.Do(videosReq)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to get videos: %v", err)
			return nil, err
		}
		defer videosResp.Body.Close()

		if videosResp.StatusCode != http.StatusOK {
			log.Printf("[YouTube] ERROR: Videos request failed with status: %d", videosResp.StatusCode)
			return nil, fmt.Errorf("failed to get videos: %s", videosResp.Status)
		}

		videosBody, err := io.ReadAll(videosResp.Body)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to read videos response: %v", err)
			return nil, err
		}

		var videosData struct {
			Items []struct {
				Snippet struct {
					ResourceId struct {
						VideoId string `json:"videoId"`
					} `json:"resourceId"`
					Title       string `json:"title"`
					Description string `json:"description"`
					PublishedAt string `json:"publishedAt"`
				} `json:"snippet"`
			} `json:"items"`
		}

		err = json.Unmarshal(videosBody, &videosData)
		if err != nil {
			log.Printf("[YouTube] ERROR: Failed to unmarshal videos data: %v", err)
			return nil, err
		}

		log.Printf("[YouTube] Found %d videos", len(videosData.Items))

		items := arrutil.Map[struct {
			Snippet struct {
				ResourceId struct {
					VideoId string `json:"videoId"`
				} `json:"resourceId"`
				Title       string `json:"title"`
				Description string `json:"description"`
				PublishedAt string `json:"publishedAt"`
			} `json:"snippet"`
		}, map[string]any](videosData.Items, func(input struct {
			Snippet struct {
				ResourceId struct {
					VideoId string `json:"videoId"`
				} `json:"resourceId"`
				Title       string `json:"title"`
				Description string `json:"description"`
				PublishedAt string `json:"publishedAt"`
			} `json:"snippet"`
		}) (target map[string]any, find bool) {
			videoInfo := map[string]any{
				"id":   input.Snippet.ResourceId.VideoId,
				"name": input.Snippet.Title,
			}

			log.Printf("[YouTube] Processing video: ID=%s, Title=%s", input.Snippet.ResourceId.VideoId, input.Snippet.Title)

			// Add published date for additional context
			if input.Snippet.PublishedAt != "" {
				// Parse and format the date for better readability
				if t, err := time.Parse(time.RFC3339, input.Snippet.PublishedAt); err == nil {
					videoInfo["description"] = fmt.Sprintf("Published on %s", t.Format("Jan 2, 2006"))
				}
			}

			return videoInfo, true
		})

		log.Printf("[YouTube] Returning %d video items", len(items))
		return ctx.Respond(items, len(items))
	}

	helpText := "Select a video from the YouTube channel"

	return form.SelectField(fieldName, label).
		Placeholder("Select a video").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getVideoList)).
				WithFieldReference("channel_id", "channel_id").
				WithSearchSupport().
				WithPagination(20).
				End().
				RefreshOn("channel_id").
				GetDynamicSource(),
		).
		HelpText(helpText)
}
