package shopify

import (
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

/*var viewStyleOptions = []*sdkcore.AutoFormSchema{
	{Const: "shopify", Title: "Shopify"},
	{Const: "custom", Title: "Shopify (Custom/Private)"},
}*/

var (
	sharedAuth = autoform.NewCustomAuthField().
		SetFields(map[string]*sdkcore.AutoFormSchema{
			"domain": autoform.NewShortTextField().SetDisplayName("Domain Name").
				SetDescription("The domain name of the shopify app.").
				SetRequired(true).
				Build(),
			"token": autoform.NewShortTextField().SetDisplayName("Authentication Token").
				SetDescription("The token used to authenticate the shopify app.").
				Build(),
			/*// will be enabled when dropdown cab show in dialog
			// "appMode": autoform.NewSelectField().
			//	SetDisplayName("Application Mode").
			//	SetOptions(viewStyleOptions).
			//	SetRequired(true).
			//	SetDescription("The application mode of the shopify app.").
			//	Build(),*/
		}).
		Build()
)
