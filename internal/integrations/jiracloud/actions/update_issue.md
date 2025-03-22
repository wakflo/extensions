# Update Issue

## Description
Update an existing issue in Jira.

## Properties
| Name        | Type   | Required | Description                                     |
|-------------|--------|----------|-------------------------------------------------|
| projectId   | string | Yes      | The project the issue belongs to                |
| issueId     | string | Yes      | The ID of the issue to update                   |
| summary     | string | No       | Updated summary of the issue                    |
| description | string | No       | Updated description of the issue                |
| priority    | string | No       | Updated priority of the issue                   |
| assigneeId  | string | No       | ID of the user to assign the issue to           |
| ParentKey   | string | No       | If this issue is a subtask, insert the parent issue key |

## Details
- **Type**: sdkcore.ActionTypeNormal