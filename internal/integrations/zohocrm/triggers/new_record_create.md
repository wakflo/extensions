# Record Updated

## Description

Triggers when a new record is updated in a specified Zoho CRM module.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name   | Type   | Required | Description                                                                             |
|--------|--------|----------|-----------------------------------------------------------------------------------------|
| module | string | Yes      | The Zoho CRM module to monitor for new records (e.g., Leads, Contacts, Accounts)        |

## Sample JSON Data Format

```json
{
  "module": "Leads"
}

{
  "records": [
    {
      "id": "3477061000000419001",
      "Last_Name": "Smith",
      "First_Name": "John",
      "Email": "john.smith@example.com",
      "Created_Time": "2023-01-15T15:45:30+05:30"
    }
  ],
  "count": 1,
  "module": "Leads"
}
