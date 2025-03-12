# Delete Record

## Description

Permanently removes a record from a specified Zoho CRM module.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name     | Type   | Required | Description                                                                 |
|----------|--------|----------|-----------------------------------------------------------------------------|
| module   | string | Yes      | The Zoho CRM module where the record exists (e.g., Leads, Contacts, Accounts)|
| recordId | string | Yes      | The ID of the record to delete                                              |

## Sample JSON Data Format

```json
{
  "module": "Leads",
  "recordId": "3477061000000419001"
}