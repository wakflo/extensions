# Smart Data Extractor

Extract structured data from any unstructured text using AI.

## Key Features
- **Custom Schema**: Define exactly what data structure you need
- **Multiple Extraction Types**: Invoices, emails, forms, contracts, etc.
- **Validation**: Ensure extracted data matches your requirements
- **Batch Support**: Extract multiple items from a single document
- **Confidence Scores**: Get confidence levels for extracted data

## Common Use Cases

### Invoice Processing
Extract vendor, date, line items, totals, tax information

### Email Parsing
Extract sender info, intent, key dates, action items

### Resume/CV Processing
Extract contact info, experience, skills, education

### Form Digitization
Convert paper forms or PDFs into structured data

### Contract Analysis
Extract parties, dates, terms, obligations, amounts

### Customer Feedback
Extract sentiment, topics, issues, suggestions

## Schema Definition
Define your schema as a JSON object describing the structure you want:
{
  "vendor": "string",
  "invoice_date": "date",
  "items": [{"description": "string", "amount": "number"}],
  "total": "number"
}

## Best Practices
- Provide clear, specific schemas
- Include examples for complex extractions
- Use low temperature (0.1-0.3) for consistency
- Validate critical data after extraction