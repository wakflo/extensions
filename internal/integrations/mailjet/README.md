# MailJet Integration

## Description

Integrate MailJet's email marketing and transactional email services with your workflow automation. Send personalized emails, manage contacts, track deliveries, and automate communications based on triggers or customer actions. This integration enables you to:

* Send transactional emails directly from your workflows
* Create and update contact information in your MailJet lists
* Track email deliveries, opens, and clicks to automate follow-up actions
* Trigger workflows based on email events (sent, opened, clicked)
* Automatically manage contact subscriptions and preferences
* Access contact data to personalize other workflow steps

**MailJet Integration Documentation**

**Overview**
The MailJet integration allows you to seamlessly connect your MailJet account with our workflow automation software, enabling you to automate email tasks and streamline your workflow.

**Prerequisites**

* A MailJet account
* Our workflow automation software account
* Your MailJet API Key and Secret Key (found in your MailJet API Keys section under Account Settings)

**Setup Instructions**

1. Log in to your MailJet account and navigate to the "Account Settings" section.
2. Click on "REST API" and locate your API Key and Secret Key.
3. In our workflow automation software, go to the "Integrations" section and click on "MailJet".
4. Enter the API Key and Secret Key generated in step 2 and click "Connect".

**Available Actions**

* **Send Email**: Send transactional or marketing emails to your contacts.
* **Get Contact**: Retrieve detailed information about a specific contact.
* **Create Contact**: Add a new contact to your MailJet contact list.
* **List Contacts**: Retrieve a list of contacts from your MailJet account.

**Available Triggers**

* **Email Sent**: Triggered when an email is successfully sent through MailJet.
* **Contact Updated**: Triggered when a contact's information is updated in MailJet.

**Example Use Cases**

1. **Automated Welcome Email**: When a new user signs up on your website, automatically add them to your MailJet contacts and send a personalized welcome email.
2. **Re-engagement Campaign**: When a user performs a specific action in your application, send a targeted email to encourage further engagement.
3. **Customer Support Workflow**: When a support ticket is created, send an automated confirmation email to the customer.

**Troubleshooting Tips**

* Ensure that you have entered the correct API Key and Secret Key.
* Check MailJet's API documentation for any rate limits or usage guidelines.
* Verify that the email addresses you're sending to are valid and not on your blocklist.

## Categories

- email
- marketing

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name            | Description                                                                                     | Link                             |
|-----------------|-------------------------------------------------------------------------------------------------|----------------------------------|
| Send Email      | Send transactional or marketing emails to your contacts with personalized content.             | [docs](actions/send_email.md)    |
| Get Contact     | Retrieve detailed information about a specific contact in your MailJet account.                | [docs](actions/get_contact.md)   |
| Create Contact  | Add a new contact to your MailJet contact list with custom properties.                         | [docs](actions/create_contact.md)|
| List Contacts   | Retrieve a filtered list of contacts from your MailJet account based on specified criteria.    | [docs](actions/list_contacts.md) |

## Triggers

| Name            | Description                                                                                      | Link                                |
|-----------------|--------------------------------------------------------------------------------------------------|-------------------------------------|
| Email Sent      | Triggers a workflow when an email is successfully sent through MailJet.                          | [docs](triggers/email_sent.md)      |
| Contact Updated | Triggers a workflow when a contact's information is updated in your MailJet account.             | [docs](triggers/contact_updated.md) |