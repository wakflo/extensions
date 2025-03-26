# Add Comment

## Description
Add a comment to a Jira issue.

## Properties
| Name        | Type   | Required | Description                          |
|-------------|--------|----------|--------------------------------------|
| projectId   | string | Yes      | The project the issue belongs to     |
| issueId     | string | Yes      | The ID of the issue to comment on    |
| commentText | string | Yes      | Text of the comment to add           |

## Details
- **Type**: sdkcore.ActionTypeNormal