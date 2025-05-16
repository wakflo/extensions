# SendOwl Integration

## Description

Integrate SendOwl's digital product delivery platform with your workflow automation. Connect your digital product sales data to automate post-purchase workflows, monitor sales performance, and enhance customer experience. Use this integration to:

- Monitor sales and automatically process order completions
- Track product performance and popularity metrics
- Automate customer communication based on purchase events
- Sync order data with your CRM, email marketing, or accounting systems
- Create real-time reporting and analytics dashboards

**SendOwl Integration Documentation**

**Overview**
The SendOwl integration allows you to connect your SendOwl account with our workflow automation software, enabling you to automate tasks and streamline your workflow for digital product sales.

**Prerequisites**

- A SendOwl account
- Our workflow automation software account
- SendOwl API key and API secret (found in your SendOwl account settings)

**Setup Instructions**

1. Log in to your SendOwl account and navigate to the "Settings" tab.
2. Go to "API" section and generate or copy your existing API key and secret.
3. In our workflow automation software, go to the "Integrations" section and click on "SendOwl".
4. Enter the API key and secret from step 2 and click "Connect".

**Available Actions**

- **List Products**: Retrieve a list of all your digital products from SendOwl.
- **Get Product**: Get detailed information about a specific product.
- **List Orders**: Retrieve a list of orders with optional filtering capabilities.
- **Get Order**: Get detailed information about a specific order.

**Available Triggers**

- **Order Completed**: Trigger workflows when a new order is completed.
- **Product Updated**: Trigger workflows when a product is updated.

**Example Use Cases**

1. Customer Onboarding: When a new order is completed, automatically send a welcome email and add the customer to your CRM.
2. Sales Reporting: Automatically compile daily or weekly sales reports based on completed orders.
3. Inventory Management: When product stock falls below a threshold, trigger a workflow to notify your team.
4. Customer Segmentation: Categorize customers based on their purchase history for targeted marketing campaigns.

**Troubleshooting Tips**

- Ensure that your API key and secret are entered correctly.
- Check the SendOwl API documentation for any rate limits or usage guidelines.
- Verify that the required permissions are enabled for your API key.

**FAQs**

Q: How often does the Order Completed trigger check for new orders?
A: By default, the trigger checks for new orders every 5 minutes.

Q: Can I filter orders by date range?
A: Yes, the List Orders action allows you to filter by date range, status, and other criteria.

## Categories

- ecommerce
- sales

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name          | Description                                                                                                                                                                           | Link                             |
| ------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| List Products | Retrieve a comprehensive list of all products from your SendOwl account, allowing you to access product information for reporting, analysis, or to trigger subsequent workflow steps. | [docs](actions/list_products.md) |
| Get Product   | Retrieve detailed information about a specific product from your SendOwl account using the product ID.                                                                                | [docs](actions/get_product.md)   |
| List Orders   | Retrieve a list of orders from your SendOwl account with optional filtering by date range, status, and other criteria.                                                                | [docs](actions/list_orders.md)   |
| Get Order     | Retrieve detailed information about a specific order from your SendOwl account using the order ID.                                                                                    | [docs](actions/get_order.md)     |

## Triggers

| Name            | Description                                                                                                                                                                                                             | Link                                |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| Order Completed | Automatically trigger workflows when a new order is marked as completed in your SendOwl account. This allows you to automate post-purchase processes such as customer onboarding, fulfillment, or marketing follow-ups. | [docs](triggers/order_completed.md) |
| Product Updated | Automatically trigger workflows when a product in your SendOwl account is updated. This allows you to maintain synchronized product information across systems or notify team members about product changes.            | [docs](triggers/product_updated.md) |
