# Get Folder

## Description

Retrieves detailed information about a specific ClickUp folder by its ID, including contained lists and associated metadata.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| folder-id | Folder ID | String | The ID of the folder to retrieve. | Yes | - |

## Returns

This action returns a folder object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the folder |
| name | String | The name of the folder |
| orderindex | Number | The position of the folder within the space |
| override_statuses | Boolean | Whether the folder has custom statuses |
| hidden | Boolean | Whether the folder is hidden |
| space | Object | Information about the parent space |
| task_count | Number | Count of tasks within the folder |
| lists | Array | Lists contained within the folder |

Each list object contains:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the list |
| name | String | The name of the list |
| content | String | The description or content of the list |

## Example Usage

You can use this action to:

1. Retrieve folder details for display in your application
2. Get all lists contained within a folder for navigation
3. Check folder settings before creating or updating lists
4. Generate reports on folder structure and task counts
5. Build custom folder and list navigation interfaces

## Notes

- The folder must exist and be accessible to the authenticated user.
- This action retrieves all lists within the folder, which may be numerous in large folders.
- Task counts include all tasks across all lists in the folder.