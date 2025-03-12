# Contact Updated

## Description

The "Contact Updated" trigger monitors your HubSpot CRM for any contacts that have been created or modified. When a change is detected, this trigger will fire and allow you to initiate a series of automated tasks, such as updating records in other systems, sending notifications, or creating tasks for your team.

## Properties

| Name       | Type   | Required | Description                                                                                        |
|------------|--------|----------|----------------------------------------------------------------------------------------------------|
| properties | string | No       | Comma-separated list of contact properties to include in the response (e.g., firstname,lastname,email) |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Output

This trigger outputs a list of recently modified contacts with their properties. The output structure will include:

```json
{
  "results": [
    {
      "id": "51",
      "properties": {
        "firstname": "John",
        "lastname": "Doe",
        "email": "john.doe@example.com",
        "phone": "+1234567890",
        "createdate": "2023-03-15T09:31:40.678Z",
        "lastmodifieddate": "2023-03-15T10:45:12.412Z"
      },
      "createdAt": "2023-03-15T09:31:40.678Z",
      "updatedAt": "2023-03-15T10:45:12.412Z"
    }
  ]
}
```

## Notes

- The trigger polls the HubSpot API for contacts that have been updated since the last time the trigger ran.
- If you don't specify any properties, the response will include all default properties.
- You can optimize performance by requesting only the specific properties you need.
- HubSpot API has rate limits, so be mindful of how frequently this trigger runs.