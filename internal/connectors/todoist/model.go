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

package todoist

import "time"

type Project struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CommentCount   int     `json:"comment_count"`
	Order          int     `json:"order"`
	Color          string  `json:"color"`
	IsShared       bool    `json:"is_shared"`
	IsFavorite     bool    `json:"is_favorite"`
	IsInboxProject bool    `json:"is_inbox_project"`
	IsTeamInbox    bool    `json:"is_team_inbox"`
	ViewStyle      string  `json:"view_style"`
	URL            string  `json:"url"`
	ParentID       *string `json:"parent_id"`
}

type CreateProject struct {
	Name       string  `json:"name"`
	Color      *string `json:"color"`
	IsFavorite *bool   `json:"is_favorite"`
	ViewStyle  *string `json:"view_style"`
	ParentID   *string `json:"parent_id"`
}

type UpdateProject struct {
	Name       *string `json:"name"`
	Color      *string `json:"color"`
	IsFavorite *bool   `json:"is_favorite"`
	ViewStyle  *string `json:"view_style"`
}

type Collaborator struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Task struct {
	CreatorID    string   `json:"creator_id"`
	CreatedAt    string   `json:"created_at"`
	AssigneeID   string   `json:"assignee_id"`
	AssignerID   string   `json:"assigner_id"`
	CommentCount int      `json:"comment_count"`
	IsCompleted  bool     `json:"is_completed"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	Due          Due      `json:"due"`
	Duration     Duration `json:"duration"`
	ID           string   `json:"id"`
	Labels       []string `json:"labels"`
	Order        int      `json:"order"`
	Priority     int      `json:"priority"`
	ProjectID    string   `json:"project_id"`
	SectionID    string   `json:"section_id"`
	ParentID     string   `json:"parent_id"`
	URL          string   `json:"url"`
}

type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime"`
	String      string `json:"string"`
	Timezone    string `json:"timezone"`
}

type Duration struct {
	Amount int    `json:"amount"`
	Unit   string `json:"unit"`
}

type ProjectSection struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Order     int    `json:"order"`
	Name      string `json:"name"`
}

type UpdateTask struct {
	Content     string     `json:"content"`
	Description *string    `json:"description"`
	Labels      []string   `json:"labels"`
	Priority    *int       `json:"priority"`
	Order       *string    `json:"order"`
	DueDate     *time.Time `json:"dueDate"`
	Duration    *int       `json:"duration"`
}
