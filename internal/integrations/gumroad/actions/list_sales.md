# List Sales

## Description

Retrieves a list of sales from your Gumroad account with optional filtering.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name            | Type    | Required | Description                                     |
| --------------- | ------- | -------- | ----------------------------------------------- |
| product-id      | String  | No       | Filter sales by product ID.                     |
| email           | String  | No       | Filter sales by customer email.                 |
| after           | String  | No       | Get sales after this date (YYYY-MM-DD format).  |
| before          | String  | No       | Get sales before this date (YYYY-MM-DD format). |
| page            | Number  | No       | Page number for pagination (default: 1).        |
| limit           | Number  | No       | Number of sales per page (default: 10).         |
| refunded        | Boolean | No       | Filter by refund status.                        |
| disputed        | Boolean | No       | Filter by dispute status.                       |
| subscription-id | String  | No       | Filter sales by subscription ID.                |

## Output

Returns an array of sale objects. Each sale object contains:

- Sale ID
- Creation timestamp
- Product ID and name
- Price and currency
- Customer email
- Order number
- Refund and dispute status
- And other sale-specific details

## Example Usage

Use this action to:

- Generate reports on recent sales activity
- Filter sales by specific products or customers
- Retrieve sales data for accounting or analytics
- Find sales within a specific date range
- Identify refunded or disputed transactions
- Track subscription-related purchases
