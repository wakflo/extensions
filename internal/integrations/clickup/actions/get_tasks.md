# Get Tasks

## Description

Retrieves all tasks from a specified ClickUp list, allowing you to access task details, statuses, and metadata.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace containing the list. | Yes | - |
| space-id | Space | Select | Select a space containing the list. | Yes | - |
| folder-id | Folder | Select | Select a folder containing the list. | Yes | - |
| list-id | List | Select | Select a list to retrieve tasks from. | Yes | - |

## Returns

This action returns an object containing an array of task objects:

| Field | Type | Description |
|-------|------|-------------|
| tasks | Array | List of task objects in the specified list |

Each task object contains:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the task |
| name | String | The name/title of the task |
| status | Object | The current status of the task including color |
| priority | Object | The priority level of the task (when set) |
| date_created | String | Timestamp when the task was created |
| date_updated | String | Timestamp when the task was last updated |
| assignees | Array | List of users assigned to the task |

## Example Usage

You can use this action to:

1. Build a custom task dashboard or reporting tool
2. Sync tasks with other project management systems
3. Generate task lists for team meetings or reviews
4. Create automated workflows based on task statuses

## Notes

- Large lists may contain many tasks, which could affect performance.
- Consider using filtering parameters if you need to retrieve specific tasks.
- Tasks are returned in their default sort order from ClickUp.