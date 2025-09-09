# DocConv Integration

## Description

Extract text and convert documents from various formats including images (OCR), PDFs, Word documents, PowerPoint presentations, and more. DocConv provides powerful document processing capabilities with support for:

* OCR text extraction from images (PNG, JPEG, TIFF, etc.)
* PDF text extraction and conversion
* Microsoft Office document processing (Word, Excel, PowerPoint)
* HTML conversion from various document formats
* OpenDocument format support (ODT, ODS, ODP)
* RTF and plain text document processing
* Batch processing capabilities
* Metadata extraction from documents
* Language detection and text analysis

**DocConv Integration Documentation**

**Overview**
The DocConv integration allows you to seamlessly extract text and convert documents from various formats using the powerful docconv library. This integration supports OCR (Optical Character Recognition) for images, PDF processing, Microsoft Office documents, and many other formats.

**Prerequisites**

* No external API keys required - all processing is done locally
* Supported file formats: PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX, ODT, ODS, ODP, RTF, TXT, PNG, JPEG, TIFF, etc.
* For OCR functionality, tesseract-ocr may need to be installed on the system

**Setup Instructions**

1. The DocConv integration works out of the box without any configuration
2. For optimal OCR results, ensure tesseract-ocr is installed on your system
3. Large files may require increased memory allocation

**Available Actions**

* **Extract Text from Image**: Use OCR to extract text from image files (PNG, JPEG, TIFF, etc.)
* **Extract Text from PDF**: Extract text content from PDF documents
* **Extract Text from Document**: Extract text from Microsoft Office and OpenDocument files
* **Convert Document to HTML**: Convert various document formats to HTML

**Example Use Cases**

1. **Document Digitization**: Convert scanned documents and images to searchable text
2. **Content Migration**: Extract text from legacy document formats for modern systems
3. **Data Processing**: Extract structured data from invoices, forms, and reports
4. **Compliance**: Convert documents to standardized formats for archival
5. **Search Indexing**: Extract text content for full-text search capabilities

**Troubleshooting Tips**

* Ensure input files are not corrupted or password-protected
* For OCR, higher resolution images generally produce better results
* Some complex document layouts may require preprocessing
* Large files may take longer to process

**FAQs**

Q: What image formats are supported for OCR?
A: PNG, JPEG, TIFF, BMP, and other common image formats are supported.

Q: Can I extract text from password-protected PDFs?
A: No, password-protected documents need to be unlocked before processing.

Q: What languages are supported for OCR?
A: OCR language support depends on your tesseract installation and language packs.

## Categories

- productivity
- document  
- ocr
- converter

## Authors

- Wakflo <integrations@wakflo.com>

## Actions

| Name                      | Description                                                                                           | Link                                     |
|---------------------------|-------------------------------------------------------------------------------------------------------|------------------------------------------|
| Extract Text from Image   | Use OCR to extract text from image files including PNG, JPEG, TIFF and other common image formats  | [docs](actions/extract_text_from_image.md)   |
| Extract Text from PDF     | Extract text content from PDF documents, including both text-based and scanned PDFs                 | [docs](actions/extract_text_from_pdf.md)     |
| Extract Text from Document| Extract text from various document formats including Word, Excel, PowerPoint and OpenDocument files | [docs](actions/extract_text_from_document.md)|
| Convert Document to HTML  | Convert various document formats to HTML while preserving structure and formatting                   | [docs](actions/convert_document_to_html.md)  |

## Triggers

| Name              | Description                                                                                                    | Link                                   |
|-------------------|----------------------------------------------------------------------------------------------------------------|----------------------------------------|
| Document Uploaded | Triggers when a document is uploaded or available for processing, allowing automatic text extraction workflows | [docs](triggers/document_uploaded.md) |