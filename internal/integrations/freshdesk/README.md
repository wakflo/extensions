# Freshdesk Integration

## Description

Integrate your Freshdesk customer support platform with our workflow automation software to streamline support processes and enhance customer service capabilities. Connect Freshdesk tickets, contacts, and support activities with your workflows to provide seamless customer experiences. Use this integration to:

* Automatically create and update support tickets based on external triggers
* Sync customer information between Freshdesk and other systems
* Trigger notifications when ticket statuses change
* Generate reports on support activities and performance
* Track and prioritize urgent support issues

**Freshdesk Integration Documentation**

**Overview**
The Freshdesk integration allows you to seamlessly connect your Freshdesk account with our workflow automation software, enabling you to automate your customer support processes.

**Prerequisites**

* A Freshdesk account
* Our workflow automation software account
* A Freshdesk API key (found in your Freshdesk profile settings)
* Your Freshdesk domain (e.g., yourcompany.freshdesk.com)

**Setup Instructions**

1. Log in to your Freshdesk account and go to your profile settings.
2. Under the API section, copy your API key.
3. In our workflow automation software, go to the "Integrations" section and click on "Freshdesk".
4. Enter your Freshdesk domain and API key to connect.

**Available Actions**

* **Create Ticket**: Create a new support ticket in Freshdesk.
* **Get Ticket**: Retrieve details of a specific ticket by ID.
* **Update Ticket**: Update the properties of an existing ticket.
* **List Tickets**: Retrieve a list of tickets based on filter criteria.
* **Search Tickets**: Search for tickets using various parameters.
* **Add Note**: Add a private or public note to an existing ticket.

**Available Triggers**

* **Ticket Created**: Trigger a workflow when a new ticket is created in Freshdesk.
* **Ticket Updated**: Trigger a workflow when a ticket is updated in Freshdesk.

**Example Use Cases**

1. **Automated Ticket Creation**: When a form is submitted on your website, automatically create a support ticket in Freshdesk.
2. **Cross-platform Notification**: When a high-priority ticket is created, send notifications to team messaging platforms.
3. **SLA Tracking**: When a ticket approaches its SLA deadline, trigger escalation workflows.
4. **Customer Communication**: When a ticket status changes, automatically notify the customer via email or SMS.

**Troubleshooting Tips**

* Ensure your API key has the necessary permissions in Freshdesk.
* Check that your Freshdesk domain is entered correctly without 'https://' prefix.
* Verify your Freshdesk plan supports the API features you're trying to use.

## Categories

- customer-support
- helpdesk

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name | Description | Link |
|------|-------------|------|
| Create Ticket | Create a new support ticket in the Freshdesk system with customizable fields. | [docs](actions/create_ticket.md) |
| Get Ticket | Retrieve detailed information about a specific ticket by its ID. | [docs](actions/get_ticket.md) |
| Update Ticket | Update the properties and fields of an existing Freshdesk ticket. | [docs](actions/update_ticket.md) |
| List Tickets | Retrieve a list of tickets based on filter criteria. | [docs](actions/list_tickets.md) |
| Search Tickets | Search for tickets using various parameters like keywords, statuses, and priorities. | [docs](actions/search_tickets.md) |
| Add Note | Add a private or public note to an existing ticket for internal communication or customer updates. | [docs](actions/add_note.md) |

## Triggers

| Name | Description | Link |
|------|-------------|------|
| Ticket Created | Trigger a workflow when a new ticket is created in Freshdesk. | [docs](triggers/ticket_created.md) |
| Ticket Updated | Trigger a workflow when a ticket is updated in Freshdesk, including status changes, priority updates, or note additions. | [docs](triggers/ticket_updated.md) |
