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

package shared

type Team struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Issue struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type Assignee struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Label struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Priority struct {
	Priority int    `json:"priority,omitempty"`
	Label    string `json:"label,omitempty"`
}

type WorkflowState struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
