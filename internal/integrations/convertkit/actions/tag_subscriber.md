# Add Tag to Subscriber

## Description

Apply an existing tag to a specific subscriber to enhance segmentation and targeting capabilities.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Tag ID | Number | Yes | The ID of the tag to apply |
| Email | String | Yes* | The email address of the subscriber |


*Note: Email must be provided.

## Sample Response

```json
{
  "subscription": {
    "subscriber": {
      "id": 123456,
      "email_address": "jane@example.com",
      "first_name": "Jane",
      "state": "active"
    },
    "tag": {
      "id": 789,
      "name": "VIP Customer"
    }
  }
}
```

## Details

- **Type**: sdkcore.ActionTypeNormal