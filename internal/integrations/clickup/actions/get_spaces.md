# Get Spaces

## Description

Retrieves all spaces within a ClickUp workspace that are accessible to the authenticated user.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace to retrieve spaces from. | Yes | - |

## Returns

This action returns an object containing an array of space objects:

| Field | Type | Description |
|-------|------|-------------|
| spaces | Array | List of space objects in the specified workspace |

Each space object contains:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the space |
| name | String | The name of the space |
| private | Boolean | Whether the space is private |
| statuses | Array | List of status options available in the space |

## Example Usage

You can use this action to:

1. Populate a dropdown for space selection in your application
2. Map out the organizational structure within a workspace
3. Check available spaces before creating new folders or lists
4. Perform bulk operations across multiple spaces

## Notes

- Only spaces accessible to the authenticated user will be returned.
- Private spaces will only be included if the authenticated user has access to them.