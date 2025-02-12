package actions

import _ "embed"

//go:embed create_customer.md
var createCustomerDocs string

//go:embed create_product.md
var createProductDocs string

//go:embed find_coupon.md
var findCouponDocs string

//go:embed find_customer.md
var findCustomerDocs string

//go:embed find_product.md
var findProductDocs string

//go:embed get_customer_by_id.md
var getCustomerByIdDocs string

//go:embed list_orders.md
var listOrdersDocs string

//go:embed list_products.md
var listProductsDocs string

//go:embed update_product.md
var updateProductDocs string
