# Create Ticket

## Description

Create a new ticket in HubSpot. This action allows you to add new tickets to your HubSpot CRM by providing their details.

## Properties

| Name     | Type   | Required | Description                 |
|----------|--------|----------|-----------------------------|
| subject  | string | Yes      | Subject of the ticket       |
| content  | string | No       | Ticket description          |
| priority | string | No       | Ticket Priority             |


## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the created ticket details from HubSpot. The structure will include:

```json
{
  "id": "1234567890",
  "subject": "Ticket Subject",
  "content": "Ticket Description",
  "priority": "NORMAL"
}
```

## Notes

- Subject is the only required field for creating a contact in HubSpot
- If the priority is not provided, the default priority will be set to `NORMAL`
- The `content` field is optional and can be used to provide additional information about the ticket
- The `priority` field is optional and can be set to `HIGH`, `MEDIUM`, or `LOW`
- The `id` field in the output represents the ticket ID in HubSpot