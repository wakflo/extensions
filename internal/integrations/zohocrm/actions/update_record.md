# Update Record

## Description

Updates an existing record in a specified Zoho CRM module with the provided data.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name     | Type   | Required | Description                                                                          |
|----------|--------|----------|--------------------------------------------------------------------------------------|
| module   | string | Yes      | The Zoho CRM module where the record exists (e.g., Leads, Contacts, Accounts)        |
| recordId | string | Yes      | The ID of the record to update                                                       |
| data     | object | Yes      | The updated data for the record in JSON format                                       |

## Sample JSON Data Format

```json
{
  "module": "Leads",
  "recordId": "3477061000000419001",
  "data": {
    "Last_Name": "Smith",
    "Email": "john.smith.updated@example.com",
    "Lead_Status": "Qualified"
  }
}