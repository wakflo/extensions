# Ghost Integration

## Description

Ghost is a powerful open-source publishing platform that enables creators to build and manage modern online publications. This integration allows you to automate your Ghost CMS workflows, including:

* Create, update, and manage posts and pages programmatically
* Manage members and subscriptions for your publication
* Organize content with tags and authors
* Monitor publication activity with real-time triggers
* Upload and manage media assets
* Synchronize content across multiple platforms
* Automate newsletter distribution and member management

## Prerequisites

* A Ghost publication (self-hosted or Ghost(Pro))
* Admin API integration with Admin API Key
* The latest version of Ghost (v5.0 or higher recommended)
* Your Ghost site URL

## Setup Instructions

1. **Access Ghost Admin Panel**
   - Log in to your Ghost Admin panel at `https://yourblog.ghost.io/ghost`
   - Navigate to Settings → Integrations

2. **Create Custom Integration**
   - Click "Add custom integration"
   - Give your integration a name (e.g., "Wakflo Integration")
   - Ghost will generate an Admin API Key

3. **Copy Integration Details**
   - Copy the Admin API Key (format: `site-name:key`)
   - Note your Ghost site URL (e.g., `https://yourblog.ghost.io`)

4. **Configure in Wakflo**
   - Go to the Integrations section in Wakflo
   - Select "Ghost" from the available integrations
   - Enter your Ghost site URL
   - Enter the Admin API Key from step 3
   - Click "Connect" to establish the connection

5. **Test the Connection**
   - Try creating a test post or listing existing posts
   - Verify that the integration is working correctly

## Authentication

This integration uses Ghost's Admin API authentication with JWT tokens. The Admin API Key provides full access to create, edit, and delete content in your Ghost publication.

### Authentication Process
1. The integration splits the Admin API key into ID and secret
2. Creates a JWT token signed with the secret
3. Includes the token in the Authorization header for all requests
4. Tokens expire after 5 minutes and are regenerated automatically

## API Version

This integration uses Ghost Admin API v5.0 for maximum compatibility and feature support.

## Rate Limits

Ghost API has the following rate limits:
* **Admin API**: 120 requests per minute
* **Content API**: 300 requests per minute (not used in this integration)

The integration handles rate limiting automatically and will retry failed requests when appropriate.

## Content Formats

Ghost supports multiple content formats:

### HTML
Standard HTML content for posts and pages.

### Markdown
Write content in Markdown format, automatically converted to HTML.

### Lexical
Ghost's new rich text editor format (recommended for Ghost 5.0+).

## Error Handling

The integration provides detailed error messages for common issues:
* Invalid API key format
* Unauthorized access (expired or invalid key)
* Post/Page not found
* Validation errors (missing required fields)
* Rate limit exceeded

## Webhooks Support

While this integration primarily uses polling for triggers, Ghost also supports webhooks for real-time events. Contact support if you need webhook-based triggers.

## Categories

- cms
- publishing
- content-management
- blogging
- newsletter

## Authors

- Wakflo <integrations@wakflo.com>

## Support

For issues or questions about this integration:
- Email: integrations@wakflo.com
- Documentation: https://docs.wakflo.com/integrations/ghost
- Ghost API Docs: https://ghost.org/docs/admin-api/

## Actions

| Name | Description | Type |
|------|-------------|------|
| **Create Post** | Creates a new post in your Ghost publication with full control over content, metadata, and publishing options | Action |
| **Update Post** | Updates an existing post including content, status, tags, and all metadata | Action |
| **Delete Post** | Permanently deletes a post from your Ghost publication | Action |
| **Get Post** | Retrieves a specific post by ID or slug with all related data | Action |
| **List Posts** | Lists all posts with advanced filtering, pagination, and sorting options | Action |
| **Create Page** | Creates a new static page in Ghost for permanent content | Action |
| **Update Page** | Updates an existing page's content and metadata | Action |
| **Delete Page** | Permanently deletes a page from your publication | Action |
| **Get Page** | Retrieves a specific page by ID or slug | Action |
| **List Pages** | Lists all pages with filtering and pagination options | Action |
| **Create Member** | Creates a new member in your Ghost publication with subscription options | Action |
| **Update Member** | Updates an existing member's information and subscription status | Action |
| **Delete Member** | Removes a member from your publication | Action |
| **Get Member** | Retrieves a specific member by ID or email | Action |
| **List Members** | Lists all members with filtering and pagination options | Action |
| **Create Tag** | Creates a new tag for organizing content | Action |
| **Update Tag** | Updates an existing tag's properties and metadata | Action |
| **Delete Tag** | Deletes a tag from your Ghost publication | Action |
| **Get Tag** | Retrieves a specific tag by ID or slug | Action |
| **List Tags** | Lists all tags in your publication | Action |
| **Upload Image** | Uploads an image to Ghost's media library | Action |
| **Get Site** | Retrieves site-wide settings and configuration information | Action |
| **List Authors** | Lists all authors/users in your publication | Action |
| **Get Author** | Retrieves a specific author by ID or slug | Action |

## Triggers

| Name | Description | Type | Polling Interval |
|------|-------------|------|------------------|
| **New Post** | Triggers when a new post is created in any status | Polling | 5 minutes |
| **Post Updated** | Triggers when an existing post is updated | Polling | 5 minutes |
| **Post Published** | Triggers when a post is published (status changes to published) | Polling | 5 minutes |
| **New Member** | Triggers when a new member signs up or is added | Polling | 5 minutes |
| **Member Updated** | Triggers when a member's information is updated | Polling | 5 minutes |
| **New Page** | Triggers when a new page is created | Polling | 5 minutes |
| **Page Updated** | Triggers when an existing page is updated | Polling | 5 minutes |
| **New Tag** | Triggers when a new tag is created | Polling | 5 minutes |

## Common Use Cases

### Content Automation
- Automatically publish blog posts from other sources
- Cross-post content to multiple platforms
- Schedule posts based on external triggers

### Member Management
- Sync members with CRM systems
- Automate welcome emails for new members
- Update member segments based on behavior

### Content Organization
- Auto-tag posts based on content analysis
- Organize content into collections
- Maintain consistent taxonomy across platforms

### Analytics and Reporting
- Export post performance data
- Track member engagement
- Generate content calendars

### Backup and Migration
- Regular content backups
- Migrate content between Ghost instances
- Archive published content

## Troubleshooting

### Common Issues

**Authentication Failed**
- Verify your Admin API key format (should be `site-name:key`)
- Ensure the integration has not been deleted in Ghost
- Check that your Ghost site URL is correct

**Post/Page Not Found**
- Verify the ID or slug is correct
- Ensure the content hasn't been deleted
- Check if you have the right permissions

**Rate Limit Exceeded**
- Wait a few minutes before retrying
- Consider implementing request batching
- Reduce polling frequency for triggers

**Invalid Content Format**
- Ensure HTML is properly escaped
- Check Markdown syntax for errors
- Verify Lexical JSON structure

## Changelog

### Version 0.0.1
- Initial release
- Support for all major Ghost Admin API endpoints
- Posts, Pages, Members, Tags management
- Image upload capability
- Polling-based triggers
- SDK v2 compatibility

## License

This integration is provided under the MIT License. See LICENSE file for details.

## Contributing

We welcome contributions! Please see our contributing guidelines for more information on how to get involved.

## Disclaimer

This integration is not officially affiliated with Ghost Foundation. Ghost is a trademark of Ghost Foundation.