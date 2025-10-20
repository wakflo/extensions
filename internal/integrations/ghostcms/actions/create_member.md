# Create Member

Creates a new member in your **Ghost CMS** publication.

---

## 🧭 Overview

This action allows you to create a new member record in Ghost using the Admin API.  
You can specify email, name, notes, labels, newsletters, and subscription options.  
It’s useful for onboarding new subscribers or syncing external contact lists.

---

## ⚙️ Action ID

---

## 🧩 Input Properties

| Field | Type | Required | Description |
|-------|------|-----------|-------------|
| **email** | string | ✅ Yes | The member’s email address. This is required and must be unique. |
| **name** | string | ❌ No | The full name of the member. |
| **note** | string | ❌ No | An internal note about the member. Not visible publicly. |
| **labels** | array of strings | ❌ No | Labels (tags) to assign to the member, e.g. `["Premium", "Newsletter"]`. |
| **newsletters** | array of strings | ❌ No | IDs of newsletters the member should be subscribed to. |
| **subscribed** | boolean | ❌ No | Whether the member is subscribed to newsletters. Defaults to `true`. |
| **comped** | boolean | ❌ No | Whether to give the member a complimentary paid subscription. Defaults to `false`. |

---

## 🧾 Example Input

```json
{
  "email": "john.doe@example.com",
  "name": "John Doe",
  "note": "VIP customer",
  "labels": ["Premium", "Active"],
  "newsletters": ["newsletter-uuid-1"],
  "subscribed": true,
  "comped": false
}
{
  "id": "63f8d9a0e4b0d7001c3e3b1b",
  "uuid": "b2c3d4e5-f6a7-8901-bcde-f23456789012",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "note": "VIP customer",
  "status": "free",
  "subscribed": true,
  "created_at": "2024-02-24T10:30:00.000Z",
  "updated_at": "2024-02-24T10:30:00.000Z",
  "labels": [
    {
      "id": "1",
      "name": "Premium",
      "slug": "premium"
    }
  ],
  "email_count": 0,
  "email_opened_count": 0,
  "email_open_rate": 0
}
