# Get Folders

## Description

Retrieves all folders within a specified ClickUp space, providing a view of the organizational structure and folder metadata.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace containing the space. | Yes | - |
| space-id | Space | Select | Select a space to retrieve folders from. | Yes | - |

## Returns

This action returns an object containing an array of folder objects:

| Field | Type | Description |
|-------|------|-------------|
| folders | Array | List of folder objects in the specified space |

Each folder object contains:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the folder |
| name | String | The name of the folder |
| orderindex | Number | The position of the folder within the space |
| override_statuses | Boolean | Whether the folder has custom statuses |
| hidden | Boolean | Whether the folder is hidden |
| task_count | Number | Count of tasks within the folder |

## Example Usage

You can use this action to:

1. Display the folder structure of a space in your application
2. Generate a navigation menu for folders in a space
3. Check available folders before creating lists
4. Perform operations across multiple folders
5. Create reports on folder organization and task distribution

## Notes

- This action returns only the folder metadata, not the lists contained within each folder.
- Use the Get Folder action to retrieve lists within a specific folder.
- Hidden folders may still be included in the results depending on user permissions.