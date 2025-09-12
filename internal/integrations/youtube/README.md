# YouTube Integration

## Description

Integrate YouTube with your workflow automation software to streamline video management and analysis. This integration connects your YouTube channel with our platform, allowing you to automate tasks related to video publishing, monitoring, and analysis. Use this integration to:

- Automatically upload videos based on your workflow triggers
- Monitor new video uploads and comments
- Track video performance metrics and engagement
- Update video metadata, titles, descriptions, and privacy settings
- Create automated responses to comments based on specific criteria
- Sync video data with other systems and applications

## Overview

The YouTube integration allows you to connect your YouTube channel to our workflow automation platform, enabling you to automate tasks related to video management and content creation.

## Prerequisites

- A Google account with a YouTube channel
- Our workflow automation software account
- YouTube Data API v3 access (enabled in Google Cloud Console)
- OAuth 2.0 credentials for the YouTube Data API

## Setup Instructions

1. Log in to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the YouTube Data API v3 for your project
4. Create OAuth 2.0 credentials:
   - Go to "Credentials" in the API & Services section
   - Click "Create Credentials" and select "OAuth client ID"
   - Set up the OAuth consent screen if prompted
   - Select "Web application" as the application type
   - Add authorized redirect URIs (provided by our platform)
   - Download the credentials
5. In our workflow automation platform, navigate to the Integrations section
6. Select the YouTube integration
7. Enter your OAuth credentials (Client ID and Client Secret)
8. Complete the OAuth flow to authorize access to your YouTube account

## Available Actions

- **List Videos**: Retrieve a list of videos from your YouTube channel with filtering options
- **Get Video**: Retrieve detailed information about a specific video
- **Upload Video**: Upload a new video to your YouTube channel
- **Update Video**: Update metadata for an existing video

## Available Triggers

- **Video Uploaded**: Trigger workflows when a new video is uploaded to your YouTube channel

## Example Use Cases

1. **Automated Publishing Workflow**: Automatically upload videos to YouTube when they're approved in your content management system
2. **Cross-Platform Promotion**: When a new video is uploaded to YouTube, automatically share it on other social media platforms
3. **Comment Moderation**: Monitor new comments on videos and take action based on sentiment analysis
4. **Performance Tracking**: Automatically export video performance data to analytics tools or dashboards
5. **Content Repurposing**: When a video reaches a certain view count, trigger a workflow to repurpose it for other platforms

## Troubleshooting Tips

- Ensure your OAuth credentials are correctly configured with the necessary scopes
- Check that your Google Cloud project has the YouTube Data API v3 enabled
- Verify that your redirect URIs match exactly what's configured in Google Cloud Console
- Remember that YouTube API quotas are limited - design your workflows to minimize API calls

## Categories

- social
- video

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name         | Description                                                                            | Link                            |
| ------------ | -------------------------------------------------------------------------------------- | ------------------------------- |
| List Videos  | Retrieve a list of videos from your YouTube channel with various filtering options     | [docs](actions/list_videos.md)  |
| Get Video    | Retrieve detailed information about a specific video by its ID                         | [docs](actions/get_video.md)    |
| Upload Video | Upload a new video to your YouTube channel with customizable metadata                  | [docs](actions/upload_video.md) |
| Update Video | Update metadata for an existing video, including title, description, tags, and privacy | [docs](actions/update_video.md) |

## Triggers

| Name           | Description                                                            | Link                               |
| -------------- | ---------------------------------------------------------------------- | ---------------------------------- |
| Video Uploaded | Trigger workflows when a new video is uploaded to your YouTube channel | [docs](triggers/video_uploaded.md) |
