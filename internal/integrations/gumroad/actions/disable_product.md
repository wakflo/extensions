# Disable Product

## Description

Disables (unpublishes) a product in your Gumroad store, making it unavailable for purchase. This is useful for temporarily removing products from sale without deleting them completely.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name       | Type   | Required | Description                               |
| ---------- | ------ | -------- | ----------------------------------------- |
| product-id | String | Yes      | The ID of the Gumroad product to disable. |

## Output

Returns an object containing the updated details of the product, including:

- Product ID
- Name
- Description
- URL and permalink
- Publication status (which will be set to false)
- Creation and update timestamps
- And other product metadata

## Example Usage

Use this action to:

- Remove seasonal or limited-time products from your store
- Take products offline during maintenance or updates
- Unpublish products that are out of stock
- Hide products that are pending approval or review
- Manage product availability in automated workflows

## Notes

When a product is disabled:

- Existing customers who have already purchased the product will still have access
- The product will not appear in your public Gumroad store
- Direct links to the product will show a "not available" message
- The product can be re-enabled at any time using the "Update Product" action
