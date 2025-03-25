# List Subscribers

## Description

Retrieve a list of active subscribers from a specific Campaign Monitor list. This action supports pagination.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name | Display Name | Type | Description | Required | Default Value |
|------|--------------|------|-------------|----------|---------------|
| listId | List ID | String | The ID of the list for which to retrieve subscribers. | Yes | - |
| page | Page | Number | The page number to retrieve (for pagination). | No | 1 |


## Returns

This action returns an object containing the following fields:

| Field | Type | Description |
|-------|------|-------------|
| Results | Array | List of subscriber objects |
| ResultsOrderedBy | String | The field by which the results are ordered |
| OrderDirection | String | The direction in which the results are ordered |
| PageNumber | Number | The current page number |
| PageSize | Number | The number of subscribers per page |
| RecordsOnThisPage | Number | The number of subscribers on the current page |
| TotalNumberOfRecords | Number | The total number of subscribers matching the criteria |
| NumberOfPages | Number | The total number of pages |

Each subscriber object in the Results array contains:

| Field | Type | Description |
|-------|------|-------------|
| EmailAddress | String | The subscriber's email address |
| Name | String | The subscriber's name |
| Date | String | The date the subscriber was added or last updated |
| State | String | The subscriber's state (usually "Active") |
| CustomFields | Array | List of custom field objects with Key and Value properties |

## Example Usage

You can use this action to retrieve a list of subscribers for various purposes such as:

1. Generating reports on subscriber growth and engagement
2. Creating custom segments based on subscriber data
3. Exporting subscriber data to other systems
4. Performing analysis on subscriber demographics or preferences

To retrieve all subscribers across multiple pages, you can use this action in a loop, incrementing the page number with each iteration until you've processed all pages.