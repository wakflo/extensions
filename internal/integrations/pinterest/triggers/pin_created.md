# Pin Created Trigger

This trigger fires when a new pin is created in your Pinterest account.

## Configuration

### Board Filter (Optional)

You can optionally filter pins by board. If no board is specified, the trigger will monitor pins from all boards.

### Page Size

Configure how many pins to check per polling interval (1-100). Default is 25.

## Pinterest API Notes

This trigger uses the Pinterest API v5 endpoints:

- **List Pins**: GET /v5/pins - Returns pins owned by the authenticated user
- **List Board Pins**: GET /v5/boards/{board_id}/pins - Returns pins from a specific board

The trigger polls for new pins by comparing the created_at timestamp with the last run time.

## Rate Limits

Please be aware of Pinterest API rate limits:

- Standard rate limit: 1000 requests per hour per user
- Some endpoints may have specific limits

## Required Permissions

Your Pinterest app needs the following permissions:

- pins:read - Read access to pins

## Sample Output

The trigger returns an array of pin objects with the following structure:

\`\`\`json
{
"id": "123456789012345678",
"created_at": "2024-08-20T14:30:00Z",
"link": "https://www.pinterest.com/pin/123456789012345678",
"title": "Beautiful Sunset Photography",
"description": "Amazing sunset captured at the beach during golden hour",
"alt_text": "Sunset at beach",
"note": "Private note about this pin",
"board_id": "987654321098765432",
"board_section_id": "876543210987654321",
"media": {
"media_type": "image",
"images": {
"150x150": {
"url": "https://i.pinimg.com/150x150/example.jpg",
"width": 150,
"height": 150
},
"400x300": {
"url": "https://i.pinimg.com/400x300/example.jpg",
"width": 400,
"height": 300
},
"600x": {
"url": "https://i.pinimg.com/600x/example.jpg",
"width": 600,
"height": 450
}
}
},
"parent_pin_id": null,
"is_standard": true,
"has_been_promoted": false,
"creative_type": "REGULAR",
"product_tags": []
}
\`\`\`

## Limitations

- The trigger only detects pins created by the authenticated user
- Pinterest API v5 does not provide webhooks, so this trigger uses polling
- Only pins created after the trigger is activated will be detected
