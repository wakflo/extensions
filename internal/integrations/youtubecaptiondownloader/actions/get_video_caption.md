# Get YouTube Caption

The Get YouTube Caption action allows you to extract captions/transcripts from any public YouTube video without requiring authentication or API keys.

**Important Notes:**

- Works with any public YouTube video - no ownership or permissions required
- No authentication needed - uses public YouTube data
- No API quota costs or rate limits beyond reasonable usage
- Cannot access private, unlisted, or age-restricted videos

## Features

### Video Input

- **YouTube URL**: Supports multiple URL formats
  - Standard: `https://www.youtube.com/watch?v=VIDEO_ID`
  - Short: `https://youtu.be/VIDEO_ID`
  - Embedded: `https://www.youtube.com/embed/VIDEO_ID`
  - Shorts: `https://youtube.com/shorts/VIDEO_ID`
  - Direct ID: Just the 11-character video ID

### Output Options

- **Both**: Returns structured captions with timestamps + full concatenated text
- **Structured**: Returns only the caption array with timing information
- **Text**: Returns only the full transcript as plain text

### Language Support

- **Language Selection**: Specify which language caption track to extract
- **Auto-detection**: Defaults to English if not specified
- **Fallback Logic**:
  - Tries exact language match first (e.g., "en")
  - Falls back to language variants (e.g., "en-US" for "en")
  - Falls back to English if requested language unavailable
- **Error Reporting**: Shows available languages if requested language not found

## Response Structure

The action returns:

- **video_id**: The extracted YouTube video ID
- **video_url**: The original input URL
- **video_title**: The title of the video
- **language**: The language code of the extracted captions
- **total_segments**: Number of caption segments
- **full_text**: Complete transcript as a single string (if format is "both" or "text")
- **captions**: Array of caption objects (if format is "both" or "structured")
  - **text**: The caption text
  - Additional timing fields (structure depends on library version)

## Supported Languages

The action supports any language that has captions available on the video. Common language codes include:

- English: `en`
- Spanish: `es`
- French: `fr`
- German: `de`
- Italian: `it`
- Portuguese: `pt`
- Russian: `ru`
- Japanese: `ja`
- Korean: `ko`
- Chinese (Simplified): `zh-CN`
- Chinese (Traditional): `zh-TW`
- Arabic: `ar`
- Hindi: `hi`
- Dutch: `nl`
- Polish: `pl`

## Common Use Cases

1. **Content Analysis**: Extract transcripts for sentiment analysis, keyword extraction, or content summarization
2. **Accessibility**: Create text versions of video content for hearing-impaired users
3. **Translation Preparation**: Get original captions for translation workflows
4. **Search & Discovery**: Index video content based on spoken words for better searchability
5. **Note Taking**: Automatically generate study notes from educational videos
6. **Content Repurposing**: Convert video content to blog posts, articles, or social media posts
7. **Research**: Analyze speech patterns, word frequency, or linguistic features
8. **SEO Optimization**: Extract video content for search engine optimization

## Example Usage

### Basic Usage

```json
{
	"videoURL": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"outputFormat": "both",
	"language": "en"
}
```

### Response Example

```json
{
	"video_id": "dQw4w9WgXcQ",
	"video_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	"video_title": "Rick Astley - Never Gonna Give You Up",
	"language": "en",
	"total_segments": 125,
	"full_text": "We're no strangers to love You know the rules and so do I...",
	"captions": [
		{
			"text": "We're no strangers to love",
			"start": 0,
			"duration": 2500
		},
		{
			"text": "You know the rules and so do I",
			"start": 2500,
			"duration": 2800
		}
	]
}
```

## Error Handling

Common errors and solutions:

### Invalid URL Format

- **Error**: "Failed to extract video ID: unsupported YouTube URL format or invalid video ID"
- **Solution**: Ensure you're using a supported URL format or valid video ID

### Video Not Found

- **Error**: "Failed to fetch video: ..."
- **Solution**: Check if the video exists and is publicly accessible

### No Captions Available

- **Error**: "No captions available for this video"
- **Solution**: The video doesn't have any caption tracks. Try a different video

### Language Not Available

- **Error**: "Captions not available in language 'xx'. Available languages: en, es, fr..."
- **Solution**: Use one of the listed available languages or leave empty for default

### Network Issues

- **Error**: "Failed to fetch transcript: ..."
- **Solution**: Check your internet connection and try again

## Technical Details

- Uses the `kkdai/youtube/v2` library for YouTube access
- No API keys or OAuth tokens required
- Processes captions client-side without external API calls
- Respects YouTube's terms of service by only accessing public data

## Limitations

1. **Public Videos Only**: Cannot access private, unlisted, or age-restricted videos
2. **Caption Availability**: Only works with videos that have captions enabled
3. **Auto-generated Captions**: Quality depends on YouTube's auto-caption accuracy
4. **No Caption Upload**: This action only downloads/extracts captions, cannot upload new ones
5. **Format Preservation**: Caption styling and advanced formatting may be lost in plain text output
