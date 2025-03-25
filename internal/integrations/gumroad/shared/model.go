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

type Response struct {
	Success  bool      `json:"success"`
	Products []Product `json:"products"`
}

type Product struct {
	CustomPermalink      interface{}            `json:"custom_permalink"`
	CustomReceipt        interface{}            `json:"custom_receipt"`
	CustomSummary        string                 `json:"custom_summary"`
	CustomFields         []interface{}          `json:"custom_fields"`
	CustomizablePrice    interface{}            `json:"customizable_price"`
	Description          string                 `json:"description"`
	Deleted              bool                   `json:"deleted"`
	MaxPurchaseCount     interface{}            `json:"max_purchase_count"`
	Name                 string                 `json:"name"`
	PreviewURL           interface{}            `json:"preview_url"`
	RequireShipping      bool                   `json:"require_shipping"`
	SubscriptionDuration interface{}            `json:"subscription_duration"`
	Published            bool                   `json:"published"`
	URL                  string                 `json:"url"`
	ID                   string                 `json:"id"`
	Price                int                    `json:"price"`
	PPPrices             PPPrices               `json:"purchasing_power_parity_prices"`
	Currency             string                 `json:"currency"`
	ShortURL             string                 `json:"short_url"`
	ThumbnailURL         string                 `json:"thumbnail_url"`
	Tags                 []string               `json:"tags"`
	FormattedPrice       string                 `json:"formatted_price"`
	FileInfo             map[string]interface{} `json:"file_info"`
	SalesCount           string                 `json:"sales_count"`
	SalesUSDCents        string                 `json:"sales_usd_cents"`
	IsTieredMembership   bool                   `json:"is_tiered_membership"`
	Recurrences          []string               `json:"recurrences"`
	Variants             []Variant              `json:"variants"`
}

type PPPrices struct {
	US int `json:"US"`
	IN int `json:"IN"`
	EC int `json:"EC"`
}

type Variant struct {
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type Option struct {
	Name             string           `json:"name"`
	PriceDifference  int              `json:"price_difference"`
	PPPrices         PPPrices         `json:"purchasing_power_parity_prices"`
	IsPayWhatYouWant bool             `json:"is_pay_what_you_want"`
	RecurrencePrices RecurrencePrices `json:"recurrence_prices"`
}

type RecurrencePrices struct {
	Monthly MonthlyPrice `json:"monthly"`
}

type MonthlyPrice struct {
	PriceCents          int      `json:"price_cents"`
	SuggestedPriceCents *int     `json:"suggested_price_cents"`
	PPPrices            PPPrices `json:"purchasing_power_parity_prices"`
}
