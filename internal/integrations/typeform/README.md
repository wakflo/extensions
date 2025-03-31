# Typeform Integration

## Description

Integrate Typeform with your workflow automation to create dynamic forms and process responses efficiently. Connect your forms to other systems, automate follow-up actions, and analyze data based on form submissions. Use this integration to:

- Automatically process form submissions
- Create personalized follow-ups based on responses
- Sync form data with CRM, marketing, or customer support tools
- Analyze form responses and generate custom reports
- Trigger workflows when forms are submitted or updated

**Typeform Integration Documentation**

**Overview**
The Typeform integration allows you to connect your Typeform account with our workflow automation software, enabling you to automate tasks based on form submissions and manage your forms programmatically.

**Prerequisites**

- A Typeform account
- Our workflow automation software account
- A Typeform Personal Access Token (found in your Typeform account settings)

**Setup Instructions**

1. Log in to your Typeform account and navigate to the "Settings" > "Personal tokens" section.
2. Create a new Personal Access Token with the required scopes (forms:read, responses:read).
3. In our workflow automation software, go to the "Integrations" section and click on "Typeform".
4. Enter the Personal Access Token generated in step 2 and click "Connect".

**Available Actions**

- **Get Form**: Retrieve details about a specific form.
- **List Forms**: Get a list of all forms in your Typeform account.
- **Get Responses**: Retrieve responses for a specific form.
- **Create Form**: Create a new form in your Typeform account.

**Available Triggers**

- **Form Response**: Triggers when a new response is submitted to a form.

**Example Use Cases**

1. Send a personalized email to new form respondents when they submit your contact form.
2. Create a task in your project management tool when customers submit a request form.
3. Add new leads to your CRM when they submit a lead generation form.
4. Notify your team in Slack when a critical form is submitted.

**Troubleshooting Tips**

- Ensure your Personal Access Token has the correct scopes and hasn't expired.
- Check the Typeform API documentation for any rate limits or usage guidelines.

**FAQs**

Q: How quickly will my workflow run after a form is submitted?
A: Workflows triggered by form submissions typically run within a few seconds of submission.

Q: Can I filter form responses before triggering a workflow?
A: Yes, you can set up conditions based on specific answers to determine whether a workflow should be triggered.

## Categories

- forms
- data-collection

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name          | Description                                                                                   | Link                             |
| ------------- | --------------------------------------------------------------------------------------------- | -------------------------------- |
| Get Form      | Retrieves details about a specific Typeform form including questions, settings, and metadata. | [docs](actions/get_form.md)      |
| List Forms    | Retrieves a list of all forms in your Typeform account with basic information.                | [docs](actions/list_forms.md)    |
| Get Responses | Retrieves responses submitted to a specific form with answer details.                         | [docs](actions/get_responses.md) |
| Create Form   | Creates a new form in your Typeform account with customizable fields and settings.            | [docs](actions/create_form.md)   |

## Triggers

| Name          | Description                                                                                                          | Link                              |
| ------------- | -------------------------------------------------------------------------------------------------------------------- | --------------------------------- |
| Form Response | Triggers when a new response is submitted to a specified form, providing all answer data and respondent information. | [docs](triggers/form_response.md) |
