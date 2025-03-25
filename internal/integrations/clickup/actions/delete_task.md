# Delete Task

## Description

Deletes a task from ClickUp permanently, removing it from its list and all associated views.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| task-id | Task ID | String | The ID of the task to delete. | Yes | - |

## Returns

This action returns a simple status object:

| Field | Type | Description |
|-------|------|-------------|
| success | Boolean | True if the task was successfully deleted |

## Example Usage

You can use this action to:

1. Remove completed or obsolete tasks programmatically
2. Clean up tasks as part of automated workflows
3. Delete tasks that meet certain criteria
4. Implement "purge" functionality in custom applications
5. Remove test or temporary tasks

## Notes

- Deleting a task is permanent and cannot be undone.
- All subtasks, comments, and attachments will also be removed.
- The authenticated user must have permission to delete tasks in the list.
- Consider archiving tasks instead of deleting if you may need the information later.