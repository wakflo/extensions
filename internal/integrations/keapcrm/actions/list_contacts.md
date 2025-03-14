# List Contacts
## Description
Retrieve a list of contacts from Keap with optional filtering, sorting, and additional property selection.

## Properties
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `limit` | Number | No | Maximum number of contacts to return (default: 50, max: 1000) |
| `offset` | Number | No | Number of contacts to skip for pagination |
| `email` | String | No | Filter contacts by email address |
| `first_name` | String | No | Filter contacts by first name |
| `last_name` | String | No | Filter contacts by last name |
| `order` | Select | No | Field to order contacts by (options: email, given_name, family_name, last_updated) |
| `order_direction` | Select | No | Direction of ordering (ascending or descending) |

### Optional Properties
Select additional fields to include in the response:
- Email
- Phone Number
- Lead Source
- Job Title
- Last Updated
- Tags

### Ordering Options
- Order By: Email, First Name, Last Name, Last Updated
- Order Direction: Ascending, Descending (default: Ascending)

## Details
- **Type**: sdkcore.ActionTypeNormal
- **Icon**: mdi:account-multiple

## Potential Errors
- **400**: Bad Request - Invalid parameters
- **401**: Unauthorized - Invalid or expired access token

## Sample Output
```json
{
  "contacts": [
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
      "last_updated": "2023-01-15T10:30:00Z"
    },
    {
      "id": "67890",
      "given_name": "Jane",
      "family_name": "Smith",
      "email_addresses": [
        {
          "email": "jane.smith@example.com",
          "type": "PRIMARY"
        }
      ],
      "last_updated": "2023-02-20T14:45:00Z"
    }
  ],
  "total": 2,
  "limit": 50,
  "offset": 0
}
```

## Behavior
- Default limit is 50 contacts if not specified
- Supports flexible filtering by name, email
- Allows adding extra properties to the response
- Provides pagination through limit and offset
- Supports sorting of results by various fields