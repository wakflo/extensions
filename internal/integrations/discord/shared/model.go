package shared

type Guild struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon,omitempty"`
	Owner       bool     `json:"owner,omitempty"`
	Permissions string   `json:"permissions,omitempty"`
	Features    []string `json:"features,omitempty"`
}

type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	Position int    `json:"position,omitempty"`
	Topic    string `json:"topic,omitempty"`
	NSFW     bool   `json:"nsfw,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
}

type Role struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Color        int    `json:"color"`
	Hoist        bool   `json:"hoist"`
	Icon         string `json:"icon,omitempty"`
	UnicodeEmoji string `json:"unicode_emoji,omitempty"`
	Position     int    `json:"position"`
	Permissions  string `json:"permissions"`
	Managed      bool   `json:"managed"`
	Mentionable  bool   `json:"mentionable"`
}
