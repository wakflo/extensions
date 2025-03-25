package actions

import _ "embed"

//go:embed list_products.md
var listProductsDocs string

// go:embed get_product.md
var getProductDocs string

// go:embed disable_product.md
var disableProductDocs string

// go:embed enable_product.md
var enableProductDocs string

// go:embed get_sales.md
var listSalesDocs string

// go:embed get_sale.md
var getSaleDocs string

// go:embed delete_product.md
var deleteProductDocs string

// go:embed mark_as_shipped.md
var markasShippedDocs string
