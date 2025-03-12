# Search Records

## Description

Searches for records in a specified Zoho CRM module based on search criteria.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name       | Type   | Required | Default | Description                                                                      |
|------------|--------|----------|---------|----------------------------------------------------------------------------------|
| module     | string | Yes      | -       | The Zoho CRM module to search in (e.g., Leads, Contacts, Accounts)               |
| searchType | string | Yes      | -       | Type of search to perform (word or criteria)                                     |
| criteria   | string | Yes      | -       | The search criteria, depending on the search type                                |
| page       | number | No       | 1       | The page number for pagination (1-based)                                         |
| perPage    | number | No       | 100     | The number of records to return per page (max 200)                               |

## Search Types

1. **Word**: Simple search for records containing the specified word.
2. **Criteria**: Advanced search using Zoho CRM's criteria format.

## Criteria Format Examples

For criteria search, use the following format:

- Single criterion: `(field:operator:value)`
- Multiple criteria with AND: `(field1:operator:value1)AND(field2:operator:value2)`
- Multiple criteria with OR: `(field1:operator:value1)OR(field2:operator:value2)`

Available operators:
- equals
- starts_with
- ends_with
- contains
- not_equals
- greater_than
- less_than
- greater_equals
- less_equals

## Sample JSON Data Format

```json
{
  "module": "Leads",
  "searchType": "criteria",
  "criteria": "(Last_Name:equals:Smith)OR(Email:contains:example.com)",
  "page": 1,
  "perPage": 100
}