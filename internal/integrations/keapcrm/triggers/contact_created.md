# Contact Created

## Description

Trigger a workflow when a new contact is created in Keap.

## Properties

This trigger doesn't require any configuration properties.

## Details

- **Type**: sdkcore.TriggerTypePolling
- **Icon**: mdi:account-plus-outline

## Output

Returns a list of newly created contacts since the last execution of the trigger.

## Sample Output

```json
{
  "contacts": [
    {
      "id": "123",
      "given_name": "John",
      "family_name": "Doe",
      "email": "john.doe@example.com",
      "date_created": "2023-06-15T14:30:25Z",
      "last_updated": "2023-06-15T14:30:25Z"
    }
  ]
}
```

## Notes

- The trigger uses polling to check for new contacts created since the last execution.
- If this is the first execution, it will look for contacts created in the last 24 hours.
- The trigger identifies truly new contacts by comparing the `date_created` and `last_updated` fields - when they match, it indicates a newly created contact rather than an updated one.