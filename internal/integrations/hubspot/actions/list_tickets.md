# List Tickets

## Description

List tickets in HubSpot. This action allows you to list tickets in your HubSpot CRM by providing their details.

## Properties

| Name     | Type   | Required | Description              |
|----------|--------|----------|--------------------------|
| limit    | int    | no       | limit of tickets to show |



## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the list of tickets from HubSpot. The structure will include:

```json
{
  "results": [
    {
      "id": "51",
      "properties": {
        "subject": "Ticket 1",
        "status": "Open",
        "createdate": "2022-01-01T00:00:00Z"
      }
    },
    {
      "id": "52",
      "properties": {
        "subject": "Ticket 2",
        "status": "Closed",
        "createdate": "2022-01-02T00:00:00Z"
      }
    }
  ],
  "paging": {
    "next": {
      "after": "NTI=",
      "link": "https://api.hubapi.com/crm/v3/objects/tickets?after=NTI="
    }
  }
}
```

## Notes

- The `id` field is the unique identifier of the ticket in HubSpot
- The `subject` field is the title of the ticket
- The `status` field is the current status of the ticket
- Possible values are `Open`, `Closed`, `Pending`, `Resolved`, `Waiting on Customer`, `Waiting on Third Party`, `Custom`