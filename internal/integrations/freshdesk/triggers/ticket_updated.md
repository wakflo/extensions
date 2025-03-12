# Ticket Updated

## Description

Trigger a workflow when a ticket is updated in Freshdesk, including status changes, priority updates, or note additions.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Sample Response

```json
{
  "id": 123,
  "subject": "Support Request",
  "description": "I need help with your product",
  "status": 3,
  "priority": 2,
  "requester_id": 456,
  "created_at": "2023-12-01T12:30:45Z",
  "updated_at": "2023-12-01T14:45:20Z"
}
```

## Notes

This trigger polls the Freshdesk API at regular intervals to check for updated tickets. By default, it retrieves tickets updated in the last 24 hours on the first run, and then tickets updated since the last polling time on subsequent runs.