package getchats

import (
	"ChannelReply/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func GetChats(b ext.Bot, u *gotgbot.Update) error {
	var message string
	message = "<b>All chats list</b>\n\n"
	if !utils.IsUserOwner(u.EffectiveUser.Id) {
		if !utils.IsUserSudo(u.EffectiveUser.Id) {
			return nil
		}
	}
	chats := utils.GetAllChats()
	for _, i := range chats {
		chat, _ := b.GetChat(i)
		if len(chat.Username) != 0 {
			message += fmt.Sprintf("<b>%v</b>\n%v\n<code>%v</code>\n", chat.Title, chat.Username, chat.Id)
		} else {
			message += fmt.Sprintf("<b>%v</b>\n<code>%v</code>\n\n", chat.Title, chat.Id)
		}
	}
	_, _ = b.ReplyHTML(u.EffectiveChat.Id, message, u.EffectiveMessage.MessageId)
	return nil
}

func LoadGetChatsHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("GetChats Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetGetChatsCommand(), GetChats))
}
