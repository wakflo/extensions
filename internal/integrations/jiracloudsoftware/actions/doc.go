package actions

import _ "embed"

//go:embed create_issue.md
var createIssueDocs string

//go:embed get_issue.md
var getIssueDocs string

//go:embed list_issue.md
var listIssuesDocs string

//go:embed update_issue.md
var updateIssueDocs string

//go:embed transition_issue.md
var transitionIssueDocs string

//go:embed add_comment.md
var addCommentDocs string
