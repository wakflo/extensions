# New Response

## Description

Triggers when a new response is submitted to a specific survey in your SurveyMonkey account. This trigger polls the SurveyMonkey API at regular intervals to check for new responses since the last run.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name      | Type   | Description                                       | Required | Default |
| --------- | ------ | ------------------------------------------------- | -------- | ------- |
| survey_id | String | The ID of the survey to monitor for new responses | Yes      | -       |

## Output

This trigger outputs an array of response objects that were created since the last trigger execution:

```json
[
	{
		"id": "9876543210",
		"survey_id": "123456789",
		"collector_id": "12345",
		"response_status": "completed",
		"date_created": "2023-04-15T14:22:35+00:00",
		"date_modified": "2023-04-15T14:30:10+00:00",
		"ip_address": "192.168.1.1",
		"href": "https://api.surveymonkey.com/v3/surveys/123456789/responses/9876543210"
	}
]
```
