# Contact Updated Trigger

Polls Mailjet to detect contacts that have been created or updated since the last poll and triggers a workflow for each one.

## Configuration Options

- **Contacts List ID**: Filter contacts that belong to a specific list
- **Exclude From Campaigns**: Filter contacts based on their exclusion status
- **Limit**: Maximum number of contacts to process per poll (default: 50, max: 1000)
- **Lookback Hours**: For the first run, how many hours back to check (default: 24)

## Output

Each triggered workflow receives information about the contact including:
- Contact ID
- Email address
- Name (if available)
- Creation/update timestamps
- Campaign exclusion status
- Activity metrics

## Notes

- This trigger uses polling to check for updated contacts at regular intervals
- Only contacts updated since the last poll will trigger the workflow
- All filter parameters are optional