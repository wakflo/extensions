## Extract Structured Data with Function Calling

This action uses Gemini's function calling to extract structured data from text.

### Use Cases
- Extract contact information from emails
- Parse invoice details from documents
- Extract product specifications from descriptions
- Convert natural language to API parameters
- Extract entities and relationships from text

### Function Schema Format
Provide a JSON schema defining the structure you want to extract.

Example:
{
  "name": {"type": "string", "description": "Person's full name"},
  "email": {"type": "string", "description": "Email address"},
  "phone": {"type": "string", "description": "Phone number"},
  "age": {"type": "number", "description": "Person's age"},
  "is_customer": {"type": "boolean", "description": "Whether they are a customer"}
}