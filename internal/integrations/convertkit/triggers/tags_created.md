# Tag Added

## Description

Triggers a workflow when a tag is added to a subscriber in your ConvertKit account, enabling you to create targeted automations based on subscriber segmentation and behavior.

## Properties
This trigger doesn't require any additional properties.

## Sample Response

```json
{
  "tag_id": 789,
  "subscriptions": [
    {
      "id": 12345,
      "subscriber": {
        "id": 123456,
        "first_name": "Jane",
        "email_address": "jane@example.com"
      },
      "created_at": "2023-01-15T10:30:00Z"
    },
    {
      "id": 67890,
      "subscriber": {
        "id": 789012,
        "first_name": "John",
        "email_address": "john@example.com"
      },
      "created_at": "2023-01-16T14:20:00Z"
    }
  ],
  "total_new_subscriptions": 2
}
```

## Details

- **Type**: sdkcore.TriggerTypePolling
- **Polling Interval**: 5 minutes

## Usage Notes

This trigger polls the ConvertKit API every 5 minutes to check for new subscribers that have been tagged with the specified tag since the last run. The first time the trigger runs, it will fetch subscribers tagged within the last 24 hours.