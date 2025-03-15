# ConvertKit Integration

## Description

Integrate your ConvertKit email marketing platform with our workflow automation software to streamline your email marketing campaigns and subscriber management. Use this integration to:

* Automatically add new subscribers to your ConvertKit account when they fill out forms in other systems
* Create and apply tags to organize your subscribers
* Trigger workflows when new subscribers are added
* Automate email sequences based on user actions in other platforms
* Sync subscriber data between ConvertKit and your CRM, e-commerce, or other business tools

**ConvertKit Integration Documentation**

**Overview**
The ConvertKit integration allows you to seamlessly connect your ConvertKit account with our workflow automation software, enabling you to automate email marketing tasks and streamline your workflow.

**Prerequisites**

* A ConvertKit account
* Our workflow automation software account
* ConvertKit API key (found in your ConvertKit account settings)

**Setup Instructions**

1. Log in to your ConvertKit account and navigate to the "Account Settings" page.
2. Click on "API" and copy your API key.
3. In our workflow automation software, go to the "Integrations" section and click on "ConvertKit".
4. Enter your API key and click "Connect".

**Available Actions**

* **List Subscribers**: Retrieve a list of subscribers from your ConvertKit account.
* **Get Subscriber**: Retrieve details for a specific subscriber.
* **Create Subscriber**: Add a new subscriber to your ConvertKit account.
* **Create Tag**: Create a new tag in your ConvertKit account.
* **Add Tag to Subscriber**: Apply a tag to a specific subscriber.

**Available Triggers**

* **Subscriber Created**: Trigger a workflow when a new subscriber is added to your ConvertKit account.
* **Tag Added**: Trigger a workflow when a tag is added to a subscriber.

**Example Use Cases**

1. When a new customer makes a purchase in your e-commerce platform, automatically add them as a subscriber in ConvertKit and apply a "customer" tag.
2. When a subscriber is added to ConvertKit, create a corresponding contact in your CRM system.
3. When a lead submits a form on your website, add them to ConvertKit and trigger a welcome email sequence.

**Troubleshooting Tips**

* Ensure that your API key is correct and has not been reset.
* Check that your ConvertKit account has the necessary permissions for the actions you're trying to perform.
* Verify that the fields you're mapping match the expected format in ConvertKit.

**FAQs**

Q: How quickly will new subscribers in ConvertKit trigger my workflows?
A: The integration polls for new subscribers every 5 minutes, so workflows will be triggered within that timeframe.

Q: Can I use custom fields from ConvertKit in my workflows?
A: Yes, all custom fields associated with subscribers are available for use in your workflows.

## Categories

- marketing
- email

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name                  | Description                                                                                                    | Link                                     |
|-----------------------|----------------------------------------------------------------------------------------------------------------|------------------------------------------|
| List Subscribers      | Retrieve a list of subscribers from your ConvertKit account with their details and tags.                        | [docs](actions/list_subscribers.md)      |
| Get Subscriber        | Retrieve detailed information about a specific subscriber using their ID or email address.                      | [docs](actions/get_subscriber.md)        |
| Create Subscriber     | Add a new subscriber to your ConvertKit account with email, first name, and optional custom fields.             | [docs](actions/create_subscriber.md)     |
| Create Tag            | Create a new tag in your ConvertKit account to help organize and segment your subscribers.                      | [docs](actions/create_tag.md)            |
| Add Tag to Subscriber | Apply an existing tag to a specific subscriber to enhance segmentation and targeting capabilities.              | [docs](actions/add_tag_to_subscriber.md) |

## Triggers

| Name               | Description                                                                                                                                                                       | Link                                |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| Subscriber Created | Triggers a workflow when a new subscriber is added to your ConvertKit account, allowing you to automate follow-up actions or sync the data with other systems.                     | [docs](triggers/subscriber_created.md) |
| Tag Added          | Triggers a workflow when a tag is added to a subscriber in your ConvertKit account, enabling you to create targeted automations based on subscriber segmentation and behavior.     | [docs](triggers/tag_added.md)          |