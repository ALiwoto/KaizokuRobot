package help

import (
	"github.com/AnimeKaizoku/KaizokuRobot/utils"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func HelpHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
		return ext.ContinueGroups
	}

	txt := mdparser.GetBold("Help section of KaizokuRobot\n\n")
	txt.Normal("Available commands are as follows: \n\n")
	txt.Mono(utils.GetStartCommand()).Normal(" : ").Italic("The all boring start command.\n\n")
	txt.Mono(utils.GetSendCommand()).Normal(" : ").Italic("The actual useful send command used for sending messages to channels.\n\n")
	txt.Mono(utils.GetAddCommand()).Normal(" : ").Italic("The add command for adding new chat id to json.\n\n")
	txt.Mono(utils.GetRemoveCommand()).Normal(" : ").Italic("The remove command for removing chat id from json.\n\n")
	txt.Mono(utils.GetGetChatsCommand()).Normal(" : ").Italic("The command which is used to get all chats in json.\n\n")
	txt.Mono(utils.GetSudoCommand()).Normal(" : ").Italic("The command which is used to promote sudo users in chats/channels.\n\n")
	txt.Mono(utils.GetHelpCommand()).Normal(" : ").Italic("Click on this to get more info about this command.\n\n")

	_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})
	return ext.EndGroups
}

func LoadHelpHandler(d *ext.Dispatcher, t []rune) {
	getHelpCommand := handlers.NewCommand(utils.GetHelpCommand(), HelpHandler)

	getHelpCommand.Triggers = t

	d.AddHandler(getHelpCommand)
}
