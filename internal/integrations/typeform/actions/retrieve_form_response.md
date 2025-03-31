# Get Responses

## Description

Retrieves responses submitted to a specific form with answer details.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name      | Type   | Required | Description                                                                          |
| --------- | ------ | -------- | ------------------------------------------------------------------------------------ |
| form-id   | string | Yes      | The ID of the form to get responses from                                             |
| page-size | number | No       | Number of responses per page (default: 25, max: 1000)                                |
| since     | string | No       | Limit to responses submitted after this date (ISO 8601 format: YYYY-MM-DDThh:mm:ss)  |
| until     | string | No       | Limit to responses submitted before this date (ISO 8601 format: YYYY-MM-DDThh:mm:ss) |
| after-id  | string | No       | Limit to responses submitted after a specific response ID (for pagination)           |
| before-id | string | No       | Limit to responses submitted before a specific response ID (for pagination)          |
