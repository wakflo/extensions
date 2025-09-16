# List YouTube Videos

This action allows you to search and retrieve videos from YouTube using various criteria.

## Search Methods

You can use one of the following methods to find videos:

1. **Search Query**: Search for videos using keywords
2. **Channel ID**: List videos from a specific channel
3. **Playlist ID**: Get videos from a specific playlist
4. **Video IDs**: Retrieve specific videos by their IDs

## Filters

When using search, you can apply the following filters:

- **Max Results**: Number of videos to return (1-50)
- **Sort Order**: How to order the results (relevance, date, rating, etc.)
- **Video Duration**: Filter by video length
- **Video Type**: Filter by episode or movie
- **Published Date Range**: Filter by publication date

## Output

The action returns:

- List of videos with detailed information
- Next page token (for pagination)
- Total number of results

Each video includes:

- Video ID and title
- Description
- Channel information
- Publication date
- Duration
- View, like, and comment counts
- Thumbnail URLs
