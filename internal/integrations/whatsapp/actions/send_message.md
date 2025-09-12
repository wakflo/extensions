# Send Media

## Description
Send media files (images, videos, documents, audio, or stickers) to a WhatsApp number. The recipient must have previously opted in to receive messages from your business or messaged you in the last 24 hours.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Recipient's Phone Number | String | Yes | The recipient's phone number in international format (e.g., +1XXXXXXXXXX). |
| Media Type | Select | Yes | The type of media you want to send. Options: Image, Video, Audio, Document, Sticker. |
| Media URL | String | No* | Public URL of the media file. Either Media URL or Media ID is required. |
| Caption | String | No | Caption for the media. Supported for image, video, and document types only. |
| Filename | String | No** | Filename for document type media. Required when sending documents. |

*One of these fields is required  
**Required when Media Type is "Document"

## Details
* **Type**: core.ActionTypeAction

## Sample Response

```json
{
  "messaging_product": "whatsapp",
  "contacts": [
    {
      "input": "+1234567890",
      "wa_id": "1234567890"
    }
  ],
  "messages": [
    {
      "id": "wamid.HBgMMTIzNDU2Nzg5MBUCABIYFjNFQjBEN0EwNTAwNEUwRkI0MzU4QTkA"
    }
  ]
}
```

## Notes

### General Restrictions
* The WhatsApp Business API has restrictions on when you can send messages to users. You can send messages to users who have:
  * Explicitly opted in to receive messages from your business
  * Sent a message to your business within the last 24 hours (inside the "customer service window")
* Outside these conditions, you must use pre-approved message templates.
* The message ID in the response can be used to track message delivery status.

### Media Requirements
* **Media URL**: Must be a publicly accessible HTTPS URL. WhatsApp will download the media from this URL.
* **Media ID**: Use this when you've previously uploaded media to WhatsApp's servers using the Media Upload API.
* **File Size Limits**:
  * Images: 5MB (JPEG, PNG)
  * Videos: 16MB (MP4, 3GPP)
  * Audio: 16MB (AAC, AMR, MP3, MP4 Audio, OGG)
  * Documents: 100MB (PDF, DOC, DOCX, PPT, PPTX, XLS, XLSX)
  * Stickers: 100KB (WEBP - static) or 500KB (WEBP - animated)

### Supported Media Formats
* **Images**: JPEG, PNG
* **Videos**: MP4, 3GPP
* **Audio**: AAC, AMR, MP3, M4A, OGG (Opus codec only)
* **Documents**: PDF, DOC, DOCX, PPT, PPTX, XLS, XLSX
* **Stickers**: WEBP (must be 512x512 pixels)

### Caption Support
* Captions are supported for: Images, Videos, and Documents
* Captions are NOT supported for: Audio and Stickers
* Maximum caption length: 1024 characters

### Document Filename
* When sending documents, you should provide a filename to display to the recipient
* If not provided, WhatsApp will use a default filename
* The filename should include the appropriate file extension