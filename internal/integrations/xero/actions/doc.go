package actions

import _ "embed"

//go:embed create_invoice.md
var createInvoiceDocs string

//go:embed email_invoice.md
var emailInvoiceDocs string

//go:embed get_invoice.md
var getInvoiceDocs string

//go:embed list_invoices.md
var listInvoicesDocs string
