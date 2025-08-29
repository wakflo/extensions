# Pinterest Integration

## Description

Integrate your Pinterest marketing and content creation with our workflow automation software. Connect your Pinterest account to automatically create pins, manage boards, and track engagement metrics. Use this integration to:

- Create and schedule pins automatically from your content workflow
- Organize boards based on triggers from other integrated systems
- Track engagement metrics and sync with your analytics tools
- Automate marketing campaigns between multiple platforms
- Keep your social media strategy consistent with automated pin creation

**Pinterest Integration Documentation**

**Overview**
The Pinterest integration allows you to connect your Pinterest account with our workflow automation software, enabling you to automate tasks and streamline your Pinterest marketing workflow.

**Prerequisites**

- A Pinterest business account
- Our workflow automation software account
- Pinterest API credentials (found in your Pinterest Developer Portal)

**Setup Instructions**

1. Log in to your Pinterest Developer Portal and create a new app
2. Configure the OAuth settings and permissions for your app
3. Obtain the Client ID and Client Secret
4. In our workflow automation software, go to the "Integrations" section and click on "Pinterest"
5. Enter the Client ID and Client Secret generated in step 3 and authorize the connection

**Available Actions**

- **Get Board**: Retrieve details of a specific Pinterest board by ID
- **Get Pin**: Retrieve details of a specific pin by ID
- **Create Pin**: Create a new pin on a specified board with customizable content
- **List Boards**: Retrieve a list of boards from a Pinterest account
- **Search Pins**: Search for pins based on keywords and filters

**Available Triggers**

- **Pin Created**: Trigger a workflow when a new pin is created on a specified board

**Example Use Cases**

1. Content Automation: When a new blog post is published, automatically create a Pinterest pin with the featured image and link
2. Cross-Platform Promotion: When content is shared on Instagram, automatically create matching pins on relevant Pinterest boards
3. Marketing Campaigns: Schedule and post a series of pins based on your content calendar triggers

**Troubleshooting Tips**

- Ensure you have the correct permission scopes enabled for your Pinterest app
- Check that your OAuth credentials are entered correctly
- Verify that your Pinterest business account is properly connected to the app

**FAQs**

Q: Can I schedule pins for future publication using this integration?
A: Yes, you can create pins that will be published at a scheduled future date.

Q: How many pins can I create at once?
A: Pinterest has rate limits that allow creation of approximately 10 pins per minute and 200 pins per day.

## Categories

- marketing
- social-media

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name        | Description                                                                                  | Link                           |
| ----------- | -------------------------------------------------------------------------------------------- | ------------------------------ |
| Get Board   | Retrieves detailed information about a specific Pinterest board using its unique identifier. | [docs](actions/get_board.md)   |
| Get Pin     | Retrieves detailed information about a specific Pinterest pin using its unique identifier.   | [docs](actions/get_pin.md)     |
| Create Pin  | Creates a new pin on Pinterest with specified image, title, description, and board.          | [docs](actions/create_pin.md)  |
| List Boards | Retrieves all boards associated with the authenticated Pinterest account.                    | [docs](actions/list_boards.md) |
| Search Pins | Searches for pins on Pinterest matching specific criteria or keywords.                       | [docs](actions/search_pins.md) |

## Triggers

| Name        | Description                                                                   | Link                            |
| ----------- | ----------------------------------------------------------------------------- | ------------------------------- |
| Pin Created | Triggers a workflow when a new pin is created on a specified Pinterest board. | [docs](triggers/pin_created.md) |
