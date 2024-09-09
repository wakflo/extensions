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

package clickup

type Team struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Color   string   `json:"color"`
	Avatar  string   `json:"avatar"`
	Members []Member `json:"members"`
}

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type Space struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Private           bool        `json:"private"`
	Color             *string     `json:"color,omitempty"`
	Avatar            *string     `json:"avatar,omitempty"`
	AdminCanManage    bool        `json:"admin_can_manage"`
	Archived          bool        `json:"archived"`
	Members           []Member    `json:"members"`
	Statuses          []Status    `json:"statuses"`
	MultipleAssignees bool        `json:"multiple_assignees"`
	Features          interface{} `json:"features"`
}

type SpacesResponse struct {
	Spaces []Space `json:"spaces"`
}

type Task struct {
	ID                  string      `json:"id"`
	CustomItemID        interface{} `json:"custom_item_id"`
	Name                string      `json:"name"`
	Status              Status      `json:"status"`
	MarkdownDescription string      `json:"markdown_description"`
	OrderIndex          string      `json:"orderindex"`
	DateCreated         string      `json:"date_created"`
	DateUpdated         string      `json:"date_updated"`
	DateClosed          interface{} `json:"date_closed"`
	DateDone            interface{} `json:"date_done"`
	Parent              interface{} `json:"parent"`
	Priority            interface{} `json:"priority"`
	DueDate             interface{} `json:"due_date"`
	StartDate           interface{} `json:"start_date"`
	Points              int         `json:"points"`
	TimeEstimate        interface{} `json:"time_estimate"`
	TimeSpent           interface{} `json:"time_spent"`
	List                List        `json:"list"`
	Folder              Folder      `json:"folder"`
	Space               Space       `json:"space"`
	URL                 string      `json:"url"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks,omitempty"`
}

type Folder struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	OrderIndex       int    `json:"orderindex"`
	OverrideStatuses bool   `json:"override_statuses"`
	Hidden           bool   `json:"hidden"`
	Space            Space  `json:"space"`
	TaskCount        string `json:"task_count"`
	Lists            []List `json:"lists"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}

type List struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	OrderIndex       int         `json:"orderindex"`
	Content          string      `json:"content"`
	Status           Status      `json:"status"`
	Priority         Priority    `json:"priority"`
	Assignee         interface{} `json:"assignee"`
	TaskCount        interface{} `json:"task_count"`
	DueDate          string      `json:"due_date"`
	StartDate        interface{} `json:"start_date"`
	Folder           Folder      `json:"folder"`
	Space            Space       `json:"space"`
	Archived         bool        `json:"archived"`
	OverrideStatuses bool        `json:"override_statuses"`
	PermissionLevel  string      `json:"permission_level"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}

type Status struct{}

type Priority struct{}

type MembersResponse struct {
	Members []Member `json:"members"`
}

type Member struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Color          *string `json:"color"`
	Initials       string  `json:"initials"`
	ProfilePicture string  `json:"profilePicture"`
}
