# YouTube Caption Downloader Integration Documentation

The YouTube Caption Downloader integration provides tools for extracting captions and transcripts from YouTube videos within your workflows. This integration is especially useful for content analysis, accessibility, and transforming video content into text format without requiring any authentication or API keys.

## Overview

YouTube is the world's largest video platform with billions of hours of content. The YouTube Caption Downloader integration offers a simple way to access the textual content of videos through their caption tracks, enabling:

- Extracting captions from any public YouTube video
- Converting video content to searchable text
- Analyzing spoken content programmatically
- Creating accessible versions of video content

## Available Actions

### Get Video Caption

Extracts captions/transcripts from YouTube videos in various formats. This action can handle multiple YouTube URL formats and provides options for structured or plain text output. See the [full documentation](get_video_caption.md) for details.

## Requirements

To use the YouTube Caption Downloader integration, you need:

- A Wakflo account with integration capabilities enabled
- YouTube video URLs or video IDs
- Videos must be public and have captions available
- No API keys or authentication required

## Common Use Cases

### Content Analysis & Research

Extract video transcripts for various analytical purposes:

1. Use Get Video Caption to extract transcripts from educational videos
2. Process the text data for keyword analysis or topic modeling
3. Generate summaries or key points from long-form content
4. Analyze speech patterns or linguistic features

### Accessibility & Compliance

Make video content accessible to all users:

1. Extract captions from videos for hearing-impaired users
2. Create text-based alternatives to video content
3. Generate searchable transcripts for video libraries
4. Ensure compliance with accessibility standards

### Content Repurposing

Transform video content into other formats:

1. Convert video tutorials into written documentation
2. Create blog posts from video content
3. Generate social media snippets from longer videos
4. Build searchable knowledge bases from video archives

### SEO & Discovery

Improve content discoverability:

1. Extract video transcripts for search engine optimization
2. Create searchable indexes of video content
3. Generate metadata and keywords from video captions
4. Build content recommendation systems based on transcript analysis

## Implementation Example

Here's a simple example of how to use the YouTube Caption Downloader in a workflow:

```go
// Import the YouTube Caption Downloader actions
import (
    ycd "github.com/wakflo/integrations/youtube-caption-downloader"
)

// In your setup code
func RegisterActions(connector *sdk.Connector) {
    // Register the Get Video Caption action
    connector.AddAction(ycd.NewGetVideoCaptionAction())
}

// Example workflow extracting captions
func ExtractVideoContent(videoURL string, language string) (map[string]interface{}, error) {
    // Create the action
    captionAction := ycd.NewGetVideoCaptionAction()

    // Prepare inputs
    inputs := map[string]interface{}{
        "videoURL": videoURL,
        "outputFormat": "both",
        "language": language,
    }

    // Execute the action
    result, err := captionAction.Perform(ctx)
    if err != nil {
        return nil, err
    }

    // Extract results
    videoData := map[string]interface{}{
        "title": result["video_title"],
        "transcript": result["full_text"],
        "captions": result["captions"],
        "language": result["language"],
    }

    return videoData, nil
}

// Example: Process multiple videos
func ProcessVideoPlaylist(videoURLs []string) ([]map[string]interface{}, error) {
    var results []map[string]interface{}

    for _, url := range videoURLs {
        // Extract captions for each video
        videoData, err := ExtractVideoContent(url, "en")
        if err != nil {
            // Log error but continue with other videos
            fmt.Printf("Error processing %s: %v\n", url, err)
            continue
        }

        // Process the transcript (e.g., sentiment analysis, keyword extraction)
        processedData := processTranscript(videoData["transcript"].(string))
        videoData["analysis"] = processedData

        results = append(results, videoData)
    }

    return results, nil
}
```

## Best Practices

1. **URL Validation**: Always validate YouTube URLs before processing to avoid unnecessary API calls.

2. **Language Selection**:

   - Specify the language code when you know it for faster processing
   - Use the default (English) fallback for general use cases
   - Check available languages in error messages when captions aren't found

3. **Error Handling**:

   - Implement retry logic for network failures
   - Handle cases where captions aren't available gracefully
   - Log errors for debugging but don't let single failures stop batch processing

4. **Output Format Selection**:

   - Use "text" format for simple text processing or display
   - Use "structured" format when you need timing information
   - Use "both" format when building comprehensive databases

5. **Rate Limiting**: While there's no strict API limit, be respectful:

   - Implement reasonable delays between requests
   - Cache results to avoid repeated extractions
   - Consider batch processing during off-peak hours

6. **Data Processing**:
   - Clean and normalize extracted text before analysis
   - Handle special characters and formatting appropriately
   - Consider language-specific processing for non-English content

## Troubleshooting

Common issues and their solutions:

### Invalid URL Errors

**Problem**: "Failed to extract video ID: unsupported YouTube URL format"

- **Solution**: Ensure you're using a supported URL format (standard, short, embedded, or direct video ID)
- **Check**: Remove any additional parameters or anchors from the URL

### Caption Not Available

**Problem**: "No captions available for this video"

- **Solution**: Verify the video has captions by checking on YouTube directly
- **Alternative**: Try auto-generated captions if manual captions aren't available

### Language Issues

**Problem**: "Captions not available in language 'xx'"

- **Solution**: Use one of the available languages listed in the error message
- **Tip**: Leave language empty to use the default available caption track

### Video Access Issues

**Problem**: "Failed to fetch video"

- **Solution**: Ensure the video is public and not age-restricted
- **Check**: Verify you can access the video in a browser without logging in

### Network Errors

**Problem**: Connection or timeout errors

- **Solution**: Check internet connectivity and retry
- **Implement**: Exponential backoff for retry attempts

## Limitations

1. **Public Videos Only**: Cannot access private, unlisted, or age-restricted videos
2. **Caption Dependency**: Only works with videos that have captions enabled
3. **No Upload Capability**: This integration only extracts captions, cannot create or upload new ones
4. **Format Limitations**: Some advanced caption formatting (colors, positioning) is not preserved
5. **Auto-caption Quality**: Accuracy depends on YouTube's automatic caption generation quality

## Support

For additional assistance with the YouTube Caption Downloader integration:

- Check the detailed action documentation:
  - [Get Video Caption action documentation](./get_video_caption.md)
- Contact Wakflo support at support@wakflo.com
- Visit the [Wakflo documentation](https://docs.wakflo.com) for more information
- Report issues on [GitHub](https://github.com/wakflo/integrations/issues)
