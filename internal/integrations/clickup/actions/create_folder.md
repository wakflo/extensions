# Create Folder

## Description

Creates a new folder in a specified ClickUp space to help organize lists and tasks within a hierarchical structure.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace where the folder will be created. | Yes | - |
| space-id | Space | Select | Select a space where the folder will be created. | Yes | - |
| name | Folder Name | String | The name of the folder to create. | Yes | - |

## Returns

This action returns the newly created folder object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the created folder |
| name | String | The name of the folder |
| orderindex | Number | The position of the folder within the space |
| override_statuses | Boolean | Whether the folder has custom statuses |
| hidden | Boolean | Whether the folder is hidden |
| space | Object | Information about the parent space |
| task_count | Number | Count of tasks within the folder (initially 0) |

## Example Usage

You can use this action to:

1. Create organizational structures for projects or departments
2. Set up folders for different work categories or teams
3. Build hierarchical structures programmatically
4. Create folders as part of onboarding or project setup processes

## Notes

- Folders exist within spaces and can contain multiple lists.
- Folders inherit statuses from their parent space by default.
- The authenticated user must have permission to create folders in the space.