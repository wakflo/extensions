# Campaign Sent

## Description

This trigger fires when new campaigns are sent from your Campaign Monitor account. It uses polling to periodically check for campaigns that have been sent since the last time the trigger ran.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| clientId | Client ID | String | The Client ID for which to monitor sent campaigns. If not provided, the Client ID from the authentication will be used. | No | - |

## Output

This trigger returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| recentCampaigns | Array | List of campaign objects sent since the last run |
| clientId | String | The Client ID that was monitored |
| sinceDate | String | The date from which campaigns were retrieved (ISO 8601 format) |

Each campaign object in the recentCampaigns array contains:

| Field | Type | Description |
|-------|------|-------------|
| CampaignID | String | Unique identifier for the campaign |
| Name | String | The name of the campaign |
| Subject | String | The subject line of the campaign |
| FromName | String | The sender name |
| FromEmail | String | The sender email address |
| SentDate | String | The date and time when the campaign was sent |
| WebVersionURL | String | URL to the web version of the campaign |
| WebVersionTextURL | String | URL to the text version of the campaign |
| TotalRecipients | Number | Total number of recipients for the campaign |

## Example Usage

You can use this trigger to automate various workflows when campaigns are sent, such as:

1. Updating campaign records in your CRM or database
2. Notifying team members when important campaigns are sent
3. Logging campaign information for reporting and analytics
4. Triggering follow-up actions or campaigns based on the primary campaign
5. Syncing campaign data to other marketing platforms or analytics tools

## Notes

- This trigger uses polling, so there may be a slight delay between when a campaign is sent and when this trigger fires.
- The trigger uses the last run time to determine which campaigns to include. On the first run, it will include campaigns sent in the last 24 hours.
- The trigger only detects campaigns that have been sent, not campaigns that are in draft or scheduled status.
- If multiple campaigns were sent since the last run, all of them will be included in the response.