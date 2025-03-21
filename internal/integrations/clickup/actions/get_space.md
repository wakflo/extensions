# Get Space

## Description

Retrieves details of a specific ClickUp space by ID, including its name, status options, privacy settings, and enabled features.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace where the space is located. | Yes | - |
| space-id | Space | Select | Select a space to retrieve details for. | Yes | - |

## Returns

This action returns a space object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the space |
| name | String | The name of the space |
| private | Boolean | Whether the space is private |
| statuses | Array | List of status options available in the space |
| multiple_assignees | Boolean | Whether multiple assignees are allowed on tasks |
| features | Object | Features enabled for the space (due dates, time tracking, tags, etc.) |

## Example Usage

You can use this action to retrieve space details for:

1. Displaying space information in a dashboard
2. Checking available statuses before creating or updating tasks
3. Verifying space settings before performing other operations
4. Integrating space information with other systems

## Notes

- The space must exist in the specified workspace or the API will return an error.
- Private spaces may not be accessible depending on your authentication permissions.