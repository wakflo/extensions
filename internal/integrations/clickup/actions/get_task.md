# Get Task

## Description

Retrieves detailed information about a specific ClickUp task by its ID, including all metadata, custom fields, and history.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| task-id | Task ID | String | The ID of the task to retrieve. | Yes | - |

## Returns

This action returns a comprehensive task object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the task |
| name | String | The name/title of the task |
| description | String | The detailed description of the task |
| status | Object | The current status of the task including color |
| priority | Object | The priority level of the task (when set) |
| date_created | String | Timestamp when the task was created |
| date_updated | String | Timestamp when the task was last updated |
| creator | Object | Information about the user who created the task |
| assignees | Array | List of users assigned to the task |
| custom_fields | Array | Custom fields and their values for the task |

## Example Usage

You can use this action to:

1. Display detailed task information in your application
2. Create task detail pages or modals
3. Track specific tasks within automated workflows
4. Generate detailed reports for individual tasks
5. Retrieve task custom fields for integration with other systems

## Notes

- This action retrieves a single task and is more efficient than retrieving all tasks when you only need information about one specific task.
- The task ID must be valid, or the API will return an error.