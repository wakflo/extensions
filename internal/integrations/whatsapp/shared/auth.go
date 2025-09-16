package shared

import "github.com/juicycleff/smartform/v1"

const markdown = `
To Obtain a Phone Number ID and a Permanent System User Access Token, follow these steps:

1. Go to https://developers.facebook.com/
2. Make a new app, Select Other for usecase.
3. Choose Business as the type of app.
4. Add new Product -> WhatsApp.
5. Navigate to WhatsApp Settings > API Setup.
6. Select a phone number and copy its Phone number ID.
7. Login to your [Meta Business Manager](https://business.facebook.com/).
8. Click on Settings.
9. Create a new System User with access over the app and copy the access token.
`

var (
	form = smartform.NewAuthForm("whatsapp-auth", "API Access Token", smartform.AuthStrategyCustom)

	_ = form.
		TextField("token", "System User Access Token").
		Required(true).
		HelpText(markdown).
		Build()

	_ = form.TextField("phone-id", "Phone Number ID").
		Required(true).
		HelpText("The ID of your WhatsApp Business phone number. Find this in the Meta for Developers dashboard.")

	SharedAuth = form.Build()
)
