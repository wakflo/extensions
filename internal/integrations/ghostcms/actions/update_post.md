✏️ update_post

ID: update_post
Description: Updates an existing post (partial updates supported).


| Field              | Type     | Required | Description                                                              |
| ------------------ | -------- | -------- | ------------------------------------------------------------------------ |
| `post_id`          | string   | ✅        | The ID of the post to update.                                            |
| `title`            | string   | –        | New title.                                                               |
| `content`          | string   | –        | New content.                                                             |
| `content_format`   | string   | –        | `html` (default), `markdown`, `lexical` (applies if `content` provided). |
| `slug`             | string   | –        | New slug.                                                                |
| `status`           | string   | –        | New status `draft`, `published`, `scheduled`.                            |
| `featured`         | boolean? | –        | Whether post is featured (nullable boolean in payload).                  |
| `feature_image`    | string   | –        | New feature image URL.                                                   |
| `custom_excerpt`   | string   | –        | New custom excerpt.                                                      |
| `tags`             | string[] | –        | Replace tags (names).                                                    |
| `authors`          | string[] | –        | Replace authors (IDs).                                                   |
| `meta_title`       | string   | –        | SEO meta title.                                                          |
| `meta_description` | string   | –        | SEO meta description.                                                    |


{
  "post_id": "63f8d9a0e4b0d7001c3e3b1a",
  "title": "Updated Blog Post",
  "content": "<p>Updated HTML content.</p>",
  "content_format": "html",
  "status": "published",
  "featured": true,
  "tags": ["Changelog", "Release"],
  "authors": ["author-uuid-1"]
}


{
  "id": "63f8d9a0e4b0d7001c3e3b1a",
  "uuid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "title": "Updated Blog Post",
  "slug": "updated-blog-post",
  "html": "<p>This is the updated content.</p>",
  "status": "published",
  "featured": true,
  "visibility": "public",
  "created_at": "2024-02-24T10:30:00.000Z",
  "updated_at": "2024-02-24T15:45:00.000Z",
  "published_at": "2024-02-24T10:30:00.000Z",
  "url": "https://myblog.ghost.io/updated-blog-post/"
}
