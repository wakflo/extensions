# List Issues

## Description
List issues from a Jira project.

## Properties
| Name             | Type    | Required | Description                                        |
|------------------|---------|----------|----------------------------------------------------|
| projectId        | string  | Yes      | The project to list issues from                    |
| maxResults       | number  | No       | Maximum number of results to return (default: 50)  |
| orderBy          | string  | No       | How to order the results                          |
| onlyAssignedToMe | boolean | No       | Only show issues assigned to the current user     |

## Details
- **Type**: sdkcore.ActionTypeNormal