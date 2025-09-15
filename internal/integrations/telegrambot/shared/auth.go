package shared

import "github.com/juicycleff/smartform/v1"

const markdown = `
**Authentication**:

1. Begin a conversation with the [Botfather](https://telegram.me/BotFather).
2. Type in "/newbot"
3. Choose a name for your bot
4. Choose a username for your bot.
5. Copy the token value from the Botfather and use it activepieces connection.
6. Congratulations! You can now use your new Telegram connection in your flows.
`

var (
	form = smartform.NewAuthForm("telegram-auth", "API Access Token", smartform.AuthStrategyCustom)

	_ = form.
		TextField("token", "Telegram Bot Token").
		Required(true).
		HelpText(markdown).
		Build()

	SharedAuth = form.Build()
)
