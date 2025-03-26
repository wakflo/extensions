# Wrike Integration

## Description

Integrate your Wrike project management data with our workflow automation software to streamline task management and team collaboration processes. Automatically sync Wrike tasks, projects, and updates with your workflows, eliminating manual data entry and reducing errors. Use this integration to:

* Track task progress and deadlines in real-time
* Automate task assignments and updates based on Wrike project status
* Generate accurate reports using synced Wrike data
* Enhance collaboration by sharing Wrike task information with team members
* Gain insights into project performance and resource allocation

**Wrike Integration Documentation**

**Overview**
The Wrike integration allows you to seamlessly connect your Wrike account with our workflow automation software, enabling you to automate tasks and streamline your project management workflow.

**Prerequisites**

* A Wrike account
* Our workflow automation software account
* A Wrike API access token (found in your Wrike account settings)

**Setup Instructions**

1. Log in to your Wrike account and navigate to the "Account Settings" > "Apps & Integrations" section.
2. Click on "API Apps" tab and create a new app or use an existing one.
3. Generate an access token and copy it.
4. In our workflow automation software, go to the "Integrations" section and click on "Wrike".
5. Enter the API token generated in step 3 and click "Connect".

**Available Actions**

* **Get Task**: Retrieve detailed information about a specific task in Wrike by its ID.
* **List Tasks**: Retrieve a list of tasks from Wrike with optional filtering parameters.
* **Create Task**: Automatically create a new task in Wrike with specified details.
* **Update Task**: Update the properties of an existing task in Wrike.

**Available Triggers**

* **Task Created**: Trigger a workflow when a new task is created in Wrike.
* **Task Updated**: Trigger a workflow when a task is updated in Wrike.

**Example Use Cases**

1. When a new lead is added to your CRM, automatically create a task in Wrike for the sales team.
2. When a task status is changed to "Completed" in Wrike, automatically send a notification to stakeholders.
3. Sync task updates between Wrike and other project management tools.

**Troubleshooting Tips**

* Ensure that your API token is valid and has not expired.
* Check that you have the necessary permissions in Wrike to perform the requested actions.
* Verify that the task IDs provided in your workflows are correct and exist in your Wrike account.

**FAQs**

Q: Can I use this integration to create custom fields in Wrike tasks?
A: Yes, you can include custom fields when creating or updating tasks through the integration.

Q: How frequently does the Task Updated trigger check for updates?
A: The integration polls for task updates every 5 minutes by default.

## Categories

- productivity

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name         | Description                                                                 | Link                         |
|--------------|-----------------------------------------------------------------------------|------------------------------|
| Get Task     | Retrieves detailed information about a specific task in Wrike by its ID.    | [docs](actions/get_task.md)  |
| List Tasks   | Retrieve a list of tasks from Wrike with optional filtering parameters.     | [docs](actions/list_tasks.md)|
| Create Task  | Create a new task in Wrike with specified title, description, and status.   | [docs](actions/create_task.md)|
| Update Task  | Update an existing task in Wrike with new properties such as status.        | [docs](actions/update_task.md)|

## Triggers

| Name          | Description                                                             | Link                            |
|---------------|-------------------------------------------------------------------------|----------------------------------|
| Task Created  | Triggers when a new task is created in your Wrike account.              | [docs](triggers/task_created.md) |
| Task Updated  | Triggers when a task is updated in your Wrike account.                  | [docs](triggers/task_updated.md) |