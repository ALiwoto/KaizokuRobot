package add

import (
	"TGChannelGo/utils"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/ssg/ssg"
	_ "github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func AddHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	if !utils.IsUserSudo(ctx.EffectiveUser.Id) {
		return ext.EndGroups
	}

	idStr := strings.Split(message.Text, utils.GetAddCommand()+" ")[1]
	targetId := ssg.ToInt64(idStr)
	if targetId == 0 {
		txt := mdparser.GetNormal("Can't parse the given id, please check it again.")
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	}

	if utils.IsChatInJson(targetId) {
		txt := mdparser.GetNormal("User is already present in json.")
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	}

	err := utils.AddId(targetId)
	if err != nil {
		txt := mdparser.GetBold("Failed to add the Id:\n")
		txt.Mono(err.Error())
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	}

	txt := mdparser.GetNormal("Added the id ").Mono(ssg.ToBase10(targetId))
	txt.Normal("to the json.")
	txt.Mono(err.Error())
	_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})
	return ext.EndGroups
}

func LoadAddHandler(d *ext.Dispatcher, t []rune) {
	addCommand := handlers.NewCommand(utils.GetAddCommand(), AddHandler)

	addCommand.Triggers = t

	d.AddHandler(addCommand)
}
