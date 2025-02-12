# Stripe Integration

## Description

Integrate Stripe with your workflow automation software to seamlessly process payments and manage transactions within your automated workflows. With this integration, you can:

* Automatically trigger payment processing upon workflow completion
* Update customer information and subscription status in real-time
* Generate invoices and send notifications based on workflow events
* Track payment history and reconcile transactions directly from the workflow interface

Streamline your payment processing and enhance the overall customer experience by automating Stripe integrations with your workflows.

**Prerequisites**

* A Stripe account with API keys enabled
* A workflow automation software account with the necessary permissions to integrate with Stripe

**Setup**

1. In your workflow automation software, navigate to the integrations or marketplace section.
2. Search for Stripe and click on the integration tile.
3. Click "Connect" to initiate the authentication process.
4. You will be redirected to the Stripe login page. Enter your Stripe account credentials to authorize the connection.
5. Once authorized, you will be returned to your workflow automation software with a confirmation message indicating that the integration is set up.

**Configuration**

1. In the integration settings, configure the following:
	* **API Key**: Enter your Stripe API key (available in the Stripe dashboard under "Developers" > "API keys").
	* **Webhook Secret**: Enter a unique secret key to secure webhook communications between your workflow automation software and Stripe.
2. Map Stripe objects to your workflow automation software's objects by selecting the corresponding fields.

**Features**

1. **Payment Processing**: Automate payment processing for invoices, subscriptions, or one-time payments using Stripe's APIs.
2. **Customer Synchronization**: Synchronize customer information between your workflow automation software and Stripe, ensuring that both systems have up-to-date customer data.
3. **Webhooks**: Receive notifications from Stripe when events occur (e.g., payment succeeded, failed, or refunded) to trigger custom workflows in your workflow automation software.

**Troubleshooting**

1. If you encounter issues with the integration, check the Stripe dashboard for any errors or warnings related to API key authentication.
2. Verify that the webhook secret is correctly configured on both sides of the integration.

**Security**
The Stripe integration uses secure protocols and encryption to protect sensitive payment information. Ensure that your workflow automation software's security settings are also properly configured to maintain data integrity and confidentiality.

**Support**
For assistance with the Stripe integration, please contact our support team or refer to Stripe's official documentation for more information on their APIs and integrations.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>## Actions

- **Create Customer**: Create a new customer in your CRM system by providing required details such as name, email, phone number, and other relevant information. This integration action allows you to automate the process of creating new customers, reducing manual errors and increasing efficiency. ([Documentation]([Create Customer](actions/create_customer.md)))

- **Create Invoice**: Create Invoice: Automatically generates and sends professional-looking invoices to customers based on predefined templates and payment terms, streamlining your accounting process and ensuring timely payments. ([Documentation]([Create Invoice](actions/create_invoice.md)))

- **Search Customer**: Searches for a customer by their name, email, or phone number in your CRM system and retrieves relevant information such as contact details, order history, and account status. ([Documentation]([Search Customer](actions/search_customer.md)))

## Triggers

- **New Customer**: Triggered when a new customer is created in your CRM or database, this integration allows you to automate workflows and tasks immediately after a new customer is added, streamlining your sales and marketing processes. ([Documentation]([New Customer](triggers/new_customer.md)))

