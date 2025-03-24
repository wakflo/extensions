# Task Updated

## Description

Triggers a workflow whenever a task is updated in a specified Asana workspace or project.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Project ID | String | No | Optional: The GID of a specific project to monitor for updated tasks |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Usage

This trigger polls the Asana API at regular intervals to detect tasks that have been updated since the last check. You must specify a Workspace ID, and you can optionally filter to only detect tasks updated in a specific project.

## Example Payload

```json
{
  "data": [
    {
      "gid": "1234567890",
      "name": "Updated Task",
      "resource_type": "task",
      "created_at": "2023-01-10T08:00:00.000Z",
      "modified_at": "2023-01-15T09:30:00.000Z",
      "notes": "This task has been updated with new information",
      "completed": false
    }
  ]
}