# Add Subscriber

## Description

Add a new subscriber to a specific Campaign Monitor list. This action allows you to create or update a subscriber with custom fields, consent settings, and resubscription options.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| listId | List ID | String | The ID of the list to which the subscriber will be added. | Yes | - |
| email | Email | String | The email address of the subscriber. | Yes | - |
| name | Name | String | The name of the subscriber. | No | - |
| customFields | Custom Fields | Array | Custom fields to add for the subscriber. Each custom field has a key and value. | No | - |
| consentToTrack | Consent To Track | Select | Whether the subscriber has given consent to track their opens and clicks (Yes, No, Unchanged). | Yes | Yes |
| consentToSendSMS | Consent To Send SMS | Select | Whether the subscriber has given consent to receive SMS messages (Yes, No, Unchanged). | No | - |
| resubscribe | Resubscribe | Boolean | Whether to resubscribe the subscriber if they have previously unsubscribed. | No | false |
| restartAutoresponders | Restart Autoresponders | Boolean | Whether to restart autoresponders for the subscriber. | No | false |

### Custom Fields Object

Each custom field in the customFields array has the following properties:

| Name | Display Name | Type | Description | Required |
|------|--------------|------|-------------|----------|
| key | Key | String | The key/name of the custom field. | Yes |
| value | Value | String | The value of the custom field. | Yes |

## Returns

This action returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| EmailAddress | String | The email address of the added subscriber |
| Name | String | The name of the added subscriber |
| Status | String | The status of the subscriber (usually "Active") |
| Message | String | A success message confirming the subscriber was added |

## Example Usage

You can use this action to automatically add subscribers to your Campaign Monitor lists from various sources, such as:

1. Form submissions on your website
2. Customer registrations in your application
3. Event registrations or webinar sign-ups
4. User import from other systems (CRM, e-commerce platforms, etc.)

By including custom fields, you can capture additional information about your subscribers to enable better segmentation and personalization in your email campaigns.

## Notes

- If a subscriber with the provided email address already exists in the list, their information will be updated.
- Setting `resubscribe` to `true` will resubscribe users who have previously unsubscribed. Use this option carefully to avoid violating email marketing regulations.
- The `restartAutoresponders` option will restart any autoresponder sequences for the subscriber, which could result in them receiving emails they've already received before.
- Make sure to comply with privacy regulations (like GDPR, CCPA) when adding subscribers, especially when using the resubscribe feature.