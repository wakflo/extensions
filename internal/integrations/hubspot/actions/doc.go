package actions

import _ "embed"

//go:embed list_contacts.md
var listContactsDocs string

//go:embed list_tickets.md
var listTicketsDocs string

//go:embed create_contact.md
var createContactDocs string

//go:embed create_ticket.md
var createTicketDocs string

//go:embed retrieve_contact.md
var retrieveContactDocs string

//go:embed search_owner.md
var searchOwnerDocs string

//go:embed get_deal.md
var getDealDocs string
