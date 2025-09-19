package shared

import "github.com/juicycleff/smartform/v1"

const markdown = `
**Authentication**:

1. Search for [**B͟o͟t͟Fa͟t͟h͟e͟r͟**](https://telegram.me/BotFather)   in Telegram and begin a conversation with it.
2. Type in "/newbot"
3. Choose a name for your bot
4. Choose a username for your bot.
5. Copy the token value from the Botfather and use it as your Wakflo connection.
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
