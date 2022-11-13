package start

import (
	"github.com/AnimeKaizoku/KaizokuRobot/utils"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func StartHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage
	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
		return ext.EndGroups
	}

	txt := mdparser.GetNormal("Can't parse the given id, please check it again.")
	_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})

	return ext.EndGroups
}

func LoadStartHandler(d *ext.Dispatcher, t []rune) {
	startCommand := handlers.NewCommand(utils.GetStartCommand(), StartHandler)

	startCommand.Triggers = t

	d.AddHandler(startCommand)
}
