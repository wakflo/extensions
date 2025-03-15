package shared

import (
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
)

var SharedAuth = autoform.NewAuth().NewCustomAuth().
	SetDescription("API Key").
	SetLabel("ConvertKit Authentication").
	SetFields(map[string]*sdkcore.AutoFormSchema{
		"api-key": autoform.NewShortTextField().
			SetDisplayName("API Key").
			SetDescription("ConvertKit API Key").
			SetRequired(true).
			Build(),
		"api-secret": autoform.NewShortTextField().
			SetDisplayName("API Secret").
			SetDescription("ConvertKit API Secret").
			SetRequired(true).
			Build(),
	}).
	Build()
