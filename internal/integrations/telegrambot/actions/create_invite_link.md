## Create Chat Invite Link

Creates a new invite link for a chat. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights.

### Prerequisites
- **For groups**: Bot must have "can_invite_users" administrator right
- **For channels**: Bot must be an administrator

### Features
- Create custom named invite links
- Set expiration dates
- Limit the number of users who can join
- Require admin approval for new members

### Common Use Cases
1. **Event invitations**: Create time-limited links for special events
2. **Limited access**: Control group size with member limits
3. **Moderated entry**: Require approval for new members
4. **Tracking**: Named links help track invitation sources

### Important Notes
- The link will be unique and different from the chat's main invite link
- Previous links remain valid unless explicitly revoked
- member_limit and creates_join_request options are mutually exclusive