# List Contacts

## Description

Retrieve a list of contacts from your ActiveCampaign account with filtering options.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| email | String | No | Filter contacts by email address |
| list-id | String | No | Filter contacts by a specific list ID |
| tag-id | String | No | Filter contacts by a specific tag ID |
| limit | Number | No | Maximum number of contacts to return (default: 20, max: 100) |

## Output

The action outputs an array of contacts matching the specified filters. Each contact contains the following properties:

```json
[
  {
    "id": "123",
    "email": "sample1@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+1234567890",
    "cdate": "2023-01-15T15:30:00-05:00",
    "udate": "2023-02-20T10:15:00-05:00",
    "links": {
      "lists": "https://api.example.com/contacts/123/lists",
      "deals": "https://api.example.com/contacts/123/deals"
    }
  },
  {
    "id": "124",
    "email": "sample2@example.com",
    "firstName": "Jane",
    "lastName": "Smith",
    "phone": "+0987654321",
    "cdate": "2023-01-20T09:45:00-05:00",
    "udate": "2023-02-25T14:30:00-05:00",
    "links": {
      "lists": "https://api.example.com/contacts/124/lists",
      "deals": "https://api.example.com/contacts/124/deals"
    }
  }
]
```

## Notes

- The maximum number of contacts that can be returned in a single request is 100.
- Use the offset parameter with the limit parameter to implement pagination.
- For optimal performance, it's recommended to use specific filters rather than retrieving all contacts.