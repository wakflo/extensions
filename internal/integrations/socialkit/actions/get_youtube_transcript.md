# Get YouTube Transcript

## Description

Extract accurate, timestamped transcripts from YouTube videos for content analysis, accessibility, and data extraction. This action retrieves both the complete transcript and individual timestamped segments from any YouTube video.

## Properties

| Name      | Type    | Required | Description                                                                       |
| --------- | ------- | -------- | --------------------------------------------------------------------------------- |
| video_url | string  | yes      | The full URL of the YouTube video (e.g., https://youtube.com/watch?v=dQw4w9WgXcQ) |
| cache     | boolean | no       | Enable caching for faster subsequent requests (default: false)                    |
| cache_ttl | number  | no       | Cache time-to-live in seconds, between 3600 (1 hour) and 2592000 (1 month)        |

## Response

The action returns a JSON object containing:

- **transcript**: The complete transcript text of the video
- **segments**: An array of timestamped transcript segments, each containing:
  - `text`: The text content of the segment
  - `start`: Start time in seconds
  - `end`: End time in seconds
- **metadata**: Video information including:
  - `video_id`: YouTube video ID
  - `title`: Video title
  - `duration`: Video duration in seconds
  - `language`: Transcript language code
  - `channel`: Channel name (if available)

## Details

- **Type**: sdkcore.ActionTypeAction
- **Platforms**: YouTube
- **Rate Limits**: Subject to your SocialKit plan limits

## Use Cases

1. **Content Analysis**: Extract transcripts for keyword analysis, sentiment analysis, or topic modeling
2. **Accessibility**: Generate captions or subtitles for videos
3. **Content Repurposing**: Convert video content into blog posts, articles, or social media content
4. **SEO**: Analyze video transcripts for search engine optimization
5. **Research**: Collect and analyze video content for academic or market research

## Notes

- The video must have available captions or transcripts on YouTube
- Cached responses are faster but may not reflect recent changes to the video
- Cache TTL must be between 3600 seconds (1 hour) and 2592000 seconds (1 month)
