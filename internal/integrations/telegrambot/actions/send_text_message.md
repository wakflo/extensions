# Send Message

## Description

Send a text message to a specific Telegram chat, user, group, or channel.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name                    | Type    | Required | Description                                                           |
|-------------------------|---------|----------|-----------------------------------------------------------------------|
| Chat ID                 | string  | Yes      | Unique identifier for the target chat of the target channel/group/user |
| Message Text            | string  | Yes      | Text of the message to be sent                                        |
| Parse Mode              | select  | No       | Mode for parsing entities in the message text (None, Markdown, HTML)  |
| Disable Web Page Preview | boolean | No       | Disables link previews for links in this message                      |
| Disable Notification    | boolean | No       | Sends the message silently. Users will receive a notification with no sound |
| Reply to Message ID     | string  | No       | If the message is a reply, ID of the original message                 |