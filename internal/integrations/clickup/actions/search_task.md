# Search Task

## Description

Searches for tasks across a ClickUp workspace based on a query string, allowing you to find tasks that match specific criteria.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| workspace-id | Workspace | Select | Select a workspace to search within. | Yes | - |
| query | Search Query | String | The search query to find matching tasks. | Yes | - |

## Returns

This action returns an object containing an array of matching task objects:

| Field | Type | Description |
|-------|------|-------------|
| tasks | Array | List of task objects that match the search query |

Each task object contains:

| Field | Type | Description |
|-------|------|-------------|
| id | String | The unique identifier of the task |
| name | String | The name/title of the task |
| status | Object | The current status of the task including color |
| list | Object | Information about the list containing the task |
| folder | Object | Information about the folder containing the list |
| space | Object | Information about the space containing the folder |

## Example Usage

You can use this action to:

1. Find tasks matching specific keywords or criteria
2. Search for tasks across different lists and folders
3. Create a custom search interface for your ClickUp workspace
4. Set up automated actions based on finding specific tasks
5. Generate reports for tasks matching certain patterns

## Notes

- The search query can include task names, descriptions, comments, and custom fields.
- Search results may vary based on user permissions.
- For more targeted results, use specific keywords in your query.
- ClickUp's search functionality supports advanced operators - refer to ClickUp documentation for details.