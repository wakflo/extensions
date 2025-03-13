# Contact Updated

## Description

Trigger a workflow when a contact's information is updated in Keap.

## Properties

This trigger doesn't require any configuration properties.

## Details

- **Type**: sdkcore.TriggerTypePolling
- **Icon**: mdi:account-edit-outline

## Output

Returns a list of contacts that have been updated since the last execution of the trigger.

## Sample Output

```json
{
  "contacts": [
    {
      "id": "456",
      "given_name": "Jane",
      "family_name": "Smith",
      "email": "jane.smith@example.com",
      "date_created": "2023-05-10T09:15:30Z",
      "last_updated": "2023-06-15T16:45:20Z"
    }
  ]
}
```

## Notes

- The trigger uses polling to check for updated contacts since the last execution.
- If this is the first execution, it will look for contacts updated in the last 24 hours.
- The trigger identifies updated contacts by comparing the `date_created` and `last_updated` fields - when they're different, it indicates an updated contact rather than a newly created one.