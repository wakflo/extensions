# Create Contact

## Description

Create a new contact in your ActiveCampaign account with customizable fields.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| email | String | Yes | Email address of the contact |
| first-name | String | No | First name of the contact |
| last-name | String | No | Last name of the contact |
| phone | String | No | Phone number of the contact |

## Output

The action outputs the newly created contact object:

```json
{
  "id": "123",
  "email": "new@example.com",
  "firstName": "New",
  "lastName": "User",
  "phone": "+1234567890",
  "cdate": "2023-03-15T15:30:00-05:00",
  "udate": "2023-03-15T15:30:00-05:00",
  "links": {
    "lists": "https://api.example.com/contacts/123/lists",
    "deals": "https://api.example.com/contacts/123/deals"
  }
}
```

## Notes

- Email is the only required field to create a contact.
- If a contact with the provided email already exists, the API may return an error or update the existing contact depending on your ActiveCampaign account settings.