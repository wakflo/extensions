package shared

type Base struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}

type Response struct {
	Bases  []Base `json:"bases"`
	Offset string `json:"offset"`
}

type Table struct {
	Description    string  `json:"description"`
	Fields         []Field `json:"fields"`
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	PrimaryFieldID string  `json:"primaryFieldId"`
	Views          []View  `json:"views"`
}

type TableResponse struct {
	Tables []Table `json:"tables"`
}

type Field struct {
	Description string        `json:"description,omitempty"` // Omits field if no description is provided
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Options     *FieldOptions `json:"options,omitempty"` // Omits if options are not available
}

type FieldOptions struct {
	Color string `json:"color,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

type View struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type RecordAirtable struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type ResponseAirtable struct {
	Records []RecordAirtable `json:"records"`
}
