# Telegram Integration

## Description

Integrate your Telegram bot with our workflow automation software to streamline communication and notifications. Connect your Telegram bot to your workflows, enabling automated message sending and responding based on events. Use this integration to:

* Send automatic notifications to channels, groups, or individual users
* Respond to user messages with customized templates
* Share files, images, and documents through your bot
* Create interactive polls and questionnaires
* Monitor channel activities and trigger workflows based on specific messages
* Build conversational bots with advanced logic

**Telegram Integration Documentation**

**Overview**
The Telegram integration allows you to seamlessly connect your Telegram bot with our workflow automation software, enabling you to automate messaging and notification tasks.

**Prerequisites**

* A Telegram bot (created through BotFather)
* The bot's API token
* Our workflow automation software account

**Setup Instructions**

1. In Telegram, talk to [@BotFather](https://t.me/botfather) to create a new bot or use an existing one.
2. Get the bot's API token from BotFather.
3. In our workflow automation software, go to the "Integrations" section and click on "Telegram".
4. Enter the API token and click "Connect".

**Available Actions**

* **Send Message**: Send a text message to a specific chat (user, group, or channel).
* **Send Photo**: Send an image to a specific chat with optional caption.
* **Get Updates**: Retrieve recent updates (messages, etc.) from your bot.

**Available Triggers**

* **Message Received**: Triggers when your bot receives a new message.

**Example Use Cases**

1. **Notification System**: Send automated notifications to a Telegram channel when specific events occur in your other integrated services.
2. **Customer Support**: Set up a Telegram bot to collect support requests and automatically create tickets in your support system.
3. **Monitoring Alerts**: Send real-time alerts to your team's Telegram group when monitoring systems detect issues.

**Troubleshooting Tips**

* Ensure that your bot has the necessary permissions in groups or channels.
* For the Message Received trigger to work, your bot must have webhook mode enabled.
* If messages aren't being delivered, check if the user has started a conversation with your bot.

**FAQs**

Q: Can I use this integration to manage multiple bots?
A: Yes, you can set up multiple instances of the Telegram integration, each connected to a different bot.

Q: Does the bot need to be an admin in a group to send messages?
A: No, but the bot must be a member of the group. For some features like pinning messages, admin rights are required.

## Categories

- messaging
- communication

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name         | Description                                                               | Link                          |
|--------------|---------------------------------------------------------------------------|-------------------------------|
| Send Message | Send a text message to a specific Telegram chat, user, group, or channel. | [docs](actions/send_message.md) |
| Send Photo   | Send a photo with optional caption to a Telegram chat.                    | [docs](actions/send_photo.md)   |
| Get Updates  | Retrieve recent updates (messages, etc.) from your Telegram bot.          | [docs](actions/get_updates.md)  |

## Triggers

| Name             | Description                                                         | Link                                  |
|------------------|---------------------------------------------------------------------|---------------------------------------|
| Message Received | Triggered when your Telegram bot receives a new message from a user | [docs](triggers/message_received.md) |