# Retrieve Contact

## Description

Retrieve contacts by email in HubSpot. This action allows you to retrieve a contact in your HubSpot CRM by providing their email.

## Properties

| Name  | Type   | Required | Description      |
|-------|--------|----------|------------------|
| email | string | yes      | email of contact |



## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the details of the contact from HubSpot. The structure will include:

```json
{
  "id": "51",
  "properties": {
    "firstname": "John",
    "lastname": "Doe",
    "email": "johndoe@example.com",
    "phone": "+1234567890",
    "jobtitle": "tech executive",
    "company": "ACME"
  }
}
```

## Notes
- If the contact is not found, the response will be empty
- The response will include all default properties
- Some properties may be custom to your HubSpot instance