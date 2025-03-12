# Create Ticket

## Description

Create a new support ticket in the Freshdesk system with customizable fields.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| subject | string | Yes | The subject of the ticket |
| description | string | Yes | The description or content of the ticket |
| email | string | Yes | The email address of the ticket requester |
| priority | number | Yes | The priority level of the ticket (1-4, where 1=Low, 2=Medium, 3=High, 4=Urgent) |
| status | number | No | The status of the ticket (2=Open, 3=Pending, 4=Resolved, 5=Closed) |
| CCEmails | string | No | Emails to Copy for ticket |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Sample Response

```json
{
  "id": 123,
  "subject": "Sample Ticket",
  "description": "This is a sample ticket created via API",
  "status": 2,
  "priority": 1,
  "requester_id": 456,
  "created_at": "2023-12-01T12:30:45Z",
  "updated_at": "2023-12-01T12:30:45Z"
}
```