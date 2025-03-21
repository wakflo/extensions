# Update Space

## Description

Updates an existing space in ClickUp with modified details such as name and privacy settings.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace containing the space. | Yes | - |
| space-id | Space | Select | Select a space to update. | Yes | - |
| name | Space Name | String | The updated name of the space. | No | - |

## Returns

This action returns the updated space object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the space |
| name | String | The updated name of the space |
| private | Boolean | The updated privacy setting of the space |
| statuses | Array | Status options available in the space |
| multiple_assignees | Boolean | Whether multiple assignees are allowed on tasks |
| features | Object | Features enabled for the space |

## Example Usage

You can use this action to:

1. Rename spaces as projects evolve or change scope
2. Update privacy settings based on organizational needs
3. Modify space settings programmatically based on triggers
4. Apply bulk updates to spaces through automation
5. Standardize space configurations across a workspace

## Notes

- Only the fields you specify will be updated; unspecified fields will remain unchanged.
- Changing a space from public to private may affect visibility