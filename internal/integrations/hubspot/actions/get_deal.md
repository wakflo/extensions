# Get deal

## Description

Retrieve a deal in HubSpot by providing its ID.
## Properties

| Name   | Type   | Required | Description |
|--------|--------|----------|-------------|
| dealID | string | No       | the deal ID |

## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the details of the deal from HubSpot. The structure will include:

```json
{
    "deals": [
        {
        "id": "51",
        "properties": {
            "dealname": "ACME Deal",
            "dealstage": "closedwon",
            "amount": "1000",
            "closedate": "2021-01-01T00:00:00Z",
            "pipeline": "default",
            "hubspot_owner_id": "1",
            "createdate": "2021-01-01T00:00:00Z",
            "hs_lastmodifieddate": "2021-01-01T00:00:00Z"
        }
        }
    ],
    "paging": {
        "next": {
        "after": "100"
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