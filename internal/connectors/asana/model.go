// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asana

type WorkspaceNew struct {
	GID            string   `json:"gid"`
	ResourceType   string   `json:"resource_type"`
	Name           string   `json:"name"`
	EmailDomains   []string `json:"email_domains"`
	IsOrganization bool     `json:"is_organization"`
}

type WorkspaceResponse struct {
	Data []Workspace `json:"data"`
}

type ProjectResponse struct {
	Data []Project `json:"data"`
}

type Project struct {
	GID                                string               `json:"gid"`
	ResourceType                       string               `json:"resource_type"`
	Name                               string               `json:"name"`
	Archived                           bool                 `json:"archived"`
	Color                              string               `json:"color"`
	CreatedAt                          string               `json:"created_at"`
	CurrentStatus                      *CurrentStatus       `json:"current_status"`
	CurrentStatusUpdate                *CurrentStatusUpdate `json:"current_status_update"`
	CustomFieldSettings                []CustomFieldSetting `json:"custom_field_settings"`
	DefaultView                        string               `json:"default_view"`
	DueDate                            string               `json:"due_date"`
	DueOn                              string               `json:"due_on"`
	HTMLNotes                          string               `json:"html_notes"`
	Members                            []User               `json:"members"`
	ModifiedAt                         string               `json:"modified_at"`
	Notes                              string               `json:"notes"`
	PrivacySetting                     string               `json:"privacy_setting"`
	StartOn                            string               `json:"start_on"`
	DefaultAccessLevel                 string               `json:"default_access_level"`
	MinimumAccessLevelForCustomization string               `json:"minimum_access_level_for_customization"`
	MinimumAccessLevelForSharing       string               `json:"minimum_access_level_for_sharing"`
	CustomFields                       []CustomField        `json:"custom_fields"`
	Completed                          bool                 `json:"completed"`
	CompletedAt                        string               `json:"completed_at"`
	CompletedBy                        *User                `json:"completed_by"`
	Followers                          []User               `json:"followers"`
	Owner                              *User                `json:"owner"`
	Team                               *Team                `json:"team"`
	Icon                               string               `json:"icon"`
	PermalinkURL                       string               `json:"permalink_url"`
	ProjectBrief                       *ProjectBrief        `json:"project_brief"`
	CreatedFromTemplate                *Template            `json:"created_from_template"`
	Workspace                          *Workspace           `json:"workspace"`
}

type CurrentStatus struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	HTMLText     string `json:"html_text"`
	Color        string `json:"color"`
	Author       *User  `json:"author"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    *User  `json:"created_by"`
	ModifiedAt   string `json:"modified_at"`
}

type CurrentStatusUpdate struct {
	GID             string `json:"gid"`
	ResourceType    string `json:"resource_type"`
	Title           string `json:"title"`
	ResourceSubtype string `json:"resource_subtype"`
}

type CustomFieldSetting struct {
	GID          string       `json:"gid"`
	ResourceType string       `json:"resource_type"`
	Project      *Project     `json:"project"`
	IsImportant  bool         `json:"is_important"`
	Parent       *Project     `json:"parent"`
	CustomField  *CustomField `json:"custom_field"`
}

type CustomField struct {
	GID                     string       `json:"gid"`
	ResourceType            string       `json:"resource_type"`
	Name                    string       `json:"name"`
	ResourceSubtype         string       `json:"resource_subtype"`
	Type                    string       `json:"type"`
	EnumOptions             []EnumOption `json:"enum_options"`
	Enabled                 bool         `json:"enabled"`
	RepresentationType      string       `json:"representation_type"`
	IDPrefix                string       `json:"id_prefix"`
	IsFormulaField          bool         `json:"is_formula_field"`
	DateValue               *DateValue   `json:"date_value"`
	EnumValue               *EnumOption  `json:"enum_value"`
	MultiEnumValues         []EnumOption `json:"multi_enum_values"`
	NumberValue             float64      `json:"number_value"`
	TextValue               string       `json:"text_value"`
	DisplayValue            string       `json:"display_value"`
	Description             string       `json:"description,omitempty"`
	Precision               int          `json:"precision,omitempty"`
	Format                  string       `json:"format,omitempty"`
	CurrencyCode            string       `json:"currency_code,omitempty"`
	CustomLabel             string       `json:"custom_label,omitempty"`
	CustomLabelPosition     string       `json:"custom_label_position,omitempty"`
	IsGlobalToWorkspace     bool         `json:"is_global_to_workspace,omitempty"`
	HasNotificationsEnabled bool         `json:"has_notifications_enabled,omitempty"`
	AsanaCreatedField       string       `json:"asana_created_field,omitempty"`
	IsValueReadOnly         bool         `json:"is_value_read_only,omitempty"`
	CreatedBy               *User        `json:"created_by,omitempty"`
	PeopleValue             []User       `json:"people_value,omitempty"`
}

type EnumOption struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Color        string `json:"color"`
}

type DateValue struct {
	Date     string `json:"date"`
	DateTime string `json:"date_time"`
}

type User struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type Team struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type ProjectBrief struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}

type Template struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type Workspace struct {
	GID          string `json:"gid"`
	ResourceType string `json:"resource_type"`
	Name         string `json:"name"`
}

type Task struct {
	Name      string   `json:"name"`
	Workspace string   `json:"workspace"` // Add a default workspace ID
	Projects  []string `json:"projects"`
}
