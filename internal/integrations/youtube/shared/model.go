package shared

// YouTube API response structures
type YouTubeChannel struct {
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CustomURL   string `json:"customUrl"`
		Thumbnails  struct {
			Default struct {
				URL string `json:"url"`
			} `json:"default"`
		} `json:"thumbnails"`
	} `json:"snippet"`
	Statistics struct {
		ViewCount       string `json:"viewCount"`
		SubscriberCount string `json:"subscriberCount"`
		VideoCount      string `json:"videoCount"`
	} `json:"statistics"`
}

type YouTubeChannelList struct {
	Items         []YouTubeChannel `json:"items"`
	NextPageToken string           `json:"nextPageToken"`
}

type YouTubePlaylist struct {
	ID      string `json:"id"`
	Snippet struct {
		ChannelID    string `json:"channelId"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ChannelTitle string `json:"channelTitle"`
		Thumbnails   struct {
			Default struct {
				URL string `json:"url"`
			} `json:"default"`
		} `json:"thumbnails"`
	} `json:"snippet"`
	ContentDetails struct {
		ItemCount int `json:"itemCount"`
	} `json:"contentDetails"`
}

type YouTubePlaylistList struct {
	Items         []YouTubePlaylist `json:"items"`
	NextPageToken string            `json:"nextPageToken"`
}

type YouTubeVideo struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	ID   struct {
		Kind    string `json:"kind"`
		VideoID string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		PublishedAt  string `json:"publishedAt"`
		ChannelID    string `json:"channelId"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ChannelTitle string `json:"channelTitle"`
		Thumbnails   struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
	} `json:"snippet"`
}

type YouTubeVideoList struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []YouTubeVideo `json:"items"`
}
