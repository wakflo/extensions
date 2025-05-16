# List Orders

## Description

Retrieve a list of orders from your SendOwl account with optional filtering by date range, status, and other criteria.

## Properties

| Name        | Type   | Required | Description                                        |
| ----------- | ------ | -------- | -------------------------------------------------- |
| page        | number | no       | Page number for pagination, defaults to 1          |
| limit       | number | no       | Number of orders per page, defaults to 50          |
| status      | string | no       | Filter orders by status (completed, pending, etc.) |
| start_date  | string | no       | Filter orders from this date (format: YYYY-MM-DD)  |
| end_date    | string | no       | Filter orders until this date (format: YYYY-MM-DD) |
| product_id  | string | no       | Filter orders by specific product ID               |
| buyer_email | string | no       | Filter orders by buyer's email address             |

## Details

- **Type**: sdkcore.ActionTypeNormal
