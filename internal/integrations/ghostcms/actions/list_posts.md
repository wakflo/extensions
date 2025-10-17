# List Posts

Lists all posts from your **Ghost CMS** publication with filtering, sorting, and pagination options.

---

## 🧭 Overview

The **List Posts** action retrieves posts from your Ghost publication using the **Ghost Admin API**.  
It supports filtering by status, sorting, pagination, and selective field inclusion.  
This is ideal for syncing post data, creating dashboards, or dynamically loading Ghost content into other systems.

---

## ⚙️ Action ID

---

## 🧩 Input Properties

| Field | Type | Required | Description |
|-------|------|-----------|-------------|
| **limit** | integer | ❌ No | Number of posts to retrieve per page. Default: `15`. |
| **page** | integer | ❌ No | Page number for pagination. Default: `1`. |
| **status** | string | ❌ No | Filter posts by status. Options: `all`, `published`, `draft`, `scheduled`. Default: `all`. |
| **filter** | string | ❌ No | Custom [NQL](https://ghost.org/docs/content-api/#filtering) filter query (e.g. `featured:true`, `tag:news`). |
| **order** | string | ❌ No | Sort order for posts. Options:<br> - `published_at desc`<br> - `published_at asc`<br> - `created_at desc`<br> - `created_at asc`<br> - `updated_at desc`<br> - `updated_at asc`.<br>Default: `published_at desc`. |
| **include** | string | ❌ No | Comma-separated list of related objects to include (e.g. `authors,tags`). |
| **fields** | string | ❌ No | Comma-separated list of fields to return (e.g. `id,title,slug`). |
| **formats** | string | ❌ No | Comma-separated list of content formats to include. Options: `html`, `plaintext`. Default: `html`. |

---

## 🧾 Example Input

```json
{
  "limit": 10,
  "page": 1,
  "status": "published",
  "filter": "tag:news+featured:true",
  "order": "published_at desc",
  "include": "authors,tags",
  "fields": "id,title,slug,created_at,published_at",
  "formats": "html"
}
{
  "posts": [
    {
      "id": "63f8d9a0e4b0d7001c3e3b1a",
      "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "title": "Sample Blog Post",
      "slug": "sample-blog-post",
      "html": "<p>This is a sample blog post content.</p>",
      "status": "published",
      "featured": false,
      "visibility": "public",
      "created_at": "2024-02-24T10:30:00.000Z",
      "updated_at": "2024-02-24T10:30:00.000Z",
      "published_at": "2024-02-24T10:30:00.000Z",
      "url": "https://myblog.ghost.io/sample-blog-post/"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "limit": 15,
      "pages": 1,
      "total": 1,
      "next": null,
      "prev": null
    }
  }
}
