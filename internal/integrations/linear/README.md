# Linear Integration

## Description

Linear Integration:

The Linear integration feature in our workflow automation software enables seamless connection with other systems and applications by mapping data between them. This allows for efficient exchange of information, streamlining processes, and enhancing overall productivity. With linear integration, you can automate complex workflows, eliminate manual errors, and gain real-time visibility into your operations.

**Linear Integration with [Workflow Automation Software]**

**Overview**
The Linear integration with [Workflow Automation Software] enables seamless communication between your Linear account and our platform, allowing you to automate workflows and streamline processes.

**Prerequisites**

* A valid Linear account
* A [Workflow Automation Software] account with the necessary permissions

**Integration Steps**

1. **Create a new integration**: Log in to your [Workflow Automation Software] account and navigate to the Integrations section.
2. **Search for Linear**: In the search bar, type "Linear" and select the corresponding result.
3. **Configure the integration**: Fill in the required fields, including your Linear API key and secret key. You can obtain these keys from your Linear account settings.
4. **Authorize the integration**: Authorize [Workflow Automation Software] to access your Linear account by clicking the "Authorize" button.

**Available Endpoints**

The Linear integration with [Workflow Automation Software] supports the following endpoints:

* `GET /orders`: Retrieve a list of orders
* `GET /orders/{id}`: Retrieve a specific order
* `POST /orders`: Create a new order
* `PUT /orders/{id}`: Update an existing order
* `DELETE /orders/{id}`: Delete an existing order

**Example Use Cases**

1. **Automate order creation**: When a new order is placed in Linear, use [Workflow Automation Software] to create a corresponding task or workflow.
2. **Sync orders with Linear**: Use the integration to keep your Linear account up-to-date with the latest order information from [Workflow Automation Software].

**Troubleshooting**

* If you encounter issues during the integration process, refer to our troubleshooting guide for common errors and solutions.
* For more complex issues, contact our support team or reach out to Linear's support team.

**FAQs**

1. **What is the maximum number of orders that can be retrieved at once?**: The maximum number of orders that can be retrieved at once is 100.
2. **Can I use this integration for both Linear and other e-commerce platforms?**: Yes, you can use [Workflow Automation Software] to integrate with multiple e-commerce platforms, including Linear.

**Release Notes**

* Version 1.0: Initial release of the Linear integration
* Version 1.1: Added support for updating order status

By integrating your Linear account with [Workflow Automation Software], you can streamline your workflows and automate repetitive tasks, freeing up more time to focus on growing your business.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>


## Actions

| Name         | Description                                                                                                                                                                                                    | Link                            |
|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------|
| Create Issue | Create Issue: Automatically generates a new issue in your project management tool (e.g., Jira, Trello) based on specific conditions or triggers, ensuring timely and organized tracking of tasks and projects. | [docs](actions/create_issue.md) |## Actions
| Find Issues  | Automatically identifies and extracts issues from various data sources, such as logs, tickets, or databases, to provide a centralized view of problems and errors in your workflow.                            | [docs](actions/find_issues.md)  |
| Update Issue | Updates an existing issue in your project management tool with the latest information from other connected systems, ensuring seamless data synchronization and reducing manual errors.                         | [docs](actions/update_issue.md) |

## Triggers

| Name          | Description                                                                                                                                                                                                                                                    | Link                              |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------|
| Issue Created | Triggered when a new issue is created in your project management tool, this integration allows you to automate workflows and tasks immediately after an issue is reported, streamlining your team's response time and ensuring prompt attention to new issues. | [docs](triggers/issue_created.md) |## Triggers
| ------        | -------------                                                                                                                                                                                                                                                  | ------                            |
| Issue Updated | Triggered when an issue is updated in your project management tool, such as Jira or Trello. This integration allows you to automate workflows and tasks based on changes made to issues, including status updates, comments, and attachments.                  | [docs](triggers/issue_updated.md) |