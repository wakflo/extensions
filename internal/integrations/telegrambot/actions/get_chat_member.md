## Get Chat Member

Retrieves information about a member of a chat, including their status and permissions.

### Member Statuses

The action returns one of the following statuses:

- **creator**: The user is the creator/owner of the chat
- **administrator**: The user is an administrator
- **member**: The user is a regular member
- **restricted**: The user is restricted (limited permissions)
- **left**: The user left the chat
- **kicked**: The user was removed from the chat

### Information Returned

Depending on the member's status, you'll receive:

**For Administrators:**
- Administrative permissions (can_delete_messages, can_restrict_members, etc.)
- Custom title (if set)
- Whether they can be edited by the bot

**For Restricted Members:**
- Restriction details
- Until date (if temporary)
- Restricted permissions

**For Regular Members:**
- Basic user information
- Join date (if available)

### Common Use Cases

1. **Permission Checking**: Verify if a user has specific permissions before allowing actions
2. **Member Auditing**: Check the status of specific users
3. **Moderation**: Get information about restricted or kicked users
4. **Role Verification**: Confirm admin status before processing admin commands

### Important Notes

- The bot must be a member of the chat to get member information
- In channels, the bot must be an administrator
- User IDs are numerical (e.g., 123456789), not usernames