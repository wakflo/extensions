# Email Sent Trigger

Polls Mailjet to detect newly sent emails and triggers a workflow for each one.

## Configuration Options

- **From Email**: Filter by sender email address
- **To Email**: Filter by recipient email address
- **Custom ID**: Filter by email Custom ID
- **Custom Campaign**: Filter by campaign name
- **Limit**: Maximum number of emails to process per poll (default: 50, max: 1000)

## Output

Each triggered workflow receives information about the sent email including:
- Message ID and UUID
- Sender and recipient details
- Subject and timestamp
- Campaign information
- Delivery status and metrics

## Notes

- This trigger uses polling to check for new emails at regular intervals
- Only emails sent since the last poll will trigger the workflow
- All filter parameters are optional