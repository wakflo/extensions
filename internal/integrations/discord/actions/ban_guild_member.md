# Ban Guild Member

## Description

Bans a member from a Discord server and optionally deletes their recent messages. This action permanently removes the user from the guild and prevents them from rejoining unless the ban is later removed.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| guild-id | Guilds | Select | The ID of the Discord server (guild) to ban the member from. | Yes | - |
| user-id | User ID | String | The ID of the user to ban from the server. | Yes | - |
| reason | Reason | Text | Optional reason for banning the member (will appear in audit log). | No | - |


## Returns

This action returns a success response object containing:

| Field | Type | Description |
|-------|------|-------------|
| success | Boolean | Indicates whether the member was successfully banned |
| user_id | String | The ID of the user who was banned |
| guild_id | String | The ID of the Discord server where the ban occurred |
| action | String | Action type identifier ("banned") |
| reason | String | The reason provided for the ban |
| message | String | Success confirmation message |

## Example Usage

You can use this action to:

1. Implement automated moderation for rule violations
2. Ban users who engage in spam or harassment
3. Remove malicious users from your Discord community
4. Enforce server rules through programmatic moderation
5. Create automated anti-raid protection systems

## Notes

- The bot must have the "Ban Members" permission in the Discord server
- The bot's role must be higher in the hierarchy than the target user's highest role
- You cannot ban the server owner
- The user does not need to be currently in the server to be banned (preventive banning)
- `delete_message_seconds` can range from 0 to 604800 seconds (7 days)
- An audit log entry will be created with the provided reason
- Discord API rate limits apply to this action
- This is a destructive action and should be used carefully

## Sample Response

```json
{
  "success": true,
  "user_id": "123456789012345678",
  "guild_id": "857347647235678912",
  "action": "banned",
  "reason": "Spam and harassment",
  "message": "Member successfully banned from guild"
}
```

## Error Handling

Common errors that may occur:
- **403 Forbidden**: Bot lacks "Ban Members" permission or user has higher role
- **404 Not Found**: Guild or user not found
- **400 Bad Request**: Invalid parameters or attempting to ban the server owner