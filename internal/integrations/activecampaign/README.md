# ActiveCampaign Integration

## Description

Integrate your ActiveCampaign marketing automation data with our workflow automation software to streamline customer engagement and marketing processes. Automatically sync ActiveCampaign contacts, campaigns, and automations with your workflows, eliminating manual data entry and improving efficiency. Use this integration to:

* Automatically create or update contacts in ActiveCampaign when new leads are generated
* Trigger workflows based on campaign interactions or automation completions
* Segment contacts based on workflow data and update lists accordingly
* Send personalized emails based on user behavior or data changes
* Track campaign performance and engagement metrics in your dashboards

**ActiveCampaign Integration Documentation**

**Overview**
The ActiveCampaign integration allows you to seamlessly connect your ActiveCampaign account with our workflow automation software, enabling you to automate marketing tasks and streamline your customer engagement processes.

**Prerequisites**

* An ActiveCampaign account
* Our workflow automation software account
* Your ActiveCampaign API URL and API Key (found in your ActiveCampaign account settings)

**Setup Instructions**

1. Log in to your ActiveCampaign account and navigate to "Settings" > "Developer".
2. Copy your API URL and API Key.
3. In our workflow automation software, go to the "Integrations" section and click on "ActiveCampaign".
4. Enter the API URL and API Key, then click "Connect".

**Available Actions**

* **List Contacts**: Retrieve a list of contacts from your ActiveCampaign account with filtering options.
* **Get Contact**: Retrieve a specific contact by ID from your ActiveCampaign account.
* **Create Contact**: Create a new contact in your ActiveCampaign account with customizable fields.
* **Update Contact**: Update an existing contact in your ActiveCampaign account.

**Available Triggers**

* **Contact Updated**: Automatically trigger workflows when a contact is updated in ActiveCampaign.

**Example Use Cases**

1. Lead Nurturing: When a new lead is added to your CRM, create a contact in ActiveCampaign and add them to a specific automation sequence.
2. Follow-up Automation: When a contact completes a specific action on your website, update their custom fields in ActiveCampaign to trigger personalized follow-up emails.
3. List Management: Automatically add contacts to specific ActiveCampaign lists based on their behavior or information from other integrated systems.

**Troubleshooting Tips**

* Ensure that you have entered the correct API URL and API Key.
* Check that your ActiveCampaign account has sufficient permissions for the actions you're trying to perform.
* Review the ActiveCampaign API documentation for rate limits or usage guidelines.

**FAQs**

Q: How often are contacts synced from ActiveCampaign?
A: The Contact Updated trigger works on a polling basis, checking for updates every 5 minutes.

Q: Can I create custom fields through this integration?
A: Currently, this integration supports adding values to existing custom fields, but not creating new custom fields.

## Categories

- marketing
- crm

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name | Description | Link |
|------|-------------|------|
| List Contacts | Retrieve a list of contacts from your ActiveCampaign account with filtering options. | [docs](actions/list_contacts.md) |
| Get Contact | Retrieve a specific contact by ID from your ActiveCampaign account. | [docs](actions/get_contact.md) |
| Create Contact | Create a new contact in your ActiveCampaign account with customizable fields. | [docs](actions/create_contact.md) |
| Update Contact | Update an existing contact in your ActiveCampaign account. | [docs](actions/update_contact.md) |

## Triggers

| Name | Description | Link |
|------|-------------|------|
| Contact Updated | Automatically trigger workflows when a contact is updated in ActiveCampaign. | [docs](triggers/contact_updated.md) |