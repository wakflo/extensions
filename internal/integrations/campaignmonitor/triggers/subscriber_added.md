# Subscriber Added

## Description

This trigger fires when new subscribers are added to a specified Campaign Monitor list. It uses polling to periodically check for new subscribers that have been added since the last time the trigger ran.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| listId | List ID | String | The ID of the list to monitor for new subscribers. | Yes | - |

## Output

This trigger returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| newSubscribers | Array | List of subscriber objects added since the last run |
| listId | String | The ID of the list that was monitored |
| sinceDate | String | The date from which subscribers were retrieved (format: YYYY-MM-DD) |

Each subscriber object in the newSubscribers array contains:

| Field | Type | Description |
|-------|------|-------------|
| EmailAddress | String | The subscriber's email address |
| Name | String | The subscriber's name |
| Date | String | The date the subscriber was added |
| State | String | The subscriber's state (usually "Active") |
| CustomFields | Array | List of custom field objects with Key and Value properties |

## Example Usage

You can use this trigger to automate various workflows when new subscribers are added to your lists, such as:

1. Sending welcome emails or onboarding sequences
2. Adding subscribers to your CRM system
3. Updating dashboards or reports with new subscriber information
4. Notifying team members about new subscribers
5. Triggering follow-up actions based on subscriber information

## Notes

- This trigger uses polling, so there may be a slight delay between when a subscriber is added to the list and when this trigger fires.
- The trigger uses the last run time to determine which subscribers to retrieve. On the first run, it will retrieve subscribers added in the last 24 hours.
- Only active subscribers are returned by this trigger. Unsubscribed, bounced, or deleted subscribers are not included.
- The trigger may return a large number of subscribers if many were added since the last run.