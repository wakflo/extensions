The Download YouTube Caption action allows you to download caption/subtitle tracks from YouTube videos.

**Important Notes:**

- You can only download captions from videos that you own or have edit permissions for
- The YouTube API requires proper OAuth authentication with youtube.force-ssl scope
- Caption downloads have a quota cost of 200 units

## Features

### Video Selection

- **Own Channel Videos**: Select from your own channel's videos with available captions
- **Manual Video ID**: Enter a specific video ID (must be a video you own)

### Caption Options

- **Caption Track**: Select a specific caption track or let the action choose the first available
- **Format**: Download in various formats (SRT, VTT, SBV, TTML, SCC) or keep the original format
- **Translation**: Automatically translate captions to any language using Google Translate

## Supported Formats

- **SRT** (SubRip Subtitle): Most common subtitle format, widely supported
- **VTT** (Web Video Text Tracks): W3C standard for HTML5 video
- **SBV** (SubViewer Subtitle): YouTube's legacy subtitle format
- **TTML** (Timed Text Markup Language): XML-based caption format
- **SCC** (Scenarist Closed Caption): Professional broadcast format

## Language Translation

The action supports automatic translation to any ISO 639-1 language code:

- English: en
- Spanish: es
- French: fr
- German: de
- Italian: it
- Portuguese: pt
- Russian: ru
- Japanese: ja
- Korean: ko
- Chinese (Simplified): zh-CN

## Response

The action returns:

- **captionContent**: The full caption text in the requested format
- **videoId**: The ID of the video
- **captionId**: The ID of the caption track
- **format**: The format of the downloaded caption
- **originalLanguage**: The original language of the caption
- **isTranslated**: Whether the caption was translated
- **contentPreview**: First 500 characters of the caption content

## Common Use Cases

1. **Download and Convert Formats**: Download captions in a different format for use in other video editing software
2. **Translation**: Automatically translate captions to reach a global audience
3. **Backup**: Create backups of your video captions
4. **Processing**: Download captions for further processing or analysis

## Error Handling

Common errors:

- **403 Forbidden**: You don't have permission to download captions for this video
- **404 Not Found**: The caption track doesn't exist
- **400 Bad Request**: Invalid format or language code specified
