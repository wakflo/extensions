# Get List

## Description

Retrieves detailed information about a specific ClickUp list by its ID, including name, content, and associated task statistics.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| list-id | List ID | String | The ID of the list to retrieve. | Yes | - |

## Returns

This action returns a list object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the list |
| name | String | The name of the list |
| content | String | The description or content of the list |
| statuses | Array | Status options available in the list |
| task_count | Number | Count of tasks within the list |
| orderindex | Number | The position of the list within its folder |
| folder | Object | Information about the parent folder (if applicable) |
| space | Object | Information about the parent space |

## Example Usage

You can use this action to:

1. Retrieve list details for display in your application
2. Check list settings before creating or updating tasks
3. Get task counts for reporting purposes
4. Verify list statuses for workflow integration
5. Retrieve list information for documentation or auditing

## Notes

- Lists can exist either within a folder or directly in a space (folderless).
- The list must exist and be accessible to the authenticated user.
- Task counts may include tasks in various statuses including closed tasks.