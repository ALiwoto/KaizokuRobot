package sudo

import (
	"TGChannelGo/utils"
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func SelfPromoteHandler(b ext.Bot, u *gotgbot.Update) error {
	if !utils.IsUserOwner(u.EffectiveUser.Id) {
		if !utils.IsUserSudo(u.EffectiveUser.Id) {
			return gotgbot.EndGroups{}
		}
	}
	args := strings.Split(u.EffectiveMessage.Text, " ")
	if u.EffectiveChat.Type == "private" && len(args) == 1 {
		_, err := u.EffectiveMessage.ReplyText("You can't promote yourself in a private chat")
		return err
	}
	if len(args) > 1 {
		en := args[1]
		chat, err := utils.GetChat(en, b)
		if err != nil {
			_, err := u.EffectiveMessage.ReplyText(fmt.Sprintf("Error for %s: %s", en, err.Error()))
			return err
		}
		if chat.Type != "private" {
			err := promote(chat, u.EffectiveUser.Id, b)
			if err != nil {
				_, err := u.EffectiveMessage.ReplyText(fmt.Sprintf("Error: %s", err.Error()))
				return err
			}
		}

	} else {
		err := promote(u.EffectiveChat, u.EffectiveUser.Id, b)
		if err != nil {
			_, err := u.EffectiveMessage.ReplyText(fmt.Sprintf("Error: %s", err.Error()))
			return err
		}
	}
	_, err := u.EffectiveMessage.ReplyText("Promoted!")
	return err
}

func promote(chat *ext.Chat, userId int, b ext.Bot) error {
	selfMem, err := b.GetChatMember(chat.Id, b.Id)
	if err != nil {
		return err
	}
	rq := b.NewSendablePromoteChatMember(chat.Id, userId)
	rq.CanChangeInfo = selfMem.CanChangeInfo
	rq.CanDeleteMessages = selfMem.CanDeleteMessages
	rq.CanInviteUsers = selfMem.CanInviteUsers
	rq.CanPinMessages = selfMem.CanPinMessages
	rq.CanPromoteMembers = selfMem.CanPromoteMembers
	_, err = rq.Send()
	if err != nil {
		return err
	}
	return nil
}

func LoadSudoHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Sudo Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand("sudo", SelfPromoteHandler))
}
