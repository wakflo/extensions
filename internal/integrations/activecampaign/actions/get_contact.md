# Get Contact

## Description

Retrieve a specific contact by ID from your ActiveCampaign account.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| contact-id | String | Yes | The ID of the contact you want to retrieve |

## Output

The action outputs a single contact object with detailed information:

```json
{
  "id": "123",
  "email": "sample@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "phone": "+1234567890",
  "cdate": "2023-01-15T15:30:00-05:00",
  "udate": "2023-02-20T10:15:00-05:00",
  "links": {
    "lists": "https://api.example.com/contacts/123/lists",
    "deals": "https://api.example.com/contacts/123/deals"
  },
  "fieldValues": [
    {
      "field": "1",
      "value": "Value 1"
    },
    {
      "field": "2",
      "value": "Value 2"
    }
  ]
}
```

## Notes

- This action requires a valid contact ID to function correctly.
- The contact ID can be obtained from the List Contacts action or the Contact Updated trigger.
- If the contact ID does not exist, the API will return an error.