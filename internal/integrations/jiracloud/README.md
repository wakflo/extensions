# Jira Integration

## Description

Integrate Jira with your workflow automation to streamline project management, issue tracking, and team collaboration. This integration allows you to automatically create, update, and track issues, synchronize data between systems, and keep your team informed of important changes. Use this integration to:

* Create and update Jira issues automatically
* Track issue status changes and transitions
* Sync project data between Jira and other tools
* Notify team members of important issue updates
* Automate issue assignment and prioritization
* Generate custom reports based on Jira data

**Jira Integration Documentation**

**Overview**
The Jira integration allows you to seamlessly connect your Jira account with our workflow automation software, enabling you to automate tasks and streamline your workflow.

**Prerequisites**

* A Jira account (Cloud or Server)
* Our workflow automation software account
* Jira API token or Basic Authentication credentials
* Permission to access the Jira API

**Setup Instructions**

1. Log in to your Jira account and navigate to Account Settings
2. Under Security, select "Create and manage API tokens"
3. Click "Create API token" and provide a label for the token
4. Copy the generated API token
5. In our workflow automation software, navigate to the "Integrations" section and select "Jira"
6. Enter your Jira email address, API token, and Jira instance URL (e.g., https://yourcompany.atlassian.net)
7. Test the connection and save

**Available Actions**

* **Create Issue**: Create a new issue in Jira with specified details
* **Get Issue**: Retrieve detailed information about a specific issue
* **List Issues**: Retrieve a list of issues based on specified criteria
* **Update Issue**: Update the details of an existing issue
* **Transition Issue**: Move an issue through its workflow states
* **Add Comment**: Add a comment to an existing issue

**Available Triggers**

* **Issue Created**: Trigger a workflow when a new issue is created
* **Issue Updated**: Trigger a workflow when an issue is updated

**Example Use Cases**

1. **Customer Support Automation**: Automatically create Jira issues when new support tickets are received
2. **Development Workflow**: Trigger deployments when issues are transitioned to "Ready for Deployment"
3. **Project Status Reporting**: Generate weekly reports based on Jira issue status
4. **Cross-Platform Synchronization**: Keep Jira issues synchronized with tasks in other tools

**Troubleshooting Tips**

* Ensure your API token has the appropriate permissions
* Check that your Jira instance URL is correctly formatted
* Verify that the required fields for issue creation are properly specified
* For Jira Server instances, ensure the API endpoint is accessible from our service

**FAQs**

Q: Can I use this integration with both Jira Cloud and Jira Server?
A: Yes, this integration works with both Jira Cloud and Server editions, but the authentication method may differ.

Q: How frequently are the triggers checked?
A: The polling triggers (Issue Created, Issue Updated) check for new events every 5 minutes.

## Categories

- project-management
- collaboration

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name             | Description                                                              | Link                                  |
|------------------|--------------------------------------------------------------------------|---------------------------------------|
| Create Issue     | Create a new issue in Jira with specified details                        | [docs](actions/create_issue.md)       |
| Get Issue        | Retrieve detailed information about a specific issue                     | [docs](actions/get_issue.md)          |
| List Issues      | Retrieve a list of issues based on specified criteria                    | [docs](actions/list_issues.md)        |
| Update Issue     | Update the details of an existing issue                                  | [docs](actions/update_issue.md)       |
| Transition Issue | Move an issue through its workflow states                                | [docs](actions/transition_issue.md)   |
| Add Comment      | Add a comment to an existing issue                                       | [docs](actions/add_comment.md)        |

## Triggers

| Name           | Description                                               | Link                               |
|----------------|-----------------------------------------------------------|----------------------------------- |
| Issue Created  | Trigger a workflow when a new issue is created in Jira    | [docs](triggers/issue_created.md)  |
| Issue Updated  | Trigger a workflow when an issue is updated in Jira       | [docs](triggers/issue_updated.md)  |