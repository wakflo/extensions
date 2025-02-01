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

import "time"

type WorkspaceUser struct {
	ActiveProjectCount          int           `json:"active_project_count"`
	Admin                       bool          `json:"admin"`
	APIToken                    string        `json:"api_token"`
	At                          time.Time     `json:"at"`
	BusinessWS                  bool          `json:"business_ws"`
	CSVUpload                   *string       `json:"csv_upload,omitempty"`
	DefaultCurrency             string        `json:"default_currency"`
	DefaultHourlyRate           *float64      `json:"default_hourly_rate,omitempty"`
	HideStartEndTimes           bool          `json:"hide_start_end_times"`
	IcalEnabled                 bool          `json:"ical_enabled"`
	IcalURL                     string        `json:"ical_url"`
	ID                          int           `json:"id"`
	LastModified                time.Time     `json:"last_modified"`
	LogoURL                     string        `json:"logo_url"`
	Name                        string        `json:"name"`
	OnlyAdminsMayCreateProjects bool          `json:"only_admins_may_create_projects"`
	OnlyAdminsMayCreateTags     bool          `json:"only_admins_may_create_tags"`
	OnlyAdminsSeeBillableRates  bool          `json:"only_admins_see_billable_rates"`
	OnlyAdminsSeeTeamDashboard  bool          `json:"only_admins_see_team_dashboard"`
	OrganizationID              int           `json:"organization_id"`
	Permissions                 *Permissions  `json:"permissions,omitempty"`
	Premium                     bool          `json:"premium"`
	Profile                     int           `json:"profile"`
	ProjectsBillableByDefault   bool          `json:"projects_billable_by_default"`
	ProjectsEnforceBillable     bool          `json:"projects_enforce_billable"`
	ProjectsPrivateByDefault    bool          `json:"projects_private_by_default"`
	RateLastUpdated             *time.Time    `json:"rate_last_updated,omitempty"`
	ReportsCollapse             bool          `json:"reports_collapse"`
	Role                        string        `json:"role"`
	Rounding                    int           `json:"rounding"`
	RoundingMinutes             int           `json:"rounding_minutes"`
	ServerDeletedAt             *time.Time    `json:"server_deleted_at,omitempty"`
	Subscription                *Subscription `json:"subscription,omitempty"`
	SuspendedAt                 *time.Time    `json:"suspended_at,omitempty"`
	WorkingHoursInMinutes       *int          `json:"working_hours_in_minutes,omitempty"`
}

type Subscription struct{}

type Permissions struct{}
