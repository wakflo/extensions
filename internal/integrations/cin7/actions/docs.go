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

package actions

import (
	_ "embed"
)

//go:embed get_a_payment.md
var getAPaymentDocs string

//go:embed get_customers.md
var getCustomersDocs string

//go:embed get_products.md
var getProductsDocs string

//go:embed get_purchase_invoice.md
var getPurchaseInvoiceDocs string

//go:embed get_sales_list.md
var getSalesListDocs string

//go:embed get_sales_order.md
var getSalesOrderDocs string
