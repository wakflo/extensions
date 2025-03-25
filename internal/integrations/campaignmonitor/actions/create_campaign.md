# Create Campaign

## Description

Create a new email campaign in Campaign Monitor. This action allows you to specify the campaign content, recipients, and optionally send preview emails to test the campaign before launching.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| clientId | Client ID | String | The Client ID for which to create the campaign. If not provided, the Client ID from the authentication will be used. | No | - |
| name | Campaign Name | String | The name of the campaign (for internal reference). | Yes | - |
| subject | Subject | String | The subject line for the campaign. | Yes | - |
| fromName | From Name | String | The name that will appear in the From field. | Yes | - |
| fromEmail | From Email | String | The email address that will appear in the From field. | Yes | - |
| replyTo | Reply To | String | The reply-to email address for the campaign. | Yes | - |
| htmlContent | HTML Content | String | The HTML content of the campaign. | Yes | - |
| textContent | Text Content | String | The plain text content of the campaign. | No | - |
| listId | List ID | Array of Strings | The ID of the lists to which the campaign will be sent. | No | - |

## Returns

This action returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| CampaignID | String | The unique identifier for the newly created campaign |

## Example Usage

You can use this action to automate the creation of email campaigns based on specific triggers in your workflow. For example, you could create a monthly newsletter campaign that is automatically populated with content from your CMS or blog.

## Notes

- You must specify either List IDs or Segment IDs for the campaign to have recipients.
- The campaign is only created but not sent. To send the campaign, you'll need to use Campaign Monitor's interface or API to schedule or send it.
- If Send Preview is set to true and Preview Recipients are provided, a preview of the campaign will be sent to those email addresses immediately after creation.