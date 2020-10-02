package help

import (
	"ChannelReply/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func HelpHandler(b ext.Bot, u *gotgbot.Update) error {
	if !utils.IsUserOwner(u.EffectiveUser.Id) {
		if !utils.IsUserSudo(u.EffectiveUser.Id) {
			return nil
		}
	}
	message := fmt.Sprintf("<b>Help Section of ChannelGoBot</b>\n\n"+
		"Available commands are as follows: \n\n"+
		"<code>/%v</code> : <i>The all boring start command.</i>\n\n"+
		"<code>/%v</code> : <i>The actual useful send command used for sending messages to channels.</i>\n\n"+
		"<code>/%v</code> : <i>The add command for adding new chat id to json.</i>\n\n"+
		"<code>/%v</code> : <i>The remove command for removing chat id from json.</i>\n\n"+
		"<code>/%v</code> : <i>The command which is used to get all chats in json.</i>\n\n"+
		"/%v : <i>Click on this to get more info about this command.</i>", utils.GetStartCommand(), utils.GetSendCommand(), utils.GetAddCommand(), utils.GetRemoveCommand(), utils.GetGetChatsCommand(), utils.GetHelpCommand())
	b.ReplyHTML(u.EffectiveChat.Id, message, u.EffectiveMessage.MessageId)
	return nil
}

func LoadHelpHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Start Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetHelpCommand(), HelpHandler))
}
