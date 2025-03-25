# Create Contact

Creates a new contact in your Mailjet account.

## Input Parameters

- **Email** (required): Email address of the contact to create.
- **Name** (optional): Name of the contact.
- **Exclude from Campaigns** (optional): When set to true, the contact will not receive marketing campaigns. Default is false.

## Response

The response includes:

- **Count**: Number of contacts created (always 1)
- **Total**: Total count (always 1)
- **Data**: Array containing the created contact with properties:
  - ContactID: Unique identifier for the contact
  - Email: Contact's email address
  - Name: Contact's name (if provided)
  - IsExcludedFromCampaigns: Whether contact is excluded from campaigns
  - CreatedAt: Creation date and time
  - And other properties

## Notes

- If the contact already exists in your Mailjet account, the API will return the existing contact information.
- New contacts are automatically assigned a ContactID by Mailjet.
