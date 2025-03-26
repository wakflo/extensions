# List Products

## Description

Retrieves a list of all products in your Gumroad store, with optional filtering and pagination.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name      | Type    | Required | Description                                         |
| --------- | ------- | -------- | --------------------------------------------------- |
| published | Boolean | No       | If true, only returns published products.           |
| limit     | Number  | No       | Maximum number of products to return (default: 10). |
| page      | Number  | No       | Page number for pagination (default: 1).            |
| category  | String  | No       | Filter products by category.                        |

## Output

Returns an array of product objects. Each product object contains:

- Name
- ID
- Description
- Price and formatted price
- Currency
- URL and permalink
- Publication status
- Thumbnail URL
- Creation date
- And other product metadata

## Example Usage

Use this action when you need to:

- Generate a catalog of products
- Create dashboards or reports
- Filter products by category or status
- Paginate through large product catalogs
