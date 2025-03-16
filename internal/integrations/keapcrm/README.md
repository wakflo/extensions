# Keap Integration

## Description

Integrate your Keap (formerly Infusionsoft) CRM data with our workflow automation software to streamline sales and marketing processes. Automatically sync Keap contacts, leads, and opportunities with your workflows, eliminating manual data entry and reducing errors. Use this integration to:

* Automatically create or update contacts in Keap when leads come in from other channels
* Trigger workflows based on contact activities or status changes in Keap
* Send personalized follow-up emails through Keap when specific events occur
* Update contact information across multiple systems when changes occur in Keap
* Generate tasks in your project management system based on Keap contact interactions
* Create comprehensive customer profiles by combining Keap data with information from other sources

## Prerequisites

* A Keap account
* Our workflow automation software account
* Keap API credentials (OAuth2 client ID and client secret)

## Setup Instructions

1. Log in to your Keap account and navigate to the "Settings" > "API" section.
2. Create a new OAuth2 application if you haven't already.
3. Note your Client ID and Client Secret.
4. In our workflow automation software, go to the "Integrations" section and select "Keap".
5. Enter your Client ID and Client Secret, then click "Connect".
6. Follow the authorization flow to grant our workflow access to your Keap account.

## Available Actions

| Name | Description |
|------|-------------|
| Get Contact | Retrieve detailed information about a specific contact from Keap by providing their contact ID |
| List Contacts | Retrieve a list of contacts from Keap based on optional filter criteria |
| Create Contact | Create a new contact in Keap with specified details |
| Update Contact | Update an existing contact in Keap with new information |

## Available Triggers

| Name | Description |
|------|-------------|
| Contact Created | Trigger a workflow when a new contact is created in Keap |
| Contact Updated | Trigger a workflow when a contact's information is updated in Keap |

## Troubleshooting Tips

* If you encounter connection issues, verify your Client ID and Client Secret are entered correctly.
* Ensure your Keap account has the necessary permissions to perform the requested actions.
* Check the API request limits for your Keap plan to avoid rate limiting issues.

## Categories

- crm
- marketing

## Authors

- Wakflo <integrations@wakflo.com>