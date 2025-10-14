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
			"https://www.googleapis.com/auth/youtube.force-ssl",
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/youtube.upload",
			"https://www.googleapis.com/auth/youtubepartner",
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

func RegisterVideoProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getVideoID := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		input := sdk.DynamicInputToType[struct {
			ChannelID string `json:"channel_id"`
		}](ctx)

		client := &http.Client{}

		var allVideos []map[string]any

		if input.ChannelID == "" || input.ChannelID == "mine" {
			channelURL := "https://www.googleapis.com/youtube/v3/channels?part=contentDetails&mine=true"

			req, err := http.NewRequest("GET", channelURL, nil)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to create channel request: %v", err)
				return nil, err
			}

			req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
			req.Header.Set("Accept", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to get channel: %v", err)
				return nil, err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to read channel response: %v", err)
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("[YouTube] ERROR channel response: %s", string(body))
				return nil, fmt.Errorf("Channel API error %d: %s", resp.StatusCode, resp.Status)
			}

			var channelResp struct {
				Items []struct {
					ContentDetails struct {
						RelatedPlaylists struct {
							Uploads string `json:"uploads"`
						} `json:"relatedPlaylists"`
					} `json:"contentDetails"`
				} `json:"items"`
			}

			err = json.Unmarshal(body, &channelResp)
			if err != nil || len(channelResp.Items) == 0 {
				log.Printf("[YouTube] ERROR: Failed to get uploads playlist: %v", err)
				return nil, fmt.Errorf("could not get uploads playlist")
			}

			uploadsPlaylistID := channelResp.Items[0].ContentDetails.RelatedPlaylists.Uploads
			log.Printf("[YouTube] INFO: Found uploads playlist: %s", uploadsPlaylistID)

			var pageToken string
			totalFetched := 0

			for {
				playlistURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet,status&playlistId=%s&maxResults=50", uploadsPlaylistID)

				if pageToken != "" {
					playlistURL += "&pageToken=" + pageToken
				}

				req, err := http.NewRequest("GET", playlistURL, nil)
				if err != nil {
					log.Printf("[YouTube] ERROR: Failed to create playlist request: %v", err)
					return nil, err
				}

				req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
				req.Header.Set("Accept", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					log.Printf("[YouTube] ERROR: Failed to get playlist items: %v", err)
					return nil, err
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("[YouTube] ERROR: Failed to read playlist response: %v", err)
					return nil, err
				}

				if resp.StatusCode != http.StatusOK {
					log.Printf("[YouTube] ERROR playlist response: %s", string(body))
					return nil, fmt.Errorf("Playlist API error %d: %s", resp.StatusCode, resp.Status)
				}

				var playlistResp struct {
					NextPageToken string `json:"nextPageToken"`
					PageInfo      struct {
						TotalResults   int `json:"totalResults"`
						ResultsPerPage int `json:"resultsPerPage"`
					} `json:"pageInfo"`
					Items []struct {
						Snippet struct {
							Title        string `json:"title"`
							Description  string `json:"description"`
							PublishedAt  string `json:"publishedAt"`
							ChannelTitle string `json:"channelTitle"`
							ResourceId   struct {
								VideoId string `json:"videoId"`
							} `json:"resourceId"`
							Thumbnails struct {
								Default struct {
									Url string `json:"url"`
								} `json:"default"`
							} `json:"thumbnails"`
						} `json:"snippet"`
						Status struct {
							PrivacyStatus string `json:"privacyStatus"`
						} `json:"status"`
					} `json:"items"`
				}

				err = json.Unmarshal(body, &playlistResp)
				if err != nil {
					log.Printf("[YouTube] ERROR: Failed to unmarshal playlist response: %v", err)
					return nil, err
				}

				log.Printf("[YouTube] INFO: Fetched %d videos in this page, total results: %d",
					len(playlistResp.Items), playlistResp.PageInfo.TotalResults)

				for _, item := range playlistResp.Items {
					privacyEmoji := ""
					privacyLabel := ""
					switch item.Status.PrivacyStatus {
					case "private":
						privacyEmoji = "🔒"
						privacyLabel = "Private"
					case "unlisted":
						privacyEmoji = "🔗"
						privacyLabel = "Unlisted"
					case "public":
						privacyEmoji = "🌐"
						privacyLabel = "Public"
					}

					videoInfo := map[string]any{
						"id":           item.Snippet.ResourceId.VideoId,
						"name":         fmt.Sprintf("%s %s", privacyEmoji, item.Snippet.Title),
						"publishedAt":  item.Snippet.PublishedAt,
						"channel":      item.Snippet.ChannelTitle,
						"privacy":      item.Status.PrivacyStatus,
						"privacyLabel": privacyLabel,
						"description":  item.Snippet.Description,
						"thumbnail":    item.Snippet.Thumbnails.Default.Url,
					}

					allVideos = append(allVideos, videoInfo)
					totalFetched++
				}

				// Check if there are more pages
				pageToken = playlistResp.NextPageToken
				if pageToken == "" {
					log.Printf("[YouTube] INFO: Finished fetching all videos. Total: %d", totalFetched)
					break
				}

				log.Printf("[YouTube] INFO: Moving to next page, fetched so far: %d", totalFetched)
			}

			log.Printf("[YouTube] INFO: Returning %d total videos", len(allVideos))

		} else {
			log.Printf("[YouTube] INFO: Fetching videos for channel: %s", input.ChannelID)

			baseURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&maxResults=50&order=date&channelId=%s", input.ChannelID)

			req, err := http.NewRequest("GET", baseURL, nil)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to create request: %v", err)
				return nil, err
			}

			req.Header.Set("Authorization", "Bearer "+authCtx.AccessToken)
			req.Header.Set("Accept", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to send request: %v", err)
				return nil, err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to read response body: %v", err)
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("[YouTube] ERROR response body: %s", string(body))
				return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, resp.Status)
			}

			var searchResult struct {
				Items []struct {
					ID struct {
						VideoID string `json:"videoId"`
					} `json:"id"`
					Snippet struct {
						Title        string `json:"title"`
						PublishedAt  string `json:"publishedAt"`
						ChannelTitle string `json:"channelTitle"`
					} `json:"snippet"`
				} `json:"items"`
			}

			err = json.Unmarshal(body, &searchResult)
			if err != nil {
				log.Printf("[YouTube] ERROR: Failed to unmarshal response: %v", err)
				return nil, err
			}

			// For other channels, we can only see public videos
			for _, item := range searchResult.Items {
				videoInfo := map[string]any{
					"id":           item.ID.VideoID,
					"name":         fmt.Sprintf("🌐 %s", item.Snippet.Title),
					"publishedAt":  item.Snippet.PublishedAt,
					"channel":      item.Snippet.ChannelTitle,
					"privacy":      "public",
					"privacyLabel": "Public",
				}
				allVideos = append(allVideos, videoInfo)
			}
		}

		return ctx.Respond(allVideos, len(allVideos))
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
		HelpText("Select a YouTube video (🔒 Private, 🔗 Unlisted, 🌐 Public)")
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

func RegisterLanguageProps(form *smartform.FormBuilder, fieldName, label string, required bool) *smartform.FieldBuilder {
	getLanguages := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		// Create HTTP client
		client := &http.Client{}

		// Use i18nLanguages API to get list of supported languages
		url := "https://www.googleapis.com/youtube/v3/i18nLanguages?part=snippet&hl=en"

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
			return nil, err
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			// Try to parse Google API error
			var apiError struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			json.Unmarshal(body, &apiError)
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, apiError.Error.Message)
		}

		// Parse the language list response
		var langList struct {
			Items []struct {
				ID      string `json:"id"`
				Snippet struct {
					HL   string `json:"hl"`
					Name string `json:"name"`
				} `json:"snippet"`
			} `json:"items"`
		}

		err = json.Unmarshal(body, &langList)
		if err != nil {
			return nil, err
		}

		// Map to options format
		items := arrutil.Map[struct {
			ID      string `json:"id"`
			Snippet struct {
				HL   string `json:"hl"`
				Name string `json:"name"`
			} `json:"snippet"`
		}, map[string]any](langList.Items, func(input struct {
			ID      string `json:"id"`
			Snippet struct {
				HL   string `json:"hl"`
				Name string `json:"name"`
			} `json:"snippet"`
		}) (target map[string]any, find bool) {
			return map[string]any{
				"id":   input.ID,
				"name": fmt.Sprintf("%s (%s)", input.Snippet.Name, input.ID),
			}, true
		})

		// Add option for original language at the beginning
		originalOption := map[string]any{
			"id":   "",
			"name": "Original Language (No Translation)",
		}
		items = append([]map[string]any{originalOption}, items...)

		return ctx.Respond(items, len(items))
	}

	return form.SelectField(fieldName, label).
		Placeholder("Select a language for translation").
		Required(required).
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getLanguages)).
				WithSearchSupport().
				End().
				GetDynamicSource(),
		).
		HelpText("Leave as 'Original Language' to download without translation, or select a language to translate the captions").
		DefaultValue("")
}
