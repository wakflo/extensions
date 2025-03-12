# List Contacts

## Description

Retrieve a list of contacts from Freshworks CRM with options to filter and sort the results.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name        | Type   | Required | Description                                                      |
|-------------|--------|----------|------------------------------------------------------------------|
| page        | number | No       | Page number for pagination (default: 1)                          |
| per_page    | number | No       | Number of results per page (default: 25, max: 100)               |
| search_term | string | No       | Term to search for in contact fields                             |
| sort_by     | string | No       | Field to sort results by (e.g., first_name, last_name, created_at) |
| filter_by   | string | No       | JSON string with filter criteria (e.g., {"updated_at":{"gt":"2021-01-01"}}) |