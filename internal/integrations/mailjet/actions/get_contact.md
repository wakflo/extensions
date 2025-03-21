# Get Contact

Retrieves detailed information about a specific contact from your Mailjet account.

## Input Parameters

You must provide either:
- **Contact ID**: The numeric ID of the contact to retrieve, or
- **Email**: The email address of the contact to retrieve

## Response

The response includes:

- **Count**: Number of contacts retrieved (1 if found, 0 if not found)
- **Total**: Total count (same as Count)
- **Data**: Array containing the contact's information with properties:
  - ContactID: Unique identifier for the contact
  - Email: Contact's email address
  - Name: Contact's name (if available)
  - IsExcludedFromCampaigns: Whether contact is excluded from campaigns
  - CreatedAt: Creation date and time
  - DeliveredCount: Number of emails delivered to this contact
  - IsOptInPending: Whether the contact has pending opt-in confirmation
  - IsSpamComplaining: Whether the contact has marked emails as spam
  - LastActivityAt: When the contact last interacted with emails
  - And other properties

## Error Handling

- If the contact is not found, the action will return an error.
- If neither Contact ID nor Email is provided, the action will return an error.