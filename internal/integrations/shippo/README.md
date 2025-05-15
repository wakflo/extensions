# Shippo Integration Documentation
Overview
Shippo is a multi-carrier shipping API that allows you to connect with multiple shipping carriers through a single, unified API. The Shippo integration with Wakflo enables you to:

Create shipments
Generate shipping labels
Track packages
Compare rates across carriers

Authentication
To use the Shippo integration, you'll need:

A Shippo account
Your Shippo API key

The API key should be provided in the authentication configuration as "api-key".
Available Actions
Create New Shipment
Creates a new shipment in Shippo with sender and recipient information and package details.
Required Fields:
Sender Information:

Sender's Name: Full name of the sender
Sender's Street 1: Street address of the sender
Sender's Email: Email address of the sender
Sender's City: City of the sender
Sender's Country: Country of the sender (select from dropdown)
Sender's State: State or province of the sender
Sender's Zip: Postal code of the sender
Sender's Phone: Phone number of the sender

Recipient Information:

Receiver's Name: Full name of the recipient
Receiver's Street 1: Street address of the recipient
Receiver's Email: Email address of the recipient
Receiver's City: City of the recipient
Receiver's Country: Country of the recipient (select from dropdown)
Receiver's State: State or province of the recipient
Receiver's Zip: Postal code of the recipient
Receiver's Phone: Phone number of the recipient

Package Information:

Parcel length: Length of the package
Parcel width: Width of the package
Parcel height: Height of the package
Distance Unit: Unit of measurement for dimensions (e.g., in, cm)
Parcel weight: Weight of the package
Mass Unit: Unit of measurement for weight (e.g., oz, g, lb, kg)

Response:
Returns a Shippo shipment object containing all details of the created shipment.
Create Shipment Label
Creates a shipping label based on a rate object ID.
Required Fields:

Rate: The Shippo rate object ID

Response:
Returns a transaction object containing the label URL and tracking information.
Error Handling
Common errors that may occur:

Missing API key: Ensure your Shippo API key is provided in the authentication settings
Invalid address information: Verify all address fields are correctly formatted
Invalid parcel dimensions: Ensure all dimensions and weight values are valid numbers

Best Practices

Always validate address information before creating shipments
Compare rates from multiple carriers before purchasing a label
Store the tracking number and label URL for future reference
Use appropriate unit measurements based on the origin and destination countries

Additional Resources

Shippo API Documentation
Supported Carriers
Address Validation Guide