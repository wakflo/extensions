# Create Record

## Description

Creates a new record in a specified Zoho CRM module with the provided data.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name   | Type   | Required | Description                                                                           |
|--------|--------|----------|---------------------------------------------------------------------------------------|
| module | string | Yes      | The Zoho CRM module where the record will be created (e.g., Leads, Contacts, Accounts)|
| data   | object | Yes      | The data for the new record in JSON format                                            |

## Sample JSON Data Format

```json
{
  "module": "Leads",
  "data": {
    "Last_Name": "Smith",
    "First_Name": "John",
    "Email": "john.smith@example.com",
    "Company": "Example Inc.",
    "Lead_Source": "Website",
    "Phone": "+1-123-456-7890"
  }
}