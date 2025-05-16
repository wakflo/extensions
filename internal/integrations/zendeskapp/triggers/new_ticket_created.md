# New Ticket Added

This trigger activates when a new support ticket is created in your Zendesk account.

## Overview

The New Ticket Added trigger monitors your Zendesk account for newly created tickets and initiates workflows when they are detected. This enables you to automate responses, route tickets to the appropriate teams, create tasks in other systems, or perform any other action based on incoming support requests.

## Use Cases

- Send automated acknowledgments to customers when they submit tickets
- Create tasks in project management tools based on support requests
- Alert team members about high-priority tickets that require immediate attention
- Log ticket information to analytics platforms for reporting
- Trigger approval workflows for specific types of customer requests

## Authentication

This trigger requires Zendesk API authentication with:

- Email: Your Zendesk administrator email
- API Token: A valid Zendesk API token
- Subdomain: Your Zendesk account subdomain (the part before .zendesk.com)

## Configuration

| Field | Description                          | Required |
| ----- | ------------------------------------ | -------- |
| ID    | Optional identifier for this trigger | No       |

## Output

The trigger returns all new tickets created since the last execution. Each ticket includes:

- Ticket ID
- Subject
- Description
- Priority
- Status
- Requester information
- Created date/time
- Tags and custom fields

## Notes

- The trigger uses polling to check for new tickets at regular intervals
- Only tickets created after the last execution will be returned
- You can use filters in subsequent workflow steps to focus on specific ticket types
