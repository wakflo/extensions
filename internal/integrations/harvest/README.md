# Harvest Integration

## Description

Integrate your Harvest time tracking data with our workflow automation software to streamline project management and accounting processes. Automatically sync Harvest projects, tasks, and timesheets with your workflows, eliminating manual data entry and reducing errors. Use this integration to:

* Track project hours and costs in real-time
* Automate task assignments and updates based on Harvest project status
* Generate accurate invoices and reports using synced Harvest data
* Enhance collaboration by sharing Harvest project information with team members
* Gain insights into project performance and resource allocation

**Harvest Integration Documentation**

**Overview**
The Harvest integration allows you to seamlessly connect your Harvest account with our workflow automation software, enabling you to automate tasks and streamline your workflow.

**Prerequisites**

* A Harvest account
* Our workflow automation software account
* The latest version of the Harvest API token (found in your Harvest settings)

**Setup Instructions**

1. Log in to your Harvest account and navigate to the "Settings" tab.
2. Click on "API Tokens" and generate a new token or select an existing one.
3. In our workflow automation software, go to the "Integrations" section and click on "Harvest".
4. Enter the API token generated in step 2 and click "Connect".

**Available Actions**

* **Create Harvest Task**: Automatically create a new task in Harvest based on specific conditions or triggers.
* **Update Harvest Task Status**: Update the status of an existing task in Harvest based on specific conditions or triggers.
* **Get Harvest Tasks**: Retrieve a list of tasks from Harvest and use them to trigger automation workflows.

**Example Use Cases**

1. Automate task creation: When a new lead is added to your CRM, create a corresponding task in Harvest with the due date set to the next business day.
2. Update task status: When a task is marked as completed in your workflow automation software, update the status of the corresponding task in Harvest to "Completed".

**Troubleshooting Tips**

* Ensure that you have entered the correct API token and that it has not expired.
* Check the Harvest API documentation for any rate limits or usage guidelines.

**FAQs**

Q: What is the maximum number of tasks I can retrieve from Harvest at once?
A: The maximum number of tasks you can retrieve from Harvest at once is 100. If you need to retrieve more, you will need to use pagination or split your request into multiple calls.

Q: Can I use this integration to automate tasks in Harvest that are not related to my workflow automation software?
A: Yes, you can use this integration to automate any tasks in Harvest that meet the conditions or triggers set up in your workflow automation software.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>


## Actions

| Name          | Description                                                                                                                                             | Link                             |
|---------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------|
| Get Invoice   | Retrieves an invoice from the accounting system, allowing you to automate tasks that require access to invoice data.                                    | [docs](actions/get_invoice.md)   |## Actions
| List Invoices | Retrieve and list all invoices associated with a specific account or organization, allowing you to easily track and manage your financial transactions. | [docs](actions/list_invoices.md) |## Triggers

## Triggers

| Name            | Description                                                                                                                                                                                                                                                                                                                                                                                                                                  | Link                                |
|-----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| Invoice Updated | The "Invoice Updated" integration trigger is designed to automate workflows when an invoice is updated in your accounting or ERP system. This trigger can be used to initiate a series of automated tasks, such as sending notifications to stakeholders, updating CRM records, or triggering payment processing. When an invoice is updated, this trigger will fire and allow you to define the subsequent actions that need to take place. | [docs](triggers/invoice_updated.md) |