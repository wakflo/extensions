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

type MeetingRegistrant struct {
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name,omitempty"`
	Email                 string `json:"email"`
	Address               string `json:"address,omitempty"`
	City                  string `json:"city,omitempty"`
	State                 string `json:"state,omitempty"`
	Zip                   string `json:"zip,omitempty"`
	Country               string `json:"country,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	Comments              string `json:"comments,omitempty"`
	Industry              string `json:"industry,omitempty"`
	JobTitle              string `json:"job_title,omitempty"`
	NumberOfEmployees     string `json:"no_of_employees,omitempty"`
	Organization          string `json:"org,omitempty"`
	PurchasingTimeFrame   string `json:"purchasing_time_frame,omitempty"`
	RoleInPurchaseProcess string `json:"role_in_purchase_process,omitempty"`
	Language              string `json:"language,omitempty"`
}
