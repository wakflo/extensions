# List Subscribers

## Description

Retrieve a list of subscribers from your ConvertKit account with their details and tags.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Page | Number | Yes | The page number to retrieve (pagination), starting from 1 |


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
  "total_subscribers": 156,
  "page": 1,
  "page_size": 50
}
```

## Details

- **Type**: sdkcore.ActionTypeNormal