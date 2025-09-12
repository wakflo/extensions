# Discord Integration

## Description

Integrate Discord with your workflow automation software to streamline communication and moderation. This integration allows you to automate Discord server management, notifications, and event handling. Use this integration to:

* Send automated messages to channels or users
* Create, manage, and organize channels automatically
* Monitor and respond to server events like new members or messages
* Build interactive workflows triggered by Discord activities
* Update role assignments and permissions based on external events
* Deliver real-time notifications from your applications to Discord

**Discord Integration Documentation**

**Overview**
The Discord integration connects your Discord server with our workflow automation software, enabling you to automate tasks and build responsive workflows based on Discord activities.

**Prerequisites**

* A Discord account with administrative permissions on a server
* Our workflow automation software account
* A Discord Bot token (created through the Discord Developer Portal)

**Setup Instructions**

1. Visit the [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Navigate to the "Bot" section and create a new bot
4. Copy the bot token
5. In our workflow automation software, go to the "Integrations" section and click on "Discord"
6. Enter the Bot token and click "Connect"
7. Use the OAuth2 URL provided to invite the bot to your server

**Available Actions**

* **Send Message**: Send a message to a Discord channel or user
* **Create Channel**: Create a new text or voice channel in a Discord server
* **Get Guild**: Retrieve information about a Discord server (guild)
* **List Channels**: Get a list of channels in a Discord server

**Available Triggers**

* **Message Received**: Triggered when a new message is received in a monitored channel
* **Member Joined**: Triggered when a new member joins the Discord server

**Troubleshooting Tips**

* Ensure that your bot has the necessary permissions for the actions you are trying to perform
* Check that the bot is properly invited to your server
* Verify that the bot token is correctly entered in the integration settings

**FAQs**

Q: What permissions does the bot need?
A: The required permissions depend on the actions you want to perform. For sending messages, the bot needs "Send Messages" permission. For channel management, it needs "Manage Channels" permission.

Q: Can I use this integration to automate moderation?
A: Yes, you can set up workflows that respond to specific message patterns or user actions, automating some moderation tasks.

## Categories

- communication

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name | Description | Link |
|------|-------------|------|
| Send Message | Send a message to a specified Discord channel or user | [docs](actions/send_message.md) |
| Create Channel | Create a new text or voice channel in a Discord server | [docs](actions/create_channel.md) |
| Get Guild | Retrieve information about a Discord server (guild) | [docs](actions/get_guild.md) |
| List Channels | Get a list of channels in a Discord server | [docs](actions/list_channels.md) |

## Triggers

| Name | Description | Link |
|------|-------------|------|
| Message Received | Triggered when a new message is received in a monitored channel | [docs](triggers/message_received.md) |
| Member Joined | Triggered when a new member joins the Discord server | [docs](triggers/member_joined.md) |