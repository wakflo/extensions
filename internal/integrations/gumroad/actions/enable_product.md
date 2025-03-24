# Enable Product

## Description

Enables (publishes) a product in your Gumroad store, making it available for purchase. This is useful for activating or re-activating products that were previously disabled.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name       | Type   | Required | Description                              |
| ---------- | ------ | -------- | ---------------------------------------- |
| product-id | String | Yes      | The ID of the Gumroad product to enable. |

## Output

Returns an object containing the updated details of the product, including:

- Product ID
- Name
- Description
- URL and permalink
- Publication status (which will be set to true)
- Creation and update timestamps
- And other product metadata

## Example Usage

Use this action to:

- Launch new products at a specific time
- Re-enable seasonal or limited-time products
- Automate product availability based on inventory levels
- Publish products after they've passed review
- Schedule product releases as part of marketing campaigns

## Notes

When a product is enabled:

- It will appear in your public Gumroad store (if applicable)
- Direct links to the product will work
- Customers will be able to purchase the product
- The product can be disabled at any time using the "Disable Product" action
- You may want to notify your audience when re-enabling previously disabled products
