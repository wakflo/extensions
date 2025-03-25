# Campaign Monitor Integration

## Description

Integrate Campaign Monitor with your workflow automation software to seamlessly manage email marketing campaigns. Connect your Campaign Monitor account to automate subscriber management, campaign creation, and reporting processes. This integration allows you to:

* Automatically add or update subscribers based on form submissions or user actions
* Create and send email campaigns when triggered by specific events
* Sync campaign performance data with your CRM or analytics tools
* Create targeted segments based on user behavior in your application
* Get real-time notifications about campaign activities and subscriber changes

**Campaign Monitor Integration Documentation**

**Overview**
The Campaign Monitor integration allows you to seamlessly connect your Campaign Monitor account with our workflow automation software, enabling you to automate email marketing tasks and streamline your workflow.

**Prerequisites**

* A Campaign Monitor account
* Our workflow automation software account
* Campaign Monitor API key and Client ID (found in your Campaign Monitor account settings)

**Setup Instructions**

1. Log in to your Campaign Monitor account and navigate to the "Account Settings".
2. Click on "API Keys" and generate a new API key or select an existing one.
3. Note your Client ID from the account dashboard.
4. In our workflow automation software, go to the "Integrations" section and click on "Campaign Monitor".
5. Enter the API key and Client ID generated in steps 2 and 3, then click "Connect".

**Available Actions**

* **List Campaigns**: Retrieve a list of campaigns from your Campaign Monitor account.
* **Get Campaign**: Retrieve details of a specific campaign.
* **Create Campaign**: Create a new email campaign.
* **List Subscribers**: Retrieve a list of subscribers from a specific list.
* **Add Subscriber**: Add a new subscriber to a specific list.

**Available Triggers**

* **Subscriber Added**: Trigger a workflow when a new subscriber is added to a list.
* **Campaign Sent**: Trigger a workflow when a campaign is sent.

**Example Use Cases**

1. Automate welcome emails: When a new user signs up on your website, automatically add them to a Campaign Monitor list and trigger a welcome email sequence.
2. Sync with CRM: When a campaign is sent, update your CRM with the campaign information and statistics.
3. Segment customers: When a purchase is made, add the customer to a specific Campaign Monitor list based on what they purchased.

**Troubleshooting Tips**

* Ensure that you have entered the correct API key and Client ID.
* Check that your Campaign Monitor account has the necessary permissions to perform the requested actions.
* Review Campaign Monitor's API documentation for any rate limits or usage guidelines.

**FAQs**

Q: Can I use this integration to send personalized emails based on user behavior?
A: Yes, you can use the Add Subscriber action to add users to specific segments based on their behavior, and then use those segments for targeted campaigns.

Q: How real-time are the triggers?
A: The triggers are based on polling, so there may be a slight delay (usually a few minutes) between the event occurring in Campaign Monitor and the trigger firing in your workflow.

## Categories

- marketing
- email

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name             | Description                                                                                                            | Link                                |
|------------------|------------------------------------------------------------------------------------------------------------------------|-------------------------------------|
| List Campaigns   | Retrieve a list of campaigns from your Campaign Monitor account.                                                       | [docs](actions/list_campaigns.md)   |
| Get Campaign     | Retrieve details of a specific campaign.                                                                               | [docs](actions/get_campaign.md)     |
| Create Campaign  | Create a new email campaign.                                                                                           | [docs](actions/create_campaign.md)  |
| List Subscribers | Retrieve a list of subscribers from a specific list.                                                                   | [docs](actions/list_subscribers.md) |
| Add Subscriber   | Add a new subscriber to a specific list.                                                                               | [docs](actions/add_subscriber.md)   |

## Triggers

| Name             | Description                                                                                                            | Link                                     |
|------------------|------------------------------------------------------------------------------------------------------------------------|------------------------------------------|
| Subscriber Added | Trigger a workflow when a new subscriber is added to a list.                                                           | [docs](triggers/subscriber_added.md)     |
| Campaign Sent    | Trigger a workflow when a campaign is sent.                                                                            | [docs](triggers/campaign_sent.md)        |