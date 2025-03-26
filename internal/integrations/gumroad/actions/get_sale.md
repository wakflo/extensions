# Get Sale

## Description

Retrieves detailed information about a specific sale from your Gumroad account.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name    | Type   | Required | Description                 |
| ------- | ------ | -------- | --------------------------- |
| sale-id | String | Yes      | The ID of the Gumroad sale. |

## Output

Returns an object containing the details of the sale, including:

- Sale ID
- Creation timestamp
- Product ID and name
- Price and currency
- Customer email and name
- Order number
- Subscription information (if applicable)
- Refund and dispute status
- License key (if applicable)
- Customer location information
- Custom fields provided during checkout
- Payment method details
- Affiliate information (if applicable)
- Discount code used (if applicable)

## Example Usage

Use this action to:

- Access detailed purchase information for customer support
- Verify purchase details for license validation
- Track specific sales for accounting or analytics
- Check subscription status for a purchase
- Review custom field information submitted during checkout
- Generate receipts or invoices based on sale data
- Track payment and fulfillment status
- Determine if a purchase was made with a discount code

## Notes

- This action retrieves a single sale by its ID; use the "List Sales" action if you need to retrieve multiple sales
- The returned data includes sensitive customer information, so ensure you handle it according to relevant privacy regulations
- Sale IDs can be found in the Gumroad dashboard or via webhook notifications when new sales occur
