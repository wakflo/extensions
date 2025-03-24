# Task Created

## Description

Triggers a workflow whenever a new task is created in a specified Asana workspace or project.

## Properties

| Name | Type | Required | Description |
|------|------|----------|-------------|
| Project ID | String | No | Optional: The GID of a specific project to monitor for new tasks |

## Details

- **Type**: sdkcore.TriggerTypePolling

## Usage

This trigger polls the Asana API at regular intervals to detect new tasks that have been created since the last check. You must specify a Workspace ID, and you can optionally filter to only detect tasks created in a specific project.

## Example Payload

```json
{
  "data": [
    {
      "gid": "1234567890",
      "name": "New Task",
      "resource_type": "task",
      "created_at": "2023-01-15T08:00:00.000Z",
      "notes": "This is a new task that was just created",
      "completed": false
    }
  ]
}