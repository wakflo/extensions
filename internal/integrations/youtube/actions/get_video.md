# Get YouTube Video

This action retrieves detailed information about a specific YouTube video. You can either search for any public video by ID or select from your own channel's videos.

## Search Options

- **Search Only My Channel**: Toggle this option to switch between:
  - **Unchecked**: Search any public YouTube video by ID
  - **Checked**: Select from your own channel's uploaded videos

## Input Parameters

### When searching public videos:

- **Video ID**: The unique identifier of the YouTube video (required)
  - Example: `dQw4w9WgXcQ` from `youtube.com/watch?v=dQw4w9WgXcQ`

### When searching your channel:

- **Select Video**: Choose from a dropdown list of your uploaded videos

### Common parameters:

- **Data Parts**: Choose how much data to retrieve (basic, detailed, full)

## Data Parts Options

You can select different levels of detail:

- **Basic**: Essential information only
- **Detailed**: Standard information (default)
- **Full**: All available data

## Output

The action returns comprehensive information about the video:

### Video Information

- Video ID
- Title and description
- Channel information
- Publication date
- Tags and category

### Media Details

- Duration
- Video definition (HD/SD)
- Caption availability
- Thumbnail URLs (multiple resolutions)
- Embed HTML for player

### Statistics

- View count
- Like count
- Comment count
- Favorite count

### Status Information

- Privacy status (public, private, unlisted)
- Upload status
- License type
- Embeddable status
- Made for kids flag

### Additional Details

- Content rating and age restrictions
- Region restrictions (allowed/blocked countries)
- Topic categories
- Recording location and date (if available)
- Localizations for different languages
