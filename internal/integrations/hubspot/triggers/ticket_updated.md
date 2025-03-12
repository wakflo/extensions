# Ticket Created

## Description

The "Ticket Created" trigger monitors your HubSpot CRM for newly created tickets. When a new ticket is detected, this trigger will fire and allow you to initiate a series of automated tasks, such as sending notifications, creating follow-up tasks, or syncing with other systems.

## Properties

| Name       | Type   | Required | Description                                                                                        |
|------------|--------|----------|----------------------------------------------------------------------------------------------------|
| properties | string | No       | Comma-separated list of ticket properties to include in the response (e.g., subject,hs_ticket_priority,hs_pipeline_stage) |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Output

This trigger outputs a list of newly created tickets with their properties. The output structure will include:

```json
{
  "results": [
    {
      "id": "123",
      "properties": {
        "subject": "Customer Support Request",
        "hs_ticket_priority": "HIGH",
        "hs_pipeline_stage": "New",
        "createdate": "2024-03-10T14:30:45.123Z"
      },
      "createdAt": "2024-03-10T14:30:45.123Z",
      "updatedAt": "2024-03-10T14:30:45.123Z"
    }
  ]
}
```

## Notes

- The trigger polls the HubSpot API for tickets that have been created since the last time the trigger ran.
- If you don't specify any properties, the response will include default ticket properties.
- You can optimize performance by requesting only the specific properties you need.
- HubSpot API has rate limits, so be mindful of how frequently this trigger runs.
- The trigger will return up to 100 newly created tickets per run.