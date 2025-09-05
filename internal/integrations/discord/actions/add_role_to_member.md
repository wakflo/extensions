# Add Role to Member

## Description

Adds a role to a member in a Discord server. This action allows you to programmatically assign roles to users within a Discord guild.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| guild-id | Guilds | Select | The ID of the Discord server (guild) where the member is located. | Yes | - |
| user-id | User ID | String | The ID of the user to add the role to. | Yes | - |
| role-id | Roles | Select | The ID of the role to assign to the member. | Yes | - |

## Returns

This action returns a success response object containing:

| Field | Type | Description |
|-------|------|-------------|
| success | Boolean | Indicates whether the role was successfully added |
| user_id | String | The ID of the user who received the role |
| guild_id | String | The ID of the Discord server where the action occurred |
| role_id | String | The ID of the role that was added |
| action | String | Action type identifier ("role_added") |
| reason | String | The reason provided for adding the role |
| message | String | Success confirmation message |

## Example Usage

You can use this action to:

1. Automatically assign roles based on user behavior or achievements
2. Grant access permissions to specific channels or features
3. Implement role-based moderation workflows
4. Onboard new members with default roles
5. Reward users with special roles for contributions

## Notes

- The bot must have the "Manage Roles" permission in the Discord server
- The bot's role must be higher in the hierarchy than the role being assigned
- User ID and Role ID must be valid Discord snowflake IDs
- The user must be a member of the specified guild
- Discord API rate limits apply to this action
- An audit log entry will be created with the optional reason if provided

## Sample Response

```json
{
  "success": true,
  "user_id": "123456789012345678",
  "guild_id": "857347647235678912",
  "role_id": "987654321098765432",
  "action": "role_added",
  "reason": "Promoted to moderator",
  "message": "Role successfully added to member"
}
```