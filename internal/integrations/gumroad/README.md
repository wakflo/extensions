# Gumroad Integration

## Description

Integrate your Gumroad store with our workflow automation software to streamline e-commerce operations. Automatically sync Gumroad products, sales, subscriptions, and customer data with your workflows, eliminating manual data entry and reducing errors. Use this integration to:

- Track sales and revenue in real-time
- Automate customer follow-ups based on purchase behavior
- Generate comprehensive sales reports using synced Gumroad data
- Enhance customer support by accessing purchase history information
- Trigger actions when new sales, refunds, or subscriptions occur

**Gumroad Integration Documentation**

**Overview**
The Gumroad integration allows you to seamlessly connect your Gumroad account with our workflow automation software, enabling you to automate tasks and streamline your e-commerce workflow.

**Prerequisites**

- A Gumroad account
- Our workflow automation software account
- Gumroad API access token (found in your Gumroad settings)

**Setup Instructions**

1. Log in to your Gumroad account and navigate to the "Settings" page.
2. Go to the "Advanced" tab and find your API token or create a new one.
3. In our workflow automation software, go to the "Integrations" section and click on "Gumroad".
4. Enter the API token from step 2 and click "Connect".

**Available Triggers**

- **Sale Created**: Triggered when a new sale is made in your Gumroad store.
- **Refund Created**: Triggered when a refund is processed for a purchase.
- **Dispute Created**: Triggered when a payment dispute is created.
- **Subscription Created**: Triggered when a new subscription is started.
- **Subscription Cancelled**: Triggered when a subscription is cancelled.

**Available Actions**

- **Get Product**: Retrieve details of a specific product from your Gumroad store.
- **List Products**: Get a list of all products in your Gumroad store.
- **Create Product**: Create a new product in your Gumroad store.
- **Update Product**: Update details of an existing product in your Gumroad store.
- **Get Sale**: Retrieve details of a specific sale.
- **List Sales**: Get a list of sales from your Gumroad store.
- **Get Subscription**: Retrieve details of a specific subscription.
- **List Subscriptions**: Get a list of all active subscriptions.
- **Get Customer**: Retrieve details of a specific customer.
- **List Customers**: Get a list of all customers who have purchased from your store.

**Example Use Cases**

1. Automate customer onboarding: When a new sale is made, add the customer to your CRM, send a welcome email, and create a support ticket if needed.
2. Track sales performance: Create dashboards that automatically update with your latest Gumroad sales data.
3. Subscription management: Automatically send reminders to customers before their subscriptions renew or follow up with customers who cancel.

**Troubleshooting Tips**

- Ensure that your API token is valid and has not expired.
- Check the Gumroad API documentation for any rate limits or usage guidelines.
- Make sure your Gumroad account has the necessary permissions for the actions you are trying to perform.

**FAQs**

Q: How often can I use the Gumroad API to retrieve data?
A: Gumroad has rate limits in place for their API. It's recommended to keep your requests below 100 per minute to avoid hitting these limits.

Q: Can I use this integration to process refunds?
A: This integration allows you to monitor when refunds occur, but direct refund processing must be done through the Gumroad interface.

## Categories

- ecommerce

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name               | Description                                                                                         | Link                                  |
| ------------------ | --------------------------------------------------------------------------------------------------- | ------------------------------------- |
| Get Product        | Retrieves detailed information about a specific product from your Gumroad store.                    | [docs](actions/get_product.md)        |
| List Products      | Retrieves a list of all products in your Gumroad store.                                             | [docs](actions/list_products.md)      |
| Create Product     | Creates a new product in your Gumroad store with specified details.                                 | [docs](actions/create_product.md)     |
| Update Product     | Updates an existing product in your Gumroad store with new information.                             | [docs](actions/update_product.md)     |
| Get Sale           | Retrieves detailed information about a specific sale from your Gumroad account.                     | [docs](actions/get_sale.md)           |
| List Sales         | Retrieves a list of sales from your Gumroad account with optional filtering.                        | [docs](actions/list_sales.md)         |
| Get Subscription   | Retrieves detailed information about a specific subscription from your Gumroad account.             | [docs](actions/get_subscription.md)   |
| List Subscriptions | Retrieves a list of active subscriptions from your Gumroad account with optional filtering.         | [docs](actions/list_subscriptions.md) |
| Get Customer       | Retrieves detailed information about a specific customer who has purchased from your Gumroad store. | [docs](actions/get_customer.md)       |
| List Customers     | Retrieves a list of customers who have purchased from your Gumroad store with optional filtering.   | [docs](actions/list_customers.md)     |

## Triggers

| Name                   | Description                                                                                                                               | Link                                       |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| Sale Created           | Triggers when a new sale is completed in your Gumroad store, allowing you to automate post-purchase workflows.                            | [docs](triggers/sale_created.md)           |
| Refund Created         | Triggers when a refund is processed for a purchase, allowing you to update inventory, accounting systems, or follow up with customers.    | [docs](triggers/refund_created.md)         |
| Dispute Created        | Triggers when a payment dispute is created for a purchase, enabling you to automatically respond to or track payment issues.              | [docs](triggers/dispute_created.md)        |
| Subscription Created   | Triggers when a customer starts a new subscription to one of your products, allowing you to implement onboarding processes.               | [docs](triggers/subscription_created.md)   |
| Subscription Cancelled | Triggers when a customer cancels their subscription to one of your products, allowing you to automate offboarding or retention workflows. | [docs](triggers/subscription_cancelled.md) |
