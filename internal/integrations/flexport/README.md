# Flexport Integration

## Description

Seamlessly integrate with Flexport to automate your international shipping and logistics processes. With this integration, you can easily track shipments, update order status, and retrieve shipment data in real-time, streamlining your operations and reducing manual errors. Enjoy increased visibility, improved communication, and enhanced decision-making capabilities with our seamless connection to the Flexport platform.

**Flexport Integration Documentation**

**Overview**
The Flexport integration allows you to automate workflows with Flexport, a leading logistics platform. This integration enables seamless data exchange between your workflow automation software and Flexport, streamlining your shipping operations.

**Prerequisites**

* A Flexport account
* A workflow automation software account
* Basic understanding of API concepts

**Setup**

1. **Create a Flexport API Key**: Log in to your Flexport account and navigate to the "Settings" > "API Keys" section. Create a new API key or use an existing one.
2. **Configure the Integration**: In your workflow automation software, go to the "Integrations" or "Connections" page and search for "Flexport". Click on the integration tile to begin setup.
3. **Enter Flexport API Credentials**: Enter the API key created in step 1 and any additional required credentials (e.g., client ID, secret key).
4. **Map Flexport Endpoints**: Map relevant Flexport endpoints to your workflow automation software's actions or triggers. Common use cases include:
	* Creating shipments
	* Updating shipment status
	* Retrieving shipment information

**API Endpoints**

The following Flexport API endpoints are supported:

* `shipments`: Create, read, update, and delete shipments
* `shipment_status`: Update shipment status
* `shipments/{id}`: Retrieve a specific shipment by ID

**Example Use Cases**

1. **Automate Shipment Creation**: When a new order is created in your workflow automation software, use the Flexport API to create a corresponding shipment.
2. **Update Shipment Status**: When a shipment status changes (e.g., from "in transit" to "delivered"), update the status using the Flexport API.

**Troubleshooting**

* Check the Flexport API documentation for any rate limits or usage guidelines
* Verify that your workflow automation software's integration is properly configured and authenticated
* Monitor logs for errors or issues related to the integration

**FAQs**

Q: What are the benefits of integrating Flexport with my workflow automation software?
A: Streamlined shipping operations, reduced manual errors, and improved visibility into shipment status.

Q: Can I use this integration to automate multiple workflows?
A: Yes, you can map different Flexport endpoints to various actions or triggers in your workflow automation software.

**Conclusion**
The Flexport integration enables seamless automation of your shipping workflows with Flexport. By following the setup instructions and mapping relevant endpoints, you can streamline your operations and reduce manual errors.

## Categories

- app


## Authors

- Wakflo <integrations@wakflo.com>

