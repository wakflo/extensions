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

package zohoinventory

type Organization struct {
	OrganizationID       string `json:"organization_id"`
	Name                 string `json:"name"`
	ContactName          string `json:"contact_name"`
	Email                string `json:"email"`
	IsDefaultOrg         bool   `json:"is_default_org"`
	LanguageCode         string `json:"language_code"`
	FiscalYearStartMonth int    `json:"fiscal_year_start_month"`
	AccountCreatedDate   string `json:"account_created_date"`
	TimeZone             string `json:"time_zone"`
	IsOrgActive          bool   `json:"is_org_active"`
	CurrencyID           string `json:"currency_id"`
	CurrencyCode         string `json:"currency_code"`
	CurrencySymbol       string `json:"currency_symbol"`
	CurrencyFormat       string `json:"currency_format"`
	PricePrecision       int    `json:"price_precision"`
}

type Organizations struct {
	Organizations []Organization `json:"organizations"`
}
