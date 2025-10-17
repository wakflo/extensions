# Create Post

Creates a new post in your **Ghost CMS** publication.

---

## 🧭 Overview

This action allows you to create a new post in your Ghost blog or publication using the **Ghost Admin API**.  
You can include formatted content (HTML, Markdown, or Lexical), metadata, SEO fields, authors, tags, and even trigger an email newsletter when the post is published.

---

## ⚙️ Action ID

---

## 🧩 Input Properties

| Field | Type | Required | Description |
|-------|------|-----------|-------------|
| **title** | string | ✅ Yes | The title of the post. |
| **content** | string | ✅ Yes | The post content (HTML, Markdown, or Lexical JSON). |
| **content_format** | string | ❌ No | Format of the content. Options: `html`, `markdown`, `lexical`. Default: `html`. |
| **slug** | string | ❌ No | Custom slug for the post URL. Automatically generated if omitted. |
| **status** | string | ✅ Yes | Post status: `draft`, `published`, or `scheduled`. |
| **featured** | boolean | ❌ No | Mark the post as featured. Default: `false`. |
| **feature_image** | string | ❌ No | URL to the post’s feature image. |
| **custom_excerpt** | string | ❌ No | Custom excerpt or summary. |
| **tags** | array of strings | ❌ No | Tags to assign to the post. Example: `["News", "Announcements"]`. |
| **authors** | array of strings | ❌ No | Author IDs to attribute the post to. |
| **meta_title** | string | ❌ No | SEO meta title. |
| **meta_description** | string | ❌ No | SEO meta description. |
| **og_title** | string | ❌ No | Open Graph title for social platforms. |
| **og_description** | string | ❌ No | Open Graph description for social sharing. |
| **twitter_title** | string | ❌ No | Twitter card title. |
| **twitter_description** | string | ❌ No | Twitter card description. |
| **codeinjection_head** | string | ❌ No | Custom code to inject in the `<head>` section of the post. |
| **codeinjection_foot** | string | ❌ No | Custom code to inject before the closing `</body>` tag. |
| **email_subject** | string | ❌ No | Subject line for the email newsletter. |
| **send_email_when_published** | boolean | ❌ No | Whether to send an email newsletter when the post is published. Default: `false`. |

---

## 🧾 Example Input

```json
{
  "title": "My New Blog Post",
  "content": "<p>This is the content of my blog post.</p>",
  "content_format": "html",
  "slug": "my-new-blog-post",
  "status": "draft",
  "featured": false,
  "feature_image": "https://cdn.example.com/post-banner.jpg",
  "custom_excerpt": "A quick overview of my new post",
  "tags": ["Announcement",]
  }

{
  "id": "63f8d9a0e4b0d7001c3e3b1a",
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "title": "My New Blog Post",
  "slug": "my-new-blog-post",
  "html": "<p>This is the content of my blog post.</p>",
  "status": "draft",
  "featured": false,
  "visibility": "public",
  "created_at": "2024-02-24T10:30:00.000Z",
  "updated_at": "2024-02-24T10:30:00.000Z",
  "published_at": null,
  "url": "https://myblog.ghost.io/my-new-blog-post/"
}
