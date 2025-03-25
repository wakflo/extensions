# Create Folderless List

## Description

Creates a new list directly in a space without a parent folder, allowing for a flatter organization structure in ClickUp.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace where the list will be created. | Yes | - |
| space-id | Space | Select | Select a space where the list will be created. | Yes | - |
| name | List Name | String | The name of the list. | Yes | - |

## Returns

This action returns the newly created list object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the created list |
| name | String | The name of the list |
| content | String | The description or content of the list |
| orderindex | Number | The position of the list within the space |
| status | Object | The default status configuration for the list |

## Example Usage

You can use this action to:

1. Create top-level lists for specific projects or initiatives
2. Set up lists that don't fit within the folder hierarchy
3. Create flat organizational structures for simpler workflows
4. Generate lists programmatically based on external triggers

## Notes

- Folderless lists exist directly in a space, not within any folder.
- Lists inherit the statuses from their parent space unless custom statuses are specified.
- You can later move a folderless list into a folder if your organizational needs change.