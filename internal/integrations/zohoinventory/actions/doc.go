package actions

import _ "embed"

//go:embed get_invoice.md
var getInvoiceDocs string

//go:embed get_payment.md
var getPaymentDocs string

//go:embed list_invoices.md
var listInvoicesDocs string

//go:embed list_items.md
var listItemsDocs string

//go:embed list_payments.md
var listPaymentsDocs string
