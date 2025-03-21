# Update Task

## Description

Updates an existing task in ClickUp with modified details such as name, description, status, priority, and assignees.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| task-id | Task ID | String | The ID of the task to update. | Yes | - |
| name | Task Name | String | The updated name of the task. | No | - |
| description | Task Description | Text | The updated description of the task. | No | - |
| status | Status | String | The updated status of the task. | No | - |
| priority | Priority | Select | The updated priority level of the task. | No | - |
| assignee-id | Assignee ID | String | The ID of the user to assign to the task. | No | - |

## Returns

This action returns the updated task object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the task |
| name | String | The updated name/title of the task |
| description | String | The updated description of the task |
| status | Object | The updated status of the task including color |
| priority | Object | The updated priority level of the task |
| date_updated | String | Timestamp when the task was updated |
| assignees | Array | Updated list of users assigned to the task |

## Example Usage

You can use this action to:

1. Update task details from external systems
2. Change task statuses as part of automated workflows
3. Reassign tasks based on triggers or events
4. Modify task information based on user actions
5. Update tasks in bulk through programmatic processes

## Notes

- Only the fields you specify will be updated; unspecified fields will remain unchanged.
- Status values must match the available statuses in the task's list.
- Priority values: 1 (Urgent), 2 (High), 3 (Normal), 4 (Low).
- The task ID must be valid, or the API will return an error.