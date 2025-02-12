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

type Tenant struct {
	ID             string  `json:"id"`
	AuthEventID    string  `json:"authEventId"`
	TenantID       string  `json:"tenantId"`
	TenantType     string  `json:"tenantType"`
	TenantName     *string `json:"tenantName"` // Using a pointer to handle null values
	CreatedDateUtc string  `json:"createdDateUtc"`
	UpdatedDateUtc string  `json:"updatedDateUtc"`
}

type TenantsResponse []Tenant

type InvoicesResponse struct {
	Invoices []Invoice `json:"Invoices"`
}

type Invoice struct {
	Contact       Contact `json:"Contact"`
	InvoiceID     string  `json:"InvoiceID"`
	InvoiceNumber string  `json:"InvoiceNumber"`
}

type Contact struct {
	ContactID      string    `json:"ContactID"`
	ContactStatus  string    `json:"ContactStatus"`
	Name           string    `json:"Name"`
	Addresses      []Address `json:"Addresses"`
	Phones         []Phone   `json:"Phones"`
	UpdatedDateUTC string    `json:"UpdatedDateUTC"`
	IsSupplier     string    `json:"IsSupplier"`
	IsCustomer     string    `json:"IsCustomer"`
}

type Address struct {
	AddressType  string `json:"AddressType"`
	AddressLine1 string `json:"AddressLine1,omitempty"`
	AddressLine2 string `json:"AddressLine2,omitempty"`
	City         string `json:"City,omitempty"`
	PostalCode   string `json:"PostalCode,omitempty"`
}

type Phone struct {
	PhoneType string `json:"PhoneType"`
}

type LineItem struct {
	ItemCode    string     `json:"ItemCode"`
	Description string     `json:"Description"`
	Quantity    string     `json:"Quantity"`
	UnitAmount  string     `json:"UnitAmount"`
	TaxType     string     `json:"TaxType"`
	TaxAmount   string     `json:"TaxAmount"`
	LineAmount  string     `json:"LineAmount"`
	AccountCode string     `json:"AccountCode"`
	AccountID   string     `json:"AccountId"`
	Item        Item       `json:"Item"`
	Tracking    []Tracking `json:"Tracking"`
	LineItemID  string     `json:"LineItemID"`
}

type Item struct {
	ItemID string `json:"ItemID"`
	Name   string `json:"Name"`
	Code   string `json:"Code"`
}

type Tracking struct {
	TrackingCategoryID string `json:"TrackingCategoryID"`
	Name               string `json:"Name"`
	Option             string `json:"Option"`
}

type Payment struct {
	Date      string `json:"Date"`
	Amount    string `json:"Amount"`
	PaymentID string `json:"PaymentID"`
}
