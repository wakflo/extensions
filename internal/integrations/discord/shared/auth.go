package shared

import "github.com/juicycleff/smartform/v1"

const markdown = `
To obtain a token, follow these steps:
1. Go to https://discord.com/developers/applications
2. Click on Application (or create one if you don't have one)
3. Click on Bot
4. Copy the token
`

var (
	form = smartform.NewAuthForm("discord-auth", "Discord Authentication", smartform.AuthStrategyCustom)

	_ = form.
		TextField("token", "Bot Token").
		Required(true).
		HelpText(markdown).
		Build()

	SharedAuth = form.Build()
)
