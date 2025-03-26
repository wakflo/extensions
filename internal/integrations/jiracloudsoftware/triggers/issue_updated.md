# Issue Updated

## Description

Trigger a workflow when an issue is updated in Jira.

## Properties

| Name        | Type   | Required | Description                                          |
|-------------|--------|----------|------------------------------------------------------|
| project-key | string | No       | Filter updates by project key (e.g., PRJ)            |
| issue-type  | string | No       | Filter updates by issue type (e.g., Bug, Story)      |


## Details

- **Type**: sdkcore.TriggerTypePolling