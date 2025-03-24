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

// GumroadListProductsResponse represents the top-level response from the Gumroad List Products API
// type ListProductsResponse struct {
// 	Success  bool            `json:"success"`
// 	Products []ProductDetail `json:"products"`
// 	Next     string          `json:"next,omitempty"`     // Pagination token for the next page of results
// 	Previous string          `json:"previous,omitempty"` // Pagination token for the previous page of results
// }

type ListProductsResponse struct {
	Success  bool            `json:"success"`
	Products []ProductDetail `json:"products"`
	// Next     string          `json:"next,omitempty"`     // Pagination token for the next page of results
	// Previous string          `json:"previous,omitempty"` // Pagination token for the previous page of results
}

// ProductDetail represents a single product in the Gumroad API response
type ProductDetail struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Description           string                `json:"description"`
	URL                   string                `json:"url"`
	PreviewURL            string                `json:"preview_url"`
	ThumbnailURL          string                `json:"thumbnail_url"`
	OriginalPrice         int                   `json:"original_price"` // In cents, 0 for pay-what-you-want products
	DisplayPrice          string                `json:"display_price"`  // Formatted price string e.g. "$10"
	Currency              string                `json:"currency"`
	ShortURL              string                `json:"short_url"`
	Created               string                `json:"created_at"` // ISO 8601 timestamp
	Updated               string                `json:"updated_at"` // ISO 8601 timestamp
	Published             bool                  `json:"published"`
	CustomFields          []CustomField         `json:"custom_fields,omitempty"`
	CustomReceipt         string                `json:"custom_receipt,omitempty"`
	CustomPermalink       string                `json:"custom_permalink,omitempty"`
	MaxPurchaseCount      int                   `json:"max_purchase_count,omitempty"`
	RecurrenceSettings    *RecurrenceSettings   `json:"recurrence_settings,omitempty"`
	CustomDeliveryURL     string                `json:"custom_delivery_url,omitempty"`
	RequireShipping       bool                  `json:"require_shipping"`
	Variants              []Variant             `json:"variants,omitempty"`
	Sales                 int                   `json:"sales_count"`
	Sales7Days            int                   `json:"sales_count_7_days"`
	Sales30Days           int                   `json:"sales_count_30_days"`
	RevenueUSD            int                   `json:"revenue_usd_cents"` // In cents
	RevenueUSD7Days       int                   `json:"revenue_usd_cents_7_days"`
	RevenueUSD30Days      int                   `json:"revenue_usd_cents_30_days"`
	Tags                  []string              `json:"tags,omitempty"`
	Categories            []string              `json:"categories,omitempty"`
	IsFilesOnly           bool                  `json:"is_files_only"`
	FileInfo              []FileInfo            `json:"file_info,omitempty"`
	Tiered                bool                  `json:"tiered"`
	TieredVariants        []TieredVariant       `json:"tiered_variants,omitempty"`
	HasInstallments       bool                  `json:"has_installments"`
	InstallmentSettings   *InstallmentSettings  `json:"installment_settings,omitempty"`
	HasMembership         bool                  `json:"has_membership"`
	MembershipSettings    *MembershipSettings   `json:"membership_settings,omitempty"`
	HasFreeProductsAccess bool                  `json:"has_free_products_access"`
	FreeProductsSettings  *FreeProductsSettings `json:"free_products_settings,omitempty"`
}

// CustomField represents a custom form field for a product
type CustomField struct {
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Type     string   `json:"type"`              // e.g., "text", "dropdown"
	Options  []string `json:"options,omitempty"` // For dropdown type
}

// RecurrenceSettings represents subscription settings for a product
type RecurrenceSettings struct {
	IntervalCount      int    `json:"interval_count"`
	IntervalUnit       string `json:"interval_unit"` // "day", "week", "month", or "year"
	AutoRecur          bool   `json:"auto_recur"`
	RecurringPrice     int    `json:"recurring_price_cents"`
	UseTrial           bool   `json:"use_trial"`
	TrialIntervalCount int    `json:"trial_interval_count,omitempty"`
	TrialIntervalUnit  string `json:"trial_interval_unit,omitempty"` // "day", "week", "month", or "year"
}

// Variant represents a product variant
type Variant struct {
	Title         string          `json:"title"`
	Options       []VariantOption `json:"options"`
	DefaultOption string          `json:"default_option,omitempty"`
}

// VariantOption represents an option within a product variant
type VariantOption struct {
	Name          string `json:"name"`
	PriceDiff     int    `json:"price_difference_cents"`
	IsInventoried bool   `json:"is_inventoried"`
	Quantity      int    `json:"quantity,omitempty"`
}

// FileInfo represents information about a file attached to a product
type FileInfo struct {
	FileName     string `json:"file_name"`
	FileSize     int    `json:"file_size"` // In bytes
	ContentType  string `json:"content_type"`
	LastModified string `json:"last_modified"` // ISO 8601 timestamp
}

// TieredVariant represents a tier for tiered pricing
type TieredVariant struct {
	Name        string `json:"name"`
	Price       int    `json:"price_cents"`
	Description string `json:"description,omitempty"`
}

// InstallmentSettings represents settings for installment payments
type InstallmentSettings struct {
	PaymentCount        int    `json:"payment_count"`
	PaymentInterval     int    `json:"payment_interval"`
	PaymentIntervalUnit string `json:"payment_interval_unit"` // "day", "week", "month", or "year"
}

// MembershipSettings represents settings for membership products
type MembershipSettings struct {
	Duration     int    `json:"duration"`
	DurationUnit string `json:"duration_unit"` // "day", "week", "month", or "year"
}

// FreeProductsSettings represents settings for free product access
type FreeProductsSettings struct {
	ProductIDs []string `json:"product_ids"`
}

// type DynamicOptionsResponse struct {
// 	Options []map[string]any `json:"options"`
// }
