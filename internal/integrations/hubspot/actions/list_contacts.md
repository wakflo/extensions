# List Contacts

## Description

Retrieve a list of contacts from HubSpot based on optional filtering criteria. This action allows you to fetch multiple contacts at once, with pagination support for handling large contact databases.

## Properties

| Name       | Type   | Required | Description                                                                         |
|------------|--------|----------|-------------------------------------------------------------------------------------|
| limit      | number | No       | Maximum number of contacts to return (max: 100, default: 20)                        |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs a list of contacts from HubSpot with pagination information. The structure will include:

```json
{
  "results": [
    {
      "id": "51",
      "properties": {
        "firstname": "John",
        "lastname": "Doe",
        "email": "john.doe@example.com"
      }
    },
    {
      "id": "52",
      "properties": {
        "firstname": "Jane",
        "lastname": "Smith",
        "email": "jane.smith@example.com"
      }
    }
  ],
  "paging": {
    "next": {
      "after": "NTI=",
      "link": "https://api.hubapi.com/crm/v3/objects/contacts?after=NTI="
    }
  }
}
```

## Notes

- To retrieve more than 100 contacts, you'll need to use pagination with the "after" parameter
- The "after" value is provided in the response's "paging.next.after" property
- If you don't specify any properties, the response will include all default properties
- You can optimize performance by requesting only the specific properties you need
- Some properties may be custom to your HubSpot instance