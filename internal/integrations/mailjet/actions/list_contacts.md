# List Contacts

Retrieves a list of contacts from your Mailjet account with optional filtering and pagination.

## Input Parameters

- **Limit** (optional): Maximum number of contacts to return. Default is 10, maximum is 1000.
- **Offset** (optional): Number of contacts to skip (for pagination). Default is 0.
- **Filter** (optional): Filter string in the format "PropertyName=Value". 
  
  Available filters include:
  - IsExcludedFromCampaigns=true/false
  - Name=value
  - Email=value@example.com

## Response

The response includes:

- **Count**: Number of contacts returned in this request
- **Total**: Total number of contacts matching your criteria
- **Data**: Array of contact objects with properties:
  - ContactID: Unique identifier
  - Email: Contact's email address
  - Name: Contact's name (if available)
  - IsExcludedFromCampaigns: Whether contact is excluded from campaigns
  - CreatedAt: Creation date
  - LastActivityAt: Last activity date
  - And other properties

## Example Usage

Retrieve the first 10 contacts:
- Leave all fields at default values

Retrieve contacts that are excluded from campaigns:
- Filter: IsExcludedFromCampaigns=true

Paginate through contacts (second page of 50):
- Limit: 50
- Offset: 50