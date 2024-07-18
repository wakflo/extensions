package woocommerce

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	//"log"
	//"net/http"

	//"encoding/base64"
	//"fmt"
	//"strings"
	//
	//"github.com/go-resty/resty/v2"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var sharedAuth = autoform.NewCustomAuthField().
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"shop-url": autoform.NewShortTextField().SetDisplayName("ShopUrl").
			SetDescription("The base URL of your app (e.g https://mystore.com) and it should start with HTTPS only").
			SetRequired(true).
			Build(),
		"consumer-key": autoform.NewShortTextField().SetDisplayName("Consumer Key").
			SetDescription("The consumer key generated from your app.").
			Build(),
		"consumer-secret": autoform.NewShortTextField().SetDisplayName("Consumer Secret").
			SetDescription("The consumer secret generated from your app.").
			Build(),
	}).
	Build()
