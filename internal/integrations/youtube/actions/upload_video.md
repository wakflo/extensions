# Upload YouTube Video

This action uploads a new video to YouTube with comprehensive metadata and settings.

## Channel Selection

- **Channel**: Select the channel where you want to upload the video

## Video File

- **Video File**: Select the video file to upload
- **Supported Formats**: MP4, AVI, MOV, WMV, FLV, 3GPP, WebM

## Video Information

### Required Fields

- **Title**: Video title (max 100 characters)

### Optional Fields

- **Description**: Detailed video description (max 5000 characters)
- **Tags**: Comma-separated keywords to help with discovery
- **Category**: YouTube category for proper classification

## Privacy and Publishing

- **Privacy Status**:
  - Private: Only you can view
  - Unlisted: Anyone with link can view
  - Public: Everyone can find and view
- **Notify Subscribers**: Send notifications (public videos only)
- **Allow Embedding**: Permit embedding on external sites
- **Public Stats Viewable**: Show view count publicly

## Content Settings

- **Made for Kids**: COPPA compliance requirement (required field)
- **Self-Declared Made for Kids**: Additional declaration
- **License**: Standard YouTube or Creative Commons

## Advanced Settings

- **Recording Date**: When the video was recorded
- **Default Language**: Primary language (e.g., 'en')
- **Default Audio Language**: Audio track language

## Processing Options

- **Auto-levels**: Automatic lighting/color correction
- **Stabilize**: Reduce camera shake

## Upload Process

1. Video uploads in 1MB chunks for reliability
2. Processing begins after upload completes
3. Video will be in "processing" status initially
4. Final availability depends on video length and quality

## Output

Returns information about the uploaded video:

- Video ID for tracking
- Upload and processing status
- All applied settings
- Channel information

## Important Notes

- Large files may take significant time to upload
- Processing time varies based on video length/quality
- Privacy can be changed after upload
- Some features unavailable for "Made for Kids" videos
