# WhatsApp Business Integration

## Description

Integrate WhatsApp Business API with your workflows to send and receive messages, templates, and manage customer interactions. Connect with your customers on their preferred channel through automated messaging and personalized communication.

* Send text messages, media files, and template messages to customers
* Receive and respond to incoming messages automatically
* Create personalized customer experiences with rich media content
* Track message delivery and read receipts
* Automate customer service workflows and support processes
* Scale your business communications with WhatsApp's reach

**WhatsApp Business Integration Documentation**

**Overview**
The WhatsApp Business integration allows you to seamlessly connect your WhatsApp Business account with our workflow automation software, enabling you to automate messaging tasks and streamline your customer communication workflows.

**Prerequisites**

* A WhatsApp Business API account
* Our workflow automation software account
* WhatsApp Business API access token (found in your WhatsApp Business Platform dashboard)

**Setup Instructions**

1. Log in to your [Meta for Developers](https://developers.facebook.com/) account
2. Navigate to the WhatsApp Business API section and access your app
3. Generate an API token with the required permissions
4. In our workflow automation software, go to the "Integrations" section and click on "WhatsApp Business"
5. Enter the API token generated in step 3 and click "Connect"

**Available Actions**

* **Send Message**: Send a text message to a WhatsApp number
* **Send Template**: Send a pre-approved template message with dynamic parameters
* **Get Business Profile**: Retrieve your business profile information from WhatsApp

**Available Triggers**

* **Message Received**: Triggered when a message is received from a customer

**Example Use Cases**

1. **Automated Customer Support**: When a customer sends a message, automatically respond with common FAQs or route to the appropriate team
2. **Order Notifications**: Send automated order confirmations, shipping updates, and delivery notifications via WhatsApp
3. **Appointment Reminders**: Send automated reminders for upcoming appointments with options to confirm, reschedule, or cancel
4. **Lead Nurturing**: Send personalized follow-up messages to potential customers based on their interactions with your business

**Troubleshooting Tips**

* Ensure that your WhatsApp Business phone number is verified and active
* Check that your API token has the required permissions
* Verify that your message templates are pre-approved by WhatsApp before sending
* Monitor your message quality rating to ensure continued access to the WhatsApp Business API

**FAQs**

Q: Can I send messages to any WhatsApp user?
A: You can only send messages to users who have opted in to receive messages from your business or who have messaged you first within the last 24 hours.

Q: How do I create message templates?
A: Templates must be created and approved through the WhatsApp Business Platform or Meta for Business dashboard before they can be used in the integration.

Q: What media types can I send through WhatsApp?
A: You can send text, images, documents, audio files, video files, and location information.

## Categories

- messaging
- communication

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name | Description | Link |
|------|-------------|------|
| Send Message | Send a text message to a WhatsApp number | [docs](actions/send_message.md) |
| Send Template | Send a pre-approved template message with dynamic parameters | [docs](actions/send_template.md) |
| Get Business Profile | Retrieve your business profile information from WhatsApp | [docs](actions/get_business_profile.md) |

## Triggers

| Name | Description | Link |
|------|-------------|------|
| Message Received | Triggered when a message is received from a customer | [docs](triggers/message_received.md) |