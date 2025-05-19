# Create Contact List

## Description

Creates a new empty contact list in your SurveyMonkey account for email collectors. This action creates the list structure that contacts can be added to later using other API methods.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Type   | Description                  | Required | Default |
| ---- | ------ | ---------------------------- | -------- | ------- |
| name | String | The name of the contact list | Yes      | -       |

## Response

This action returns the newly created contact list object:

```json
{
	"id": "12345678",
	"name": "Marketing Contacts",
	"contact_count": 0,
	"href": "https://api.surveymonkey.com/v3/contact_lists/12345678",
	"date_created": "2023-04-25T14:30:45+00:00",
	"date_modified": "2023-04-25T14:30:45+00:00"
}
```
