package actions

import _ "embed"

//go:embed list_subscribers.md
var listSubscribersDocs string

//go:embed list_tags.md
var listTagsDocs string

//go:embed get_subscriber.md
var getSubscriberDocs string

// //go:embed create_subscriber.md
// var createSubscriberDocs string

//go:embed create_tag.md
var createTagDocs string

//go:embed tag_subscriber.md
var tagSubscriberDocs string
