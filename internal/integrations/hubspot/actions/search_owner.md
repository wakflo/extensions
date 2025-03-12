# Search Owner

## Description

Search for an owner in HubSpot by providing their email.
## Properties

| Name  | Type   | Required | Description    |
|-------|--------|----------|----------------|
| email | string | yes      | email of owner |



## Details

- **Type**: sdkcore.ActionTypeNormal

## Output

This action outputs the details of the owner from HubSpot. The structure will include:

```json
{
  "id": "51",
  "properties": {
    "firstname": "John",
    "lastname": "Doe"
  }
}
```

## Notes
- If the owner is not found, the response will be empty
- The response will include all default properties