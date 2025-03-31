# Form Response

## Description

Triggers when a new response is submitted to a specified form, providing all answer data and respondent information.

## Details

- **Type**: sdkcore.TriggerTypePolling

## Properties

| Name    | Type   | Required | Description                                     |
| ------- | ------ | -------- | ----------------------------------------------- |
| form-id | string | Yes      | The ID of the form to monitor for new responses |

## Returns

A JSON object containing:

- `total_items`: Number of new responses
- `page_count`: Number of pages
- `items`: Array of response objects, each containing:
  - Response ID
  - Submission timestamp
  - Landing timestamp
  - Metadata (user agent, platform, etc.)
  - Answers array with question details and responses
  - Calculated values (scores, etc.)
  - Hidden field values

## Usage Notes

This trigger uses polling to check for new responses. It will only trigger for responses submitted after the last time the trigger ran.
