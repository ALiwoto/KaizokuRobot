package getchats

import (
	"github.com/AnimeKaizoku/KaizokuRobot/utils"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func GetChatsHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
		return ext.ContinueGroups
	}

	md := mdparser.GetBold("All chats list\n\n")
	chats := utils.GetAllChats()
	for _, currentId := range chats {
		chat, _ := b.GetChat(currentId, nil)
		if chat == nil {
			// lets do use a hacky way here to include a dummy chat, instead of
			// making the bot panic.
			chat = &gotgbot.Chat{
				Id:    currentId,
				Title: "Unknown",
			}
		}

		if len(chat.Username) != 0 {
			md.Bold(chat.Title).Normal("\n" + chat.Username + "\n").Mono(ssg.ToBase10(chat.Id)).ElThis()
		} else {
			md.Bold(chat.Title).Mono("\n" + ssg.ToBase10(chat.Id)).Normal("\n\n")
		}
	}

	_, _ = message.Reply(b, md.ToString(), &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeMarkdownV2,
	})
	return ext.EndGroups
}

func LoadGetChatsHandler(d *ext.Dispatcher, t []rune) {
	getChatsCommand := handlers.NewCommand(utils.GetGetChatsCommand(), GetChatsHandler)

	getChatsCommand.Triggers = t

	d.AddHandler(getChatsCommand)
}
