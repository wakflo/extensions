package actions

import _ "embed"

//go:embed list_contacts.md
var listContactsDocs string

//go:embed get_contact.md
var getContactDocs string

//go:embed create_contact.md
var createContactDocs string

//go:embed update_contact.md
var updateContactDocs string
