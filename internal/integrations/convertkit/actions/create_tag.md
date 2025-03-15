# Create Tag

## Description

Create a new tag in your ConvertKit account to help organize and segment your subscribers.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Tag Name | String | Yes | The name of the new tag |

## Sample Response

```json
{
  "tag": {
    "id": 12345,
    "name": "New Customer",
    "created_at": "2023-03-20T09:15:00Z"
  }
}
```

## Details

- **Type**: sdkcore.ActionTypeNormal