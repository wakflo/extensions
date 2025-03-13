# Get Contact
## Description
Retrieve detailed information about a specific contact in Keap by their unique Contact ID.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `contact_id` | String | Yes | Unique identifier of the contact to retrieve |


## Details
- **Type**: sdkcore.ActionTypeNormal
- **Icon**: mdi:account-search

## Potential Errors
- **400**: Bad Request - Invalid contact ID or parameters
- **401**: Unauthorized - Invalid or expired access token
- **404**: Not Found - Contact does not exist

## Sample Output
```json
{
  "id": "12345",
  "given_name": "John",
  "family_name": "Doe",
  "email_addresses": [
    {
      "email": "john.doe@example.com",
      "type": "PRIMARY"
    }
  ],
  "phone_numbers": [
    {
      "number": "+1-555-123-4567", 
      "type": "MOBILE"
    }
  ],
  "last_updated": "2023-06-15T10:30:00Z"
}
```

## Behavior
- Requires a valid Contact ID to retrieve contact details
- Returns comprehensive contact information
- Provides detailed error handling for various scenarios