// actions/doc.go
package actions

import _ "embed"

//go:embed list_products.md
var listProductsDocs string

//go:embed get_product.md
var getProductDocs string

//go:embed list_orders.md
var listOrdersDocs string

//go:embed get_order.md
var getOrderDocs string
