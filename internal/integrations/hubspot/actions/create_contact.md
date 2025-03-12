# Create Contact

## Description

Create a new contact in HubSpot with the specified information. This action allows you to add new contacts to your HubSpot CRM by providing their details.

## Properties

| Name       | Type   | Required | Description                                                                       |
|------------|--------|----------|-----------------------------------------------------------------------------------|
| email      | string | Yes      | Email address of the contact                                                      |
| firstname  | string | No       | First name of the contact                                                         |
| lastname   | string | No       | Last name of the contact                                                          |
| phone      | string | No       | Phone number of the contact                                                       |
| company    | string | No       | Company name the contact works for                                                |
| jobtitle   | string | No       | Job title of the contact                                                          |
| website    | string | No       | Website URL associated with the contact                                           |
| address    | string | No       | Street address of the contact                                                     |
| city       | string | No       | City of the contact                                                               |
| state      | string | No       | State or region of the contact                                                    |
| zipcode    | string | No       | Zip or postal code of the contact                                                 |
| country    | string | No       | Country of the contact                                                            |
| properties | string | No       | JSON object with additional properties in 'property_name': 'value' format         |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the created contact details from HubSpot. The structure will include:

```json
{
  "id": "51",
  "properties": {
    "createdate": "2023-04-20T10:32:45.678Z",
    "email": "john.doe@example.com",
    "firstname": "John",
    "hs_object_id": "51",
    "lastmodifieddate": "2023-04-20T10:32:45.678Z",
    "lastname": "Doe"
  },
  "createdAt": "2023-04-20T10:32:45.678Z",
  "updatedAt": "2023-04-20T10:32:45.678Z"
}
```

## Notes

- Email is the only required field for creating a contact in HubSpot
- HubSpot may deduplicate contacts based on email address
- The properties field allows you to set custom or additional properties not covered by the standard fields
- For the properties field, use a JSON object with the format: `{"property_name1": "value1", "property_name2": "value2"}`
- Some properties might be specific to your HubSpot instance based on your custom property configurations