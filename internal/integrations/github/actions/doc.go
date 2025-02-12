package actions

import _ "embed"

//go:embed create_issue.md
var createIssueDocs string

//go:embed create_issue_comment.md
var createIssueCommentDocs string

//go:embed get_issue.md
var getIssueDocs string

//go:embed lock_issue.md
var lockIssueDocs string

//go:embed unlock_issue.md
var unlockIssueDocs string
