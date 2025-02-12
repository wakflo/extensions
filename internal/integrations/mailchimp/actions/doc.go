package actions

import _ "embed"

//go:embed add_member_to_list.md
var addMemberToListDocs string

//go:embed add_note_to_subscriber.md
var addNoteToSubscriberDocs string

//go:embed add_subscriber_to_tag.md
var addSubscriberToTagDocs string

//go:embed get_all_list.md
var getAllTagDocs string

//go:embed remove_subscriber_from_tag.md
var removeSubscriberFromTagDocs string

//go:embed update_subscriber_status.md
var updateSubscriberStatusDocs string
