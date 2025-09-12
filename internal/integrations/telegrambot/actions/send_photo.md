# Send Photo

## Description

Send a photo with optional caption to a Telegram chat.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name                  | Type    | Required | Description                                                           |
|-----------------------|---------|----------|-----------------------------------------------------------------------|
| Chat ID               | string  | Yes      | Unique identifier for the target chat or username of the target channel/group/user |
| Photo URL             | string  | Yes      | URL of the photo to send                                              |
| Caption               | string  | No       | Photo caption (may also be used when resending photos by file_id)      |
| Parse Mode            | select  | No       | Mode for parsing entities in the caption (None, Markdown, HTML)        |
| Disable Notification  | boolean | No       | Sends the message silently. Users will receive a notification with no sound |
| Reply to Message ID   | string  | No       | If the message is a reply, ID of the original message                 |