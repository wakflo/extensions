# Asana Integration

## Description

Integrate Asana project management with your workflow automation to streamline task tracking and team collaboration. Connect your Asana workspace to automatically sync tasks, projects, and updates with other business tools and processes. Use this integration to:

* Create and update tasks based on external events
* Automate task assignments and due date tracking
* Trigger notifications when tasks are created or completed
* Sync project data across platforms and departments
* Build custom project management workflows tailored to your team's needs

**Asana Integration Documentation**

**Overview**
The Asana integration allows you to seamlessly connect your Asana workspace with our workflow automation software, enabling you to automate tasks and streamline your project management workflow.

**Prerequisites**

* An Asana account
* Our workflow automation software account
* An Asana Personal Access Token (found in your Asana account settings)

**Setup Instructions**

1. Log in to your Asana account and navigate to "My Profile Settings"
2. Click on "Apps" and then "Manage Developer Apps"
3. Click "Create New Personal Access Token" and provide a description
4. Copy the generated token
5. In our workflow automation software, go to the "Integrations" section and click on "Asana"
6. Paste the Personal Access Token and click "Connect"

**Available Actions**

* **Get Task**: Retrieve details of a specific Asana task
* **List Tasks**: Get a list of tasks from a project or workspace
* **Create Task**: Create a new task in Asana
* **Update Task**: Modify an existing task's details
* **Get Project**: Retrieve details of a specific project
* **List Projects**: Get a list of projects from a workspace

**Available Triggers**

* **Task Created**: Trigger a workflow when a new task is created
* **Task Updated**: Trigger a workflow when a task is updated
* **Task Completed**: Trigger a workflow when a task is marked as completed

**Example Use Cases**

1. **Client Onboarding**: When a new client is added to your CRM, automatically create an Asana project with predefined tasks for the onboarding process.
2. **Bug Tracking**: Create Asana tasks automatically when bugs are reported in your support system.
3. **Project Updates**: Send notifications to team messaging platforms when important tasks are completed in Asana.

**Troubleshooting Tips**

* Ensure your Personal Access Token has the necessary permissions
* Check that your workspace and project GIDs are correct
* Verify that you're using the correct workspace/project when creating tasks

**FAQs**

Q: Can I use this integration with Asana Teams?
A: Yes, you can use this integration with any workspace or team that your Asana account has access to.

Q: Is there a limit to how many tasks I can create?
A: The limits are based on your Asana plan. Refer to Asana's API documentation for specific rate limits.

## Categories

- project-management
- productivity

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name | Description | Link |
|------|-------------|------|
| Get Task | Retrieves detailed information for a specific Asana task by its unique identifier. | [docs](actions/get_task.md) |
| List Tasks | Retrieves a list of tasks from a specific project or workspace based on provided filters. | [docs](actions/list_tasks.md) |
| Create Task | Creates a new task in Asana with the specified details, including name, description, assignee, and due date. | [docs](actions/create_task.md) |
| Update Task | Updates an existing Asana task with new information such as status, assignee, due date, or description. | [docs](actions/update_task.md) |


## Triggers

| Name | Description | Link |
|------|-------------|------|
| Task Created | Triggers a workflow whenever a new task is created in a specified Asana workspace or project. | [docs](triggers/task_created.md) |
| Task Updated | Triggers a workflow whenever a task is updated in a specified Asana workspace or project. | [docs](triggers/task_updated.md) |
| Task Completed | Triggers a workflow whenever a task is marked as completed in a specified Asana workspace or project. | [docs](triggers/task_completed.md) |