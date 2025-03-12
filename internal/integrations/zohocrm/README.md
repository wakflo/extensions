# Zoho CRM Integration

## Description

Integrate Zoho CRM with your workflows to automate sales processes, customer data management, and more. Connect your Zoho CRM account to sync contacts, leads, accounts, and other modules with your business applications, eliminating manual data entry and ensuring your customer information stays up-to-date across all platforms. Use this integration to:

* Create and update records in Zoho CRM automatically
* Trigger workflows when new records are created or modified
* Search for and retrieve customer information across modules
* Automate sales processes and customer communication
* Keep customer data synchronized across different business systems

**Zoho CRM Integration Documentation**

**Overview**

The Zoho CRM integration connects your Zoho CRM account with our workflow automation software, allowing you to automate CRM-related tasks and streamline your sales and customer management processes.

**Prerequisites**

* A Zoho CRM account (Professional or Enterprise plan recommended for API access)
* A Zoho API client (created in the Zoho Developer Console)
* Our workflow automation software account

**Setup Instructions**

1. Log in to your Zoho Developer Console at https://api-console.zoho.com/
2. Create a Self-Client application
3. Set the scope to `ZohoCRM.modules.ALL,ZohoCRM.settings.ALL,ZohoCRM.notifications.ALL`
4. Generate a code, then exchange it for refresh and access tokens
5. In our workflow automation software, connect to Zoho CRM using the OAuth2 authentication

**Available Modules**

Zoho CRM offers multiple modules that you can interact with:

* Leads
* Contacts
* Accounts
* Deals
* Tasks
* Campaigns
* Products
* Cases
* Solutions
* And other custom modules

**Example Use Cases**

1. Create a new lead in Zoho CRM when a form is submitted on your website
2. Update a deal's stage when a payment is processed
3. Create a task in Zoho CRM when a support ticket is opened
4. Send personalized emails to contacts based on their activity
5. Create invoices in your accounting software when a deal is closed in Zoho CRM

**Troubleshooting Tips**

* Ensure your OAuth tokens have the necessary permissions (scopes)
* Check that the module names and field names match exactly as they appear in Zoho CRM
* Verify that mandatory fields are included when creating records
* For custom modules, use the exact API name of the module

**FAQs**

Q: Can I connect to multiple Zoho CRM accounts?
A: Yes, you can create separate connections for each Zoho CRM account.

Q: How do I handle custom fields in Zoho CRM?
A: Custom fields can be accessed using their API names, which typically follow the format `Custom_Field_Label`.

Q: Are there rate limits for the Zoho CRM API?
A: Yes, Zoho CRM has rate limits based on your subscription plan. Check the [Zoho CRM API documentation](https://www.zoho.com/crm/developer/docs/api/rate-limits.html) for specific details.

## Categories

- crm
- sales

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name           | Description                                                                                                                        | Link                                  |
|----------------|------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------|
| Create Record  | Creates a new record in a specified Zoho CRM module with the provided data                                                         | [docs](actions/create_record.md)      |
| Update Record  | Updates an existing record in a specified Zoho CRM module with the provided data                                                   | [docs](actions/update_record.md)      |
| Get Record     | Retrieves a specific record from a Zoho CRM module based on the record ID                                                          | [docs](actions/get_record.md)         |
| List Records   | Retrieves a list of records from a specified Zoho CRM module with optional filtering and pagination                                | [docs](actions/list_records.md)       |
| Search Records | Searches for records in a specified Zoho CRM module based on search criteria                                                       | [docs](actions/search_records.md)     |
| Delete Record  | Permanently removes a record from a specified Zoho CRM module                                                                     | [docs](actions/delete_record.md)      |

## Triggers

| Name            | Description                                                                                                                                                                            | Link                                   |
|-----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------|
| Record Created  | Triggers when a new record is created in a specified Zoho CRM module                                                                                                                  | [docs](triggers/record_created.md)     |
| Record Updated  | Triggers when an existing record is updated in a specified Zoho CRM module                                                                                                            | [docs](triggers/record_updated.md)     |
| Record Deleted  | Triggers when a record is deleted from a specified Zoho CRM module                                                                                                                   | [docs](triggers/record_deleted.md)     |