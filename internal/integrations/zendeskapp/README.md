# Zendesk Customer Support Integration

## Description

Integrate your Zendesk Customer Support platform with our workflow automation software to streamline support processes and enhance customer service. Automate ticket management, monitor support activities, and improve response times. Use this integration to:

- Automatically create and update support tickets based on external events
- Track ticket status changes and alert team members
- Add comments to tickets as part of automated workflows
- Retrieve and search ticket information to inform decision making
- Enhance customer satisfaction through faster, more consistent support responses
- Generate custom reports on support performance metrics

**Zendesk Customer Support Integration Documentation**

**Overview**
The Zendesk Customer Support integration allows you to connect your Zendesk account with our workflow automation software, enabling you to automate ticket management tasks and enhance your customer support workflow.

**Prerequisites**

- A Zendesk account with administrator access
- Our workflow automation software account
- Your Zendesk subdomain (e.g., company.zendesk.com)
- API token or OAuth2 credentials for authentication

**Setup Instructions**

1. Log in to your Zendesk account and navigate to the Admin Center.
2. Go to Apps and Integrations > APIs > Zendesk API.
3. Generate an API token or set up OAuth2 credentials.
4. In our workflow automation software, go to the "Integrations" section and click on "Zendesk Customer Support".
5. Enter your Zendesk subdomain and API credentials, then click "Connect".

**Available Actions**

- **Create Ticket**: Create a new support ticket in Zendesk with specified details.
- **Update Ticket**: Update an existing ticket's properties such as status, priority, or assignee.
- **Add Comment**: Add a public or private comment to an existing ticket.
- **Get Ticket**: Retrieve detailed information about a specific ticket.
- **List Tickets**: Retrieve a list of tickets matching specified criteria.
- **Search Tickets**: Search for tickets using Zendesk's query language.

**Available Triggers**

- **Ticket Created**: Trigger a workflow when a new ticket is created in Zendesk.
- **Ticket Updated**: Trigger a workflow when any field of a ticket is updated.
- **Ticket Status Changed**: Trigger a workflow when a ticket's status changes.

**Example Use Cases**

1. When a high-priority ticket is created, automatically notify the support team manager.
2. Create follow-up tasks in your project management tool when tickets remain unresolved for more than 24 hours.
3. Update your CRM system when a customer submits a new ticket.
4. Send customers an SMS notification when their ticket status changes to "Solved".

**Troubleshooting Tips**

- Ensure that your API credentials have the necessary permissions for the actions you want to perform.
- Check that your Zendesk subdomain is correctly specified without the "https://" prefix.
- Verify that your Zendesk account is active and in good standing.

**FAQs**

Q: Can I use this integration with Zendesk Support Professional plan or higher?
A: Yes, this integration works with Zendesk Support Professional, Enterprise, and Enterprise Plus plans.

Q: How often does the integration check for new tickets or updates?
A: For trigger-based workflows, the integration polls Zendesk every 5 minutes by default.

## Categories

- customer-support
- helpdesk

## Authors

- Wakflo

## Actions

| Name           | Description                                                                                             | Link                              |
| -------------- | ------------------------------------------------------------------------------------------------------- | --------------------------------- |
| Create Ticket  | Creates a new ticket in your Zendesk support system with customizable fields.                           | [docs](actions/create_ticket.md)  |
| Update Ticket  | Updates an existing ticket in Zendesk with new information.                                             | [docs](actions/update_ticket.md)  |
| Get Ticket     | Retrieves detailed information about a specific ticket in your Zendesk support system.                  | [docs](actions/get_ticket.md)     |
| List Tickets   | Retrieves a list of tickets from your Zendesk account based on specified filters.                       | [docs](actions/list_tickets.md)   |
| Add Comment    | Adds a comment to an existing Zendesk ticket, with options for public or private (internal) visibility. | [docs](actions/add_comment.md)    |
| Search Tickets | Searches for tickets in your Zendesk account using Zendesk's query language.                            | [docs](actions/search_tickets.md) |

## Triggers

| Name                  | Description                                                                                                                | Link                                      |
| --------------------- | -------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------- |
| Ticket Created        | Triggers a workflow when a new support ticket is created in your Zendesk account.                                          | [docs](triggers/ticket_created.md)        |
| Ticket Updated        | Triggers a workflow when any fields of a ticket are updated in your Zendesk account.                                       | [docs](triggers/ticket_updated.md)        |
| Ticket Status Changed | Triggers a workflow when a ticket's status is changed, allowing you to automate follow-up actions based on the new status. | [docs](triggers/ticket_status_changed.md) |
