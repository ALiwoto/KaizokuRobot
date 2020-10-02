package remove

import (
	"ChannelReply/utils"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func RemoveHandler(b ext.Bot, u *gotgbot.Update) error {
	if !utils.IsUserSudo(u.EffectiveUser.Id) {
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "This command is sudo only", u.EffectiveMessage.MessageId)
		return nil
	}
	id := strings.Split(u.Message.Text, utils.GetRemoveCommand()+" ")[1]
	int_id, err := strconv.Atoi(id)
	if err != nil {
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Can't parse the given id, Please check it again", u.EffectiveMessage.MessageId)
		return nil
	}
	if !utils.IsChatInJson(int_id) {
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "User is not present in json", u.EffectiveMessage.MessageId)
		return err
	}
	err = utils.DelId(int_id)
	if err != nil {
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Got problems while deleting the given id", u.EffectiveMessage.MessageId)
		return err
	}
	_, _ = b.ReplyHTML(u.EffectiveChat.Id, fmt.Sprintf("Deleted id : <code>%v</code> from the json", int_id), u.EffectiveMessage.MessageId)
	return nil
}

func LoadRemoveHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Start Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetRemoveCommand(), RemoveHandler))
}
