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

type ListResponse struct {
	Lists []List `json:"lists"`
}

type List struct {
	ID                   string           `json:"id"`
	WebID                int              `json:"web_id"`
	Name                 string           `json:"name"`
	Contact              Contact          `json:"contact"`
	PermissionReminder   string           `json:"permission_reminder"`
	UseArchiveBar        bool             `json:"use_archive_bar"`
	CampaignDefaults     CampaignDefaults `json:"campaign_defaults"`
	NotifyOnSubscribe    bool             `json:"notify_on_subscribe"`
	NotifyOnUnsubscribe  bool             `json:"notify_on_unsubscribe"`
	DateCreated          time.Time        `json:"date_created"`
	ListRating           int              `json:"list_rating"`
	EmailTypeOption      bool             `json:"email_type_option"`
	SubscribeURLShort    string           `json:"subscribe_url_short"`
	SubscribeURLLong     string           `json:"subscribe_url_long"`
	BeamerAddress        string           `json:"beamer_address"`
	Visibility           string           `json:"visibility"`
	DoubleOptin          bool             `json:"double_optin"`
	HasWelcome           bool             `json:"has_welcome"`
	MarketingPermissions bool             `json:"marketing_permissions"`
	Modules              []string         `json:"modules"`
	Stats                Stats            `json:"stats"`
	Links                []Link           `json:"_links"`
}

type Contact struct {
	Company  string `json:"company"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	State    string `json:"state"`
	Zip      string `json:"zip"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
}

type CampaignDefaults struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Language  string `json:"language"`
}

type Stats struct {
	MemberCount               int       `json:"member_count"`
	TotalContacts             int       `json:"total_contacts"`
	UnsubscribeCount          int       `json:"unsubscribe_count"`
	CleanedCount              int       `json:"cleaned_count"`
	MemberCountSinceSend      int       `json:"member_count_since_send"`
	UnsubscribeCountSinceSend int       `json:"unsubscribe_count_since_send"`
	CleanedCountSinceSend     int       `json:"cleaned_count_since_send"`
	CampaignCount             int       `json:"campaign_count"`
	CampaignLastSent          time.Time `json:"campaign_last_sent"`
	MergeFieldCount           int       `json:"merge_field_count"`
	AvgSubRate                float64   `json:"avg_sub_rate"`
	AvgUnsubRate              float64   `json:"avg_unsub_rate"`
	TargetSubRate             float64   `json:"target_sub_rate"`
	OpenRate                  float64   `json:"open_rate"`
	ClickRate                 float64   `json:"click_rate"`
	LastSubDate               time.Time `json:"last_sub_date"`
	LastUnsubDate             time.Time `json:"last_unsub_date"`
}

type Link struct {
	Rel          string `json:"rel"`
	Href         string `json:"href"`
	Method       string `json:"method"`
	TargetSchema string `json:"targetSchema"`
	Schema       string `json:"schema"`
}

type Constraints struct {
	MayCreate             bool `json:"may_create"`
	MaxInstances          int  `json:"max_instances"`
	CurrentTotalInstances int  `json:"current_total_instances"`
}
