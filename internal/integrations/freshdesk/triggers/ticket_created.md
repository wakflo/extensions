# Ticket Created

## Description

Trigger a workflow when a new ticket is created in Freshdesk.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Sample Response

```json
{
  "id": 123,
  "subject": "New Support Request",
  "description": "I need help with your product",
  "status": 2,
  "priority": 1,
  "requester_id": 456,
  "created_at": "2023-12-01T12:30:45Z"
}
```

## Notes

This trigger polls the Freshdesk API at regular intervals to check for new tickets. By default, it retrieves tickets created in the last 24 hours on the first run, and then tickets created since the last polling time on subsequent runs.
