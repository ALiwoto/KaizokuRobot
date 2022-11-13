package sudo

import (
	"TGChannelGo/utils"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func SelfPromoteHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	user := ctx.EffectiveUser
	message := ctx.EffectiveMessage
	chat := ctx.EffectiveChat
	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
		return ext.EndGroups
	}

	args := strings.Split(message.Text, " ")
	if chat.Type == "private" && len(args) == 1 {
		_, err := message.Reply(b, "You can't promote yourself in a private chat", nil)
		return err
	}

	if len(args) > 1 {
		targetChat := args[1]
		chat, err := utils.GetChat(targetChat, b)
		if err != nil {
			txt := mdparser.GetNormal("Error for " + targetChat + ": \n")
			txt.Mono(err.Error())
			_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
				ParseMode: gotgbot.ParseModeMarkdownV2,
			})
			return err
		}

		if chat.Type != "private" {
			err := promote(b, chat, user.Id)
			if err != nil {
				txt := mdparser.GetNormal("Error when promoting in " + targetChat + ": \n")
				txt.Mono(err.Error())
				_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
					ParseMode: gotgbot.ParseModeMarkdownV2,
				})
				return err
			}
		}

	} else {
		err := promote(b, chat, user.Id)
		if err != nil {
			txt := mdparser.GetNormal("Error when promoting in " + ssg.ToBase10(chat.Id) + ": \n")
			txt.Mono(err.Error())
			_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
				ParseMode: gotgbot.ParseModeMarkdownV2,
			})
			return ext.EndGroups
		}
	}

	_, _ = message.Reply(b, "Promoted!", nil)
	return ext.EndGroups
}

func promote(b *gotgbot.Bot, chat *gotgbot.Chat, userId int64) error {
	selfMem, err := b.GetChatMember(chat.Id, b.Id, nil)
	if err != nil {
		return err
	}

	me := selfMem.MergeChatMember()

	_, err = chat.PromoteMember(b, userId, &gotgbot.PromoteChatMemberOpts{
		IsAnonymous:         me.IsAnonymous,
		CanManageChat:       me.CanManageChat,
		CanPostMessages:     me.CanPostMessages,
		CanEditMessages:     me.CanEditMessages,
		CanDeleteMessages:   me.CanDeleteMessages,
		CanManageVideoChats: me.CanManageVideoChats,
		CanRestrictMembers:  me.CanRestrictMembers,
		CanPromoteMembers:   me.CanPromoteMembers,
		CanChangeInfo:       me.CanChangeInfo,
		CanInviteUsers:      me.CanInviteUsers,
		CanPinMessages:      me.CanPinMessages,
	})

	return err
}

func LoadSudoHandler(d *ext.Dispatcher, t []rune) {
	sudoCommand := handlers.NewCommand(utils.GetSudoCommand(), SelfPromoteHandler)

	sudoCommand.Triggers = t

	d.AddHandler(sudoCommand)
}
