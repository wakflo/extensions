# Update YouTube Video

This action allows you to update various metadata fields of a YouTube video.

## Video Selection

- **Channel**: Select the channel containing the video
- **Video to Update**: Select the video you want to update (list updates based on selected channel)

## Updatable Fields

### Basic Information

- **Title**: Video title (max 100 characters)
- **Description**: Video description (max 5000 characters)
- **Tags**: Comma-separated list of keywords
- **Category**: YouTube category for the video

### Privacy and Status Settings

- **Privacy Status**: Private, Unlisted, or Public
- **Allow Embedding**: Whether the video can be embedded on other sites
- **Public Stats Viewable**: Whether video statistics are publicly visible

### Content Settings

- **Made for Kids**: Designate if video is child-directed
- **Self-Declared Made for Kids**: Channel owner's designation
- **License**: Standard YouTube or Creative Commons

### Additional Settings

- **Recording Date**: When the video was recorded
- **Default Language**: Video's primary language
- **Default Audio Language**: Audio track language

## Important Notes

- Only fill in fields you want to update
- Empty fields will keep their current values
- Changes to privacy status take effect immediately
- Some fields cannot be changed after upload (e.g., video file)

## Output

Returns the updated video information including:

- Video ID and title
- Updated description and tags
- Current status settings
- Channel information
- Language settings
