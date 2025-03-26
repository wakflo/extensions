# Create Issue

## Description
Create a new issue in Jira with specified details.

## Properties
| Name        | Type   | Required | Description                                     |
|-------------|--------|----------|-------------------------------------------------|
| projectId   | string | Yes      | The project to create the issue in              |
| issueType   | string | Yes      | The type of issue (e.g., Bug, Story, Task)      |
| summary     | string | Yes      | The summary or title of the issue               |
| description | string | No       | The detailed description of the issue           |
| priority    | string | No       | The priority of the issue                       |
| assigneeId  | string | No       | ID of the user to assign the issue to           |
| labels      | string | No       | Comma-separated list of labels to add to issue  |

## Details
- **Type**: sdkcore.ActionTypeNormal