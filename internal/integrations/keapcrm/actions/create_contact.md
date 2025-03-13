# Create Contact

## Description

Create a new contact in Keap with specified details.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `given_name` | String | Yes | Contact's first name |
| `family_name` | String | Yes | Contact's last name |
| `email` | String | Yes | Contact's email address |
| `phone_number` | String | No | Contact's phone number |
| `phone_type` | Select | No | Type of phone number (options: WORK, HOME, MOBILE, OTHER) |
| `address_line1` | String | No | First line of contact's address |
| `address_line2` | String | No | Second line of contact's address |
| `locality` | String | No | City/Locality of contact's address |
| `region` | String | No | State/Region of contact's address |
| `postal_code` | String | No | Postal/ZIP code of contact's address |
| `country` | String | No | Country of contact's address |
| `address_type` | Select | No | Type of address (options: BILLING, SHIPPING, HOME, WORK, OTHER) |
| `company_name` | String | No | Name of the contact's company |
| `job_title` | String | No | Contact's job title |
| `opt_in_reason` | String | No | Reason for the contact opting in to communications |

## Details

- **Type**: sdkcore.ActionTypeNormal
- **Icon**: mdi:account-plus

## Output

Returns the newly created contact object from Keap with all the provided information and a generated contact ID.

## Sample Output

```json
{
  "id": "12345",
  "given_name": "John",
  "family_name": "Doe",
  "email": "john.doe@example.com",
  "phone_numbers": [
    {
      "type": "WORK",
      "number": "+1 555-123-4567"
    }
  ],
  "addresses": [
    {
      "type": "BILLING",
      "line1": "123 Main St",
      "line2": "Suite 100",
      "locality": "Anytown",
      "region": "CA",
      "postal_code": "12345",
      "country": "USA"
    }
  ],
  "company": {
    "name": "Acme Inc"
  },
  "job_title": "CEO"
}
```