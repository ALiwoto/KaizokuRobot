package start

import (
	"ChannelReply/utils"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func StartHandler(b ext.Bot, u *gotgbot.Update) error {
	if u.EffectiveChat.Id != u.EffectiveUser.Id {
		return nil
	}
	if !utils.IsUserOwner(u.EffectiveUser.Id) {
		if !utils.IsUserSudo(u.EffectiveUser.Id) {
			return nil
		}
	}
	_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Yes bot is active and you can use it", u.EffectiveMessage.MessageId)

	return nil
}

func LoadStartHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Start Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand("start", StartHandler))
}
