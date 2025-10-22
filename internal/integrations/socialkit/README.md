# SocialKit Integration

## Description

Integrate SocialKit's social media video analysis API with your workflow automation. Extract video summaries, transcripts, engagement metrics, and custom insights from YouTube, TikTok, and Instagram videos. Use this integration to:

- Extract accurate transcripts from social media videos for content analysis
- Get AI-powered video summaries with key insights and main points
- Retrieve engagement metrics including views, likes, and comments
- Analyze video comments with sorting and filtering options
- Create custom AI-driven analysis with tailored prompts and response fields
- Automate content analysis workflows across multiple platforms

**SocialKit Integration Documentation**

**Overview**
The SocialKit integration allows you to connect your SocialKit account with our workflow automation software, enabling you to automate video analysis tasks and extract valuable insights from social media content.

**Prerequisites**

- A SocialKit account
- Our workflow automation software account
- SocialKit API access key (found in your SocialKit account settings)

**Setup Instructions**

1. Log in to your SocialKit account and navigate to the "API Keys" section.
2. Generate or copy your existing API access key.
3. In our workflow automation software, go to the "Integrations" section and click on "SocialKit".
4. Enter the API access key from step 2 and click "Connect".

**Available Actions**

- **Get YouTube Transcript**: Extract accurate, timestamped transcripts from YouTube videos for content analysis and accessibility.

**Example Use Cases**

1. Content Analysis: Automatically extract transcripts from educational videos for text analysis and keyword extraction.
2. Accessibility: Generate transcripts for video content to improve accessibility for hearing-impaired audiences.
3. Content Repurposing: Extract transcripts to create blog posts, social media content, or documentation from video content.
4. SEO Optimization: Analyze video transcripts to identify keywords and improve video discoverability.
5. Research & Monitoring: Track and analyze video content across multiple channels for market research or competitive analysis.

**Troubleshooting Tips**

- Ensure that your API access key is entered correctly.
- Check the SocialKit API documentation for any rate limits or usage guidelines.
- Verify that the video URL is accessible and properly formatted.
- For YouTube videos, ensure the video has captions or transcripts available.

**FAQs**

Q: What video platforms are supported?
A: Currently, the integration supports YouTube videos, with TikTok and Instagram support coming soon.

Q: Are there any rate limits?
A: Yes, please refer to your SocialKit plan for specific rate limits. Free plans include 20 credits per month.

Q: Can I cache transcripts for faster access?
A: Yes, the API supports caching with customizable TTL (Time To Live) settings.

## Categories

- social-media
- analytics
- video-analysis

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name                   | Description                                                                                                            | Link                                      |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------- | ----------------------------------------- |
| Get YouTube Transcript | Extract accurate, timestamped transcripts from YouTube videos for content analysis, accessibility, and data extraction | [docs](actions/get_youtube_transcript.md) |
