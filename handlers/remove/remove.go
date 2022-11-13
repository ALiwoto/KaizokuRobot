package remove

import (
	"TGChannelGo/utils"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func RemoveHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
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

	if !utils.IsChatInJson(targetId) {
		txt := mdparser.GetNormal("User is not present in json.")
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	}

	err := utils.DelId(targetId)
	if err != nil {
		txt := mdparser.GetBold("Failed to delete the Id:\n")
		txt.Mono(err.Error())
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	}

	txt := mdparser.GetNormal("Deleted the id ").Mono(ssg.ToBase10(targetId))
	txt.Normal("from the json.")
	txt.Mono(err.Error())
	_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})
	return ext.EndGroups

}

func LoadRemoveHandler(d *ext.Dispatcher, t []rune) {
	removeCommand := handlers.NewCommand(utils.GetRemoveCommand(), RemoveHandler)

	removeCommand.Triggers = t

	d.AddHandler(removeCommand)
}
