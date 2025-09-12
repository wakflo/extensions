# Message Received

## Description

Triggered when your Telegram bot receives a new message from a user.

This trigger polls the Telegram API for new messages and activates when new messages are found. It tracks the last processed update ID to ensure each message is only processed once.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Output

This trigger returns an object with the following properties:

| Name         | Type    | Description                                                      |
|--------------|---------|------------------------------------------------------------------|
| messages     | array   | Array of message update objects received from Telegram           |
| lastUpdateID | integer | The ID of the last processed update (used for future polling)    |

Each message in the messages array includes:

| Name        | Type    | Description                                                      |
|-------------|---------|------------------------------------------------------------------|
| update_id   | integer | The update's unique identifier                                   |
| message     | object  | The message object with details about the message                |

The message object contains:

| Name       | Type    | Description                                                       |
|------------|---------|-------------------------------------------------------------------|
| message_id | integer | Unique message identifier inside this chat                        |
| from       | object  | Sender information (id, first_name, username, etc.)               |
| chat       | object  | Information about the chat (id, type, title, etc.)                |
| date       | integer | Date the message was sent in Unix time                            |
| text       | string  | For text messages, the actual UTF-8 text of the message           |
| photo      | array   | Optional. Array of PhotoSize objects for image messages           |
| document   | object  | Optional. Information about a document/file                       |
| audio      | object  | Optional. Information about an audio file                         |
| video      | object  | Optional. Information about a video                               |