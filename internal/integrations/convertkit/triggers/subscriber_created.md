# Subscriber Created

## Description

Triggers a workflow when a new subscriber is added to your ConvertKit account, allowing you to automate follow-up actions or sync the data with other systems.

## Properties

This trigger doesn't require any additional properties.

## Sample Response

```json
{
  "subscribers": [
    {
      "id": 123456,
      "first_name": "Jane",
      "email_address": "jane@example.com",
      "state": "active",
      "created_at": "2023-01-15T10:30:00Z",
      "fields": {
        "company": "Acme Inc"
      },
      "tags": [123, 456]
    },
    {
      "id": 789012,
      "first_name": "John",
      "email_address": "john@example.com",
      "state": "active",
      "created_at": "2023-01-16T14:20:00Z",
      "fields": {
        "company": "XYZ Corp"
      },
      "tags": [123]
    }
  ],
  "total_subscribers": 2
}
```

## Details

- **Type**: sdkcore.TriggerTypePolling
- **Polling Interval**: 5 minutes

## Usage Notes

This trigger polls the ConvertKit API every 5 minutes to check for new subscribers that have been added since the last run. The first time the trigger runs, it will fetch subscribers added within the last 24 hours.