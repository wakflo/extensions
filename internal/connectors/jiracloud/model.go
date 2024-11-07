package jiracloud

type ProjectResponse struct {
	IsLast     bool      `json:"isLast"`
	MaxResults int       `json:"maxResults"`
	Self       string    `json:"self"`
	StartAt    int       `json:"startAt"`
	Total      int       `json:"total"`
	Values     []Project `json:"values"`
}

type Project struct {
	AvatarUrls     AvatarUrls `json:"avatarUrls"`
	Expand         string     `json:"expand"`
	ID             string     `json:"id"`
	IsPrivate      bool       `json:"isPrivate"`
	Key            string     `json:"key"`
	Name           string     `json:"name"`
	ProjectTypeKey string     `json:"projectTypeKey"`
	Properties     Properties `json:"properties"`
	Simplified     bool       `json:"simplified"`
	Style          string     `json:"style"`
	UUID           string     `json:"uuid"`
}

type Properties struct {
	Self string `json:"self"`
}

type User struct {
	AccountID   string     `json:"accountId"`
	AccountType string     `json:"accountType"`
	Active      bool       `json:"active"`
	AvatarUrls  AvatarUrls `json:"avatarUrls"`
	DisplayName string     `json:"displayName"`
	Key         string     `json:"key"`
	Name        string     `json:"name"`
	Self        string     `json:"self"`
}

type AvatarUrls struct {
	Size16x16 string `json:"16x16"`
	Size24x24 string `json:"24x24"`
	Size32x32 string `json:"32x32"`
	Size48x48 string `json:"48x48"`
}

type IssueType struct {
	AvatarID       int    `json:"avatarId"`
	Description    string `json:"description"`
	EntityID       string `json:"entityId"`
	HierarchyLevel int    `json:"hierarchyLevel"`
	IconURL        string `json:"iconUrl"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	Scope          Scope  `json:"scope"`
	Self           string `json:"self"`
	Subtask        bool   `json:"subtask"`
}

type Scope struct {
	Project Project `json:"project"`
	Type    string  `json:"type"`
}

type PrioritySchemeResponse struct {
	IsLast     bool             `json:"isLast"`
	MaxResults int              `json:"maxResults"`
	StartAt    int              `json:"startAt"`
	Total      int              `json:"total"`
	Values     []PriorityScheme `json:"values"`
}

type PriorityScheme struct {
	Description string      `json:"description"`
	ID          string      `json:"id"`
	IsDefault   bool        `json:"isDefault"`
	Name        string      `json:"name"`
	Priorities  PrioritySet `json:"priorities"`
	Projects    ProjectSet  `json:"projects"`
}

type PrioritySet struct {
	IsLast     bool       `json:"isLast"`
	MaxResults int        `json:"maxResults"`
	StartAt    int        `json:"startAt"`
	Total      int        `json:"total"`
	Values     []Priority `json:"values"`
}

type Priority struct {
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
	ID          string `json:"id"`
	IsDefault   bool   `json:"isDefault"`
	Name        string `json:"name"`
	StatusColor string `json:"statusColor"`
}

type ProjectSet struct {
	IsLast     bool      `json:"isLast"`
	MaxResults int       `json:"maxResults"`
	StartAt    int       `json:"startAt"`
	Total      int       `json:"total"`
	Values     []Project `json:"values"`
}

type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

type SearchIssuesResponse struct {
	Issues []Issue `json:"issues"`
}

type ProjectCategory struct {
	Description string `json:"description"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}
