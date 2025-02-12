# Todoist Integration

## Description

Integrate your Todoist tasks with our workflow automation software to streamline your productivity and simplify task management. Automatically create new workflows or update existing ones based on Todoist task changes, such as due dates, priorities, or labels. Use Zapier's Todoist integration to:

* Create new workflows for each Todoist project
* Update workflow stages based on Todoist task status (e.g., completed, in progress)
* Trigger custom actions when a Todoist task is moved to a specific stage
* Set deadlines and reminders from Todoist tasks directly into your workflows

Elevate your productivity by automating repetitive tasks, reducing manual errors, and gaining real-time visibility into your workflow's progress. Seamlessly integrate Todoist with our workflow automation software to achieve more in less time.

**Todoist Integration Documentation**

**Overview**
The Todoist integration allows you to seamlessly connect your Todoist account with our workflow automation software, enabling you to automate tasks and streamline your productivity.

**Prerequisites**

* A Todoist account
* Our workflow automation software account
* The Todoist API token (available in the Todoist settings)

**Setup Instructions**

1. Log in to your Todoist account and navigate to the Settings > Integrations page.
2. Click on the "Create API Token" button and follow the prompts to generate a new token.
3. In our workflow automation software, go to the Integrations page and click on the "Add Integration" button.
4. Select Todoist from the list of available integrations.
5. Enter your Todoist API token and click "Connect".
6. Authorize the integration by clicking the "Authorize" button.

**Available Actions**

* Create a new Todoist task: Use this action to create a new task in Todoist with a specified title, description, and due date.
* Update an existing Todoist task: Use this action to update an existing task in Todoist with a specified title, description, and due date.
* Mark a Todoist task as completed: Use this action to mark a Todoist task as completed.

**Trigger Options**

* When a new task is created: Trigger the integration when a new task is created in Todoist.
* When a task is updated: Trigger the integration when an existing task is updated in Todoist.
* When a task is marked as completed: Trigger the integration when a task is marked as completed in Todoist.

**Example Use Cases**

* Automatically create a new Todoist task when a specific event occurs in your workflow automation software.
* Update the due date of a Todoist task based on the status of a related workflow automation software task.
* Mark a Todoist task as completed when a specific condition is met in your workflow automation software.

**Troubleshooting Tips**

* Ensure that you have entered the correct API token and authorized the integration.
* Check that the Todoist account associated with the integration has the necessary permissions to create, update, or mark tasks as completed.
* Review the Todoist API documentation for any rate limits or usage guidelines that may affect your integration.

**FAQs**

Q: What is the maximum number of tasks that can be created or updated in a single integration trigger?
A: The maximum number of tasks that can be created or updated in a single integration trigger is 10. If you need to create or update more than 10 tasks, you will need to use multiple triggers.

Q: Can I use the Todoist integration with multiple workflows or projects?
A: Yes, you can use the Todoist integration with multiple workflows or projects by creating separate integrations for each one.

By following these instructions and using the available actions and trigger options, you can automate your workflow and streamline your productivity with our Todoist integration.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>


## Actions
| Name                       | Description                                                                                                                                                                                                  | Link                                          |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------|
| Create Project             | Create Project: Initiates the creation of a new project in your project management system, allowing you to start tracking tasks, milestones, and team progress from within our workflow automation platform. | [docs](actions/create_project.md)             |## Actions
| Create Task                | Create Task: Automatically generates and assigns a new task to a team member or group, allowing you to streamline workflows and ensure timely completion of tasks.                                           | [docs](actions/create_task.md)                |
| Get Active Task            | Retrieves the currently active task in the workflow, allowing you to access and manipulate its properties or trigger subsequent actions based on its status.                                                 | [docs](actions/get_active_task.md)            |
| Get Project                | Retrieves project details from the specified project management system or platform.                                                                                                                          | [docs](actions/get_project.md)                |
| List Projects              | Retrieves a list of all projects within your organization, allowing you to easily manage and track multiple projects from a single location.                                                                 | [docs](actions/list_projects.md)              |
| List Project Collaborators | Retrieves a list of users who are currently collaborating on a specific project, including their roles and permissions.                                                                                      | [docs](actions/list_project_collaborators.md) |
| List Task                  | The "List Tasks" integration action retrieves and displays a list of tasks from a specified workflow or project, allowing you to easily view and manage your task portfolio within the automation workflow.  | [docs](actions/list_task.md)                  |
| Updste Project             | Updates project information in the designated project management system, ensuring accurate and up-to-date records of project details, milestones, and tasks.                                                 | [docs](actions/update_project)             |
| Update Task                | Updates the status and details of an existing task in your workflow, allowing you to reflect changes or new information without having to recreate the task.                                                 | [docs](actions/update_task.md)                |## Triggers

## Triggers
| Name           | Description                                                                                                                          | Link                               |
|----------------|--------------------------------------------------------------------------------------------------------------------------------------|------------------------------------|
| Task Completed | Triggered when a task is marked as completed in the workflow, allowing subsequent actions to be executed based on the task's status. | [docs](triggers/task_completed.md) |