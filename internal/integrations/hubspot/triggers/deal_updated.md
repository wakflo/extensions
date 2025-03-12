# Deal Updated

## Description

The "Deal Updated" trigger monitors your HubSpot CRM for any deals that have been created or modified. When a change is detected, this trigger will fire and allow you to initiate a series of automated tasks, such as updating records in other systems, sending notifications to sales teams, or creating follow-up tasks.

## Properties

| Name       | Type   | Required | Description                                                                                     |
|------------|--------|----------|-------------------------------------------------------------------------------------------------|
| properties | string | No       | Comma-separated list of deal properties to include in the response (e.g., dealname,amount,dealstage) |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Output

This trigger outputs a list of recently modified deals with their properties. The output structure will include:

```json
{
  "results": [
    {
      "id": "123456",
      "properties": {
        "dealname": "New Enterprise Deal",
        "amount": "50000",
        "dealstage": "qualifiedtobuy",
        "pipeline": "default",
        "closedate": "2023-12-31T23:59:59.999Z",
        "createdate": "2023-03-15T09:31:40.678Z",
        "hs_lastmodifieddate": "2023-03-15T10:45:12.412Z"
      },
      "createdAt": "2023-03-15T09:31:40.678Z",
      "updatedAt": "2023-03-15T10:45:12.412Z"
    }
  ]
}
```

## Notes

- The trigger polls the HubSpot API for deals that have been updated since the last time the trigger ran.
- If you don't specify any properties, the response will include all default properties.
- You can optimize performance by requesting only the specific properties you need.
- HubSpot API has rate limits, so be mindful of how frequently this trigger runs.
- Deal stages and pipelines are specific to your HubSpot account configuration.