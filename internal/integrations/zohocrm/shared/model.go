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

// ModuleListResponse represents the response from Zoho CRM modules API
type ModuleListResponse struct {
	Modules []Module `json:"modules"`
}

// Module represents a Zoho CRM module
type Module struct {
	APIName             string `json:"api_name"`
	ModuleName          string `json:"module_name"`
	PluralLabel         string `json:"plural_label"`
	SingularLabel       string `json:"singular_label"`
	ID                  string `json:"id"`
	Visible             bool   `json:"visible"`
	APIsupported        bool   `json:"api_supported"`
	Creatable           bool   `json:"creatable"`
	Deletable           bool   `json:"deletable"`
	Editable            bool   `json:"editable"`
	Convertable         bool   `json:"convertable"`
	Viewable            bool   `json:"viewable"`
	InventoryTemplateSupported bool `json:"inventory_template_supported"`
}
