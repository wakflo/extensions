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

type Ticket struct {
	CcEmails        []string  `json:"cc_emails"`
	FwdEmails       []string  `json:"fwd_emails"`
	ReplyCcEmails   []string  `json:"reply_cc_emails"`
	DescriptionText string    `json:"description_text"`
	Description     string    `json:"description"`
	Spam            bool      `json:"spam"`
	EmailConfigID   *int      `json:"email_config_id"`
	FrEscalated     bool      `json:"fr_escalated"`
	GroupID         *int      `json:"group_id"`
	Priority        int       `json:"priority"`
	RequesterID     int       `json:"requester_id"`
	ResponderID     *int      `json:"responder_id"`
	Source          int       `json:"source"`
	Status          int       `json:"status"`
	Subject         string    `json:"subject"`
	ID              int       `json:"id"`
	Type            *string   `json:"type"`
	ToEmails        *string   `json:"to_emails"`
	ProductID       *int      `json:"product_id"`
	Attachments     []string  `json:"attachments"`
	IsEscalated     bool      `json:"is_escalated"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DueBy           time.Time `json:"due_by"`
	FrDueBy         time.Time `json:"fr_due_by"`
}

type TicketUpdate struct {
	Description string `json:"description"`
	Subject     string `json:"subject"`
	Priority    int    `json:"priority"`
	Status      int    `json:"status"`
	Tags        string `json:"tags"`
	Type        string `json:"type"`
}
