# Send Message

## Description

Send a text message to a WhatsApp number. The recipient must have previously opted in to receive messages from your business or messaged you in the last 24 hours.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Recipient's Phone Number | String | Yes | The recipient's phone number in international format (e.g., +1XXXXXXXXXX). |
| Message Text | String | Yes | The text message you want to send. |

## Details

- **Type**: sdkcore.ActionTypeNormal

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

- The WhatsApp Business API has restrictions on when you can send messages to users. You can send messages to users who have:
  - Explicitly opted in to receive messages from your business
  - Sent a message to your business within the last 24 hours (inside the "customer service window")
- Outside these conditions, you must use pre-approved message templates.
- The message ID in the response can be used to track message delivery status.