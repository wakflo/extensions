# Mailchimp Integration

## Description

Automate your email marketing workflows with seamless integration between our workflow automation software and Mailchimp. Easily sync subscriber data, trigger campaigns, and automate email sends based on custom conditions. Say goodbye to manual data imports and hello to streamlined email marketing operations.

**Mailchimp Integration Documentation**

**Overview**
The Mailchimp integration allows you to automate workflows with your email marketing campaigns in Mailchimp. With this integration, you can trigger actions in Mailchimp based on specific events or conditions in your workflow automation software.

**Prerequisites**

* A Mailchimp account
* A workflow automation software account
* The Mailchimp API key (available in the Mailchimp dashboard)

**Setup**

1. In your workflow automation software, navigate to the "Integrations" or "Connections" section.
2. Search for and select the "Mailchimp" integration.
3. Click "Connect" to initiate the setup process.
4. Enter your Mailchimp API key and click "Save".
5. Authorize the integration by clicking the "Authorize" button.

**Triggers**

The following triggers are available in the Mailchimp integration:

* **New Subscriber**: Triggered when a new subscriber is added to a list in Mailchimp.
* **Subscriber Update**: Triggered when an existing subscriber's information is updated in Mailchimp.
* **Campaign Sent**: Triggered when a campaign is sent from Mailchimp.

**Actions**

The following actions are available in the Mailchimp integration:

* **Add Subscriber**: Adds a new subscriber to a specified list in Mailchimp.
* **Update Subscriber**: Updates an existing subscriber's information in Mailchimp.
* **Remove Subscriber**: Removes a subscriber from a specified list in Mailchimp.
* **Send Campaign**: Sends a campaign from Mailchimp.

**Examples**

1. Trigger: New Subscriber
Action: Add Subscriber to List A

When a new subscriber is added to your Mailchimp list, the workflow automation software will automatically add them to List A in Mailchimp.

2. Trigger: Campaign Sent
Action: Send Follow-up Campaign

When a campaign is sent from Mailchimp, the workflow automation software will automatically send a follow-up campaign to subscribers who did not open or click on the original campaign.

**Troubleshooting**

* If you encounter issues with the integration, check that your API key is correct and that you have authorized the integration.
* If you are experiencing errors or delays in triggering actions, check the Mailchimp API documentation for any rate limits or usage guidelines.

**FAQs**

* Q: Can I integrate multiple Mailchimp accounts?
A: Yes, you can integrate multiple Mailchimp accounts by setting up separate connections in your workflow automation software.
* Q: Can I use this integration to automate workflows with other Mailchimp features (e.g. automations, segments)?
A: No, this integration is specifically designed for automating workflows with email marketing campaigns in Mailchimp.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>


## Actions


| Name                       | Description                                                                                                                                                                                                                                                                                                                  | Link                                            |
|----------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------|
| Add Member To List         | Adds a new member to an existing list in your workflow, allowing you to easily manage and track team members, stakeholders, or other relevant parties involved in the process. This integration action enables seamless addition of new members to lists, streamlining collaboration and communication within your workflow. | [docs](actions/add_member_to_list.md)           |## Actions
| Add Note To Subscribe      | Adds a note to a subscription, allowing you to record important information or comments related to the subscription. This integration action is useful for tracking updates, issues, or other relevant details about a specific subscription.                                                                                | [docs](actions/add_note_to_subscribe.md)        |
| Add Note To Subscriber     | Adds a note to a subscriber's record, allowing you to store additional information and context about the subscriber. This integration action is useful for tracking important details or updates related to a specific subscriber, making it easier to manage and analyze their activity over time.                          | [docs](actions/add_note_to_subscriber.md)       |
| Add Subscriber To Tag      | Add Subscriber to Tag: Automatically adds one or more subscribers to a specific tag in your email marketing platform, allowing you to easily manage and segment your audience.                                                                                                                                               | [docs](actions/add_subscriber_to_list.md)       |
| Get All List               | Retrieves all lists associated with a specific entity or resource, allowing you to access and utilize tag metadata in your workflow automation processes.                                                                                                                                                                    | [docs](actions/get_all_list)                    |
| Remove Subscriber From Tag | Remove Subscriber From Tag: This integration action removes a subscriber from a specific tag in your email marketing platform, ensuring that the individual is no longer associated with the designated group.                                                                                                               | [docs](actions/remove_subscriber_from_tag.md)   |
| Update Subscriber Status   | Updates the status of a subscriber in your application or database, allowing you to reflect changes in their subscription level, account information, or other relevant details.                                                                                                                                             | [docs](actions/update_subscriber_status.md)     |


## Triggers

| Name             | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 | Link                                 |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------|
| New Subscriber   | Triggers when a new subscriber is added to your application or service, allowing you to automate tasks and workflows immediately after subscription.                                                                                                                                                                                                                                                                                                                                                                                                                        | [docs](triggers/new_subscriber.md)   |## Triggers
| Unsubscriber     | The Unsubscriber integration trigger is designed to automatically remove subscribers from your workflow when they unsubscribe from a specific email list or service. This trigger can be used in conjunction with other automation workflows to ensure that unsubscribed contacts are no longer targeted for marketing campaigns, surveys, or other automated tasks. By integrating the Unsubscriber trigger with your existing workflows, you can maintain data accuracy and compliance with anti-spam laws by promptly removing unsubscribed contacts from your workflow. | [docs](triggers/unsubscriber.md)     |