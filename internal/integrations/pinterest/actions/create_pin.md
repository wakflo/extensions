# Create Pin

## Description

Creates a new pin on Pinterest with an image or video, allowing you to specify title, description, link, and other metadata.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name              | Type   | Required | Description                                                       |
| ----------------- | ------ | -------- | ----------------------------------------------------------------- |
| board_id          | String | Yes      | The unique identifier of the board where the pin will be created. |
| title             | String | No       | Title for the pin (max 100 characters).                           |
| description       | String | No       | Description for the pin (max 800 characters).                     |
| link              | String | No       | Link URL for the pin (must start with http:// or https://).       |
| alt_text          | String | No       | Alternative text for accessibility (max 500 characters).          |
| note              | String | No       | Note to add to the pin (max 500 characters).                      |
| media_source_type | String | Yes      | Type of media source: image_url, video_id, or image_base64.       |
| media_source_url  | String | No\*     | URL of the image (required when media_source_type is image_url).  |
| media_source_id   | String | No\*     | ID of the video (required when media_source_type is video_id).    |
| dominant_color    | String | No       | Dominant color in hex format (e.g., #FF5733).                     |

## Media Source Types

- **image_url**: Provide a URL to an image hosted online
- **video_id**: Use an existing video ID from Pinterest
- **image_base64**: Upload an image as base64 (not yet implemented)

## Notes

- You must have write access to the specified board
- Either media_source_url or media_source_id is required depending on the media_source_type
- The image URL must be publicly accessible
- Pinterest will download and store the image from the provided URL
- Video pins require a pre-uploaded video ID
- All text fields have character limits as specified in the properties table
