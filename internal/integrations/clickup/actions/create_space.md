# Create Space

## Description

Creates a new space in a specified ClickUp workspace with customizable name and privacy settings.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace where the space will be created. | Yes | - |
| name | Space Name | String | The name of the space to create. | Yes | - |

## Returns

This action returns the newly created space object containing:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the created space |
| name | String | The name of the space |
| private | Boolean | Whether the space is private |
| statuses | Array | Default status options available in the space |
| multiple_assignees | Boolean | Whether multiple assignees are allowed on tasks |
| features | Object | Features enabled for the space (due dates, time tracking, tags, etc.) |

## Example Usage

You can use this action to:

1. Create spaces programmatically for new projects or departments
2. Set up standardized workspace structures
3. Create dedicated spaces for clients or external stakeholders
4. Automate space creation as part of onboarding workflows

## Notes

- Private spaces are only visible to users who have been explicitly given access.
- New spaces come with default statuses which can be customized after creation.
- Each space will have default features enabled based on your ClickUp plan.