# Get Product

## Description

Retrieves detailed information about a specific product from your Gumroad store.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name       | Type   | Required | Description                    |
| ---------- | ------ | -------- | ------------------------------ |
| product-id | String | Yes      | The ID of the Gumroad product. |

## Output

Returns a detailed object containing all product information including:

- Name
- Description
- Price and currency
- URL and permalink
- Publication status
- Thumbnail URL
- Creation date
- Custom fields
- File information (if applicable)
- Custom receipt and summary text
- Categories
- And more product metadata

## Example Usage

Use this action when you need to retrieve product information to:

- Display product details in your application
- Check pricing information
- Verify product availability
- Reference product data in automation workflows
