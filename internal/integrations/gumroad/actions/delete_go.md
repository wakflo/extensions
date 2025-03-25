# Delete Product

## Description

Permanently removes a product from your Gumroad store. This action cannot be undone, so use it with caution.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name             | Type    | Required | Description                                                            |
| ---------------- | ------- | -------- | ---------------------------------------------------------------------- |
| product-id       | String  | Yes      | The ID of the Gumroad product to delete.                               |
| confirm-deletion | Boolean | Yes      | Confirmation flag to prevent accidental deletion. Must be set to true. |

## Output

Returns an object confirming the deletion with the following details:

- Success status
- Product ID that was deleted
- Timestamp of deletion

## Example Usage

Use this action to:

- Remove outdated or discontinued products from your inventory
- Clean up test products that are no longer needed
- Remove products that violate platform policies
- Delete duplicate product listings
- Implement product lifecycle management workflows

## Notes

- **Warning**: This action permanently deletes the product. This operation cannot be undone.
- Consider using "Disable Product" instead if you only want to hide the product temporarily
- Deleting a product will not affect existing purchases or customer access to already purchased products
- However, no new purchases can be made for the deleted product
- Sale records for the product will remain in your account for reporting purposes
- It's recommended to export any product data you might need before deletion
