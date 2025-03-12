# Get Record

## Description

Retrieves a specific record from a Zoho CRM module based on the record ID.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name     | Type   | Required | Description                                                                 |
|----------|--------|----------|-----------------------------------------------------------------------------|
| module   | string | Yes      | The Zoho CRM module where the record exists (e.g., Leads, Contacts, Accounts)|
| recordId | string | Yes      | The ID of the record to retrieve                                            |

## Sample JSON Data Format

```json
{
  "module": "Leads",
  "recordId": "3477061000000419001"
}