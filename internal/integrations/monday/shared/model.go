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

type Workspace struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Board struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Boards struct {
	Groups []Group `json:"groups"`
}

type Response struct {
	Data struct {
		Boards []Boards `json:"boards"`
	} `json:"data"`
}

type Group struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ItemsPage struct {
	Items []Item `json:"items"`
}

type ItemBoard struct {
	ItemsPage ItemsPage `json:"items_page"`
}

type ResponseData struct {
	Boards []ItemBoard `json:"boards"`
}

type ItemResponse struct {
	Data ResponseData `json:"data"`
}
