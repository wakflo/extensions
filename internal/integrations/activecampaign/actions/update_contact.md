# Update Contact

## Description

Update an existing contact in your ActiveCampaign account.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| contact-id | String | Yes | The ID of the contact you want to update |
| email | String | No | Email address of the contact |
| first-name | String | No | First name of the contact |
| last-name | String | No | Last name of the contact |
| phone | String | No | Phone number of the contact |
| custom-fields | Object | No | Key-value pairs of custom field IDs and their values |

## Output

The action outputs the updated contact object:

```json
{
  "id": "123",
  "email": "updated@example.com",
  "firstName": "Updated",
  "lastName": "User",
  "phone": "+1234567890",
  "cdate": "2023-01-15T15:30:00-05:00",
  "udate": "2023-03-15T16:45:00-05:00",
  "links": {
    "lists": "https://api.example.com/contacts/123/lists",
    "deals": "https://api.example.com/contacts/123/deals"
  }
}
```

## Notes

- The contact ID is required to identify which contact to update.
- At least one field (email, first-name, last-name, phone) must be provided for the update.
- If you want to update the contact's lists, provide a comma-separated string of list IDs in the `list-ids` field.
- If you want to update the contact's tags, provide a comma-separated string of tag IDs in the `tag-ids` field.
- Custom fields should be provided as an object where the keys are the custom field IDs and the values are the field values.

### Example Custom Fields Object

```json
{
  "1": "Updated value for custom field 1",
  "2": "Updated value for custom field 2"
}
```