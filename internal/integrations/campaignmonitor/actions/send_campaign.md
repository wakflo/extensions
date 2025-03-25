# Send Campaign

## Description

Schedules an existing draft campaign for sending either immediately or at a specified future date and time. For campaigns with more than 5 recipients, you must have sufficient email credits, a saved credit card, or an active monthly billed account.  This action allows you to control when your email campaign is delivered to recipients.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| campaignId | Campaign ID | String | The ID of the draft campaign to send. | Yes | - |
| confirmationEmail | Confirmation Emails | Array of Strings | Email addresses to receive confirmation when the campaign is sent. | No | - |
| sendDate | Send Date | DateTime | The date and time to send the campaign (format: YYYY-MM-DD HH:MM). Leave blank to send immediately. | No | - |

## Returns

This action returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| success | Boolean | Indicates whether the campaign was successfully scheduled |
| message | String | A message confirming the campaign was scheduled for sending |
| CampaignID | String | The ID of the campaign that was scheduled |

## Notes

- **Free sending limitation:** When sending campaigns to 5 or fewer recipients (which are free of charge), you can send to a maximum of 50 unique email addresses per day.
- For campaigns with more than 5 recipients, you must have sufficient email credits, a saved credit card, or an active monthly billed account.
- When sending to segments, to ensure they have finished calculating, it's recommended to wait approximately one hour after importing subscribers, creating segments, or updating segment rules.
- If no send date is specified, the campaign will be sent immediately.

## Example Usage

You can use this action to:

1. Schedule campaigns to be sent at optimal times for your audience
2. Set up automated campaign sending as part of a larger workflow
3. Schedule campaigns in advance for product launches, announcements, or newsletters
4. Send immediate notifications or time-sensitive updates to your subscribers

To use this action, you'll need the Campaign ID of a draft campaign that has already been created using the Create Campaign or Create Campaign From Template actions.