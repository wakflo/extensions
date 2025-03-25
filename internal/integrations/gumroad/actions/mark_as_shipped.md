# Mark As Shipped

## Description

Marks a sale of a physical product as shipped in your Gumroad account. This updates the delivery status for both you and your customer, and triggers a shipping notification email to the customer.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name             | Type    | Required | Description                                                     |
| ---------------- | ------- | -------- | --------------------------------------------------------------- |
| sale-id          | String  | Yes      | The ID of the Gumroad sale to mark as shipped.                  |
| tracking-number  | String  | No       | The tracking number for the shipment.                           |
| tracking-url     | String  | No       | The URL where the customer can track their shipment.            |
| carrier          | String  | No       | The name of the shipping carrier (e.g., USPS, FedEx, DHL).      |
| notify-customer  | Boolean | No       | Whether to notify the customer by email (default: true).        |
| additional-notes | String  | No       | Additional notes to include in the shipping notification email. |

## Output

Returns an object containing the updated details of the sale, including:

- Sale ID
- Product details
- Customer information
- Shipping status (now marked as shipped)
- Tracking information if provided
- Timestamp of when the shipment was marked
- Whether a notification was sent to the customer

## Example Usage

Use this action to:

- Update shipping status when an order is fulfilled
- Integrate with your shipping or fulfillment system
- Provide customers with tracking information automatically
- Keep your Gumroad sales records up to date
- Automate the notification process for shipped orders
- Complete your order fulfillment workflow

## Notes

- This action only applies to physical products that require shipping
- Customers will receive an email notification by default (unless opted out via the notify-customer parameter)
- Adding tracking information is highly recommended to improve customer experience
- You can customize the notification email by adding additional notes
- For sales with multiple items, this marks the entire order as shipped
- Sales already marked as shipped can be updated with new tracking information
