# Issue Created

## Description
Trigger a workflow when a new issue is created in Jira.

## Properties
| Name      | Type   | Required | Description                            |
|-----------|--------|----------|----------------------------------------|
| projectId | string | No       | Filter for issues in a specific project |
| issueType | string | No       | Filter for specific issue types        |

## Details
- **Type**: sdkcore.TriggerTypePolling
- **Polling Interval**: 5 minutes

## Trigger Output
The trigger returns newly created issues matching the specified filters:

```json
{
  "issues": [
    {
      "id": "12345",
      "key": "PRJ-123",
      "self": "https://yourcompany.atlassian.net/rest/api/3/issue/12345",
      "fields": {
        "summary": "New issue created",
        "description": {
          "type": "doc",
          "content": [...]
        },
        "status": {
          "name": "To Do"
        },
        "creator": {
          "displayName": "John Doe",
          "emailAddress": "john.doe@example.com"
        },
        "created": "2023-05-05T12:34:56.789Z",
        "priority": {
          "name": "Medium"
        },
        "assignee": {
          "displayName": "Jane Smith",
          "emailAddress": "jane.smith@example.com"
        }
      }
    }
  ],
  "total": 1,
  "maxResults": 50,
  "startAt": 0
}
```