package zohoinventory

import (
	"github.com/wakflo/go-sdk/autoform"
)

var (
	// #nosec
	tokenURL   = "https://accounts.zoho.com/oauth/v2/token"
	authURL    = "https://accounts.zoho.com/oauth/v2/auth"
	sharedAuth = autoform.NewOAuthField(authURL, &tokenURL, []string{
		"ZohoInventory.FullAccess.all",
	}).
		Build()
)
