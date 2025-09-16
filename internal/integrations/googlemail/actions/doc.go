package actions

import _ "embed"

//go:embed get_mail.md
var getMailDocs string

//go:embed get_thread.md
var getThreadDocs string

//go:embed list_mails.md
var listMailsDocs string

//go:embed send_email.md
var sendEmailDocs string

//go:embed send_email_template.md
var sendEmailTemplateDocs string

//go:embed extract_email_text.md
var extractMailTextDocs string

//go:embed extract_document_from_email.md
var extractDocumentsFromEmailDocs string
