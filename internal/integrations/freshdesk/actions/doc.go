package actions

import _ "embed"

//go:embed create_ticket.md
var createTicketDocs string

//go:embed get_ticket.md
var getTicketDocs string

//go:embed update_ticket.md
var updateTicketDocs string

//go:embed search_ticket.md
var searchTicketDocs string

//go:embed list_tickets.md
var listTicketsDocs string
