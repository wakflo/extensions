# Transition Issue

## Description
Move a Jira issue to a different status.

## Properties
| Name         | Type   | Required | Description                                    |
|--------------|--------|----------|------------------------------------------------|
| projectId    | string | Yes      | The project the issue belongs to               |
| issueId      | string | Yes      | The ID of the issue to transition              |
| transitionId | string | Yes      | The ID of the transition to apply              |
| comment      | string | No       | Comment to add when transitioning the issue    |

## Details
- **Type**: sdkcore.ActionTypeNormal