package send

import (
	"bytes"
	"fmt"
	"github.com/AnimeKaizoku/KaizokuRobot/utils"
	"html"
	"os"
	"strings"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/AnimeKaizoku/ssg/ssg/rangeValues"
	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func MakeKeyboards(button []tg_md2html.Button) [][]gotgbot.InlineKeyboardButton {
	var allButtons [][]gotgbot.InlineKeyboardButton
	for _, v := range button {
		var button1 []gotgbot.InlineKeyboardButton
		tempBtn := gotgbot.InlineKeyboardButton{
			Text: v.Name,
			Url:  v.Content,
		}
		if !v.SameLine {
			button1 = append(button1, tempBtn)
			allButtons = append(allButtons, button1)
		} else {
			if allButtons == nil {
				button1 = append(button1, tempBtn)
				allButtons = append(allButtons, button1)
			} else {
				button1 = append(allButtons[len(allButtons)-1], tempBtn)
				allButtons = allButtons[:len(allButtons)-1]
				allButtons = append(allButtons, button1)
			}
		}
	}
	return allButtons
}

func GetImage(b *gotgbot.Bot, ctx *ext.Context) ([]byte, string) {
	message := ctx.EffectiveMessage
	imageLast := message.ReplyToMessage.Photo[len(message.ReplyToMessage.Photo)-1]
	imageFile, _ := b.GetFile(imageLast.FileId, nil)
	file, err := utils.DownloadFile(strings.Split(imageFile.FilePath, "/")[1], imageFile.FilePath)
	if err != nil {
		txt := mdparser.GetBold("Failed to download the image file:\n")
		txt.Mono(err.Error())
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return nil, ""
	}

	dat, _ := os.ReadFile(file)
	_ = os.Remove(strings.Split(imageFile.FilePath, "/")[1])
	name := strings.Split(imageFile.FilePath, "/")[1]
	return dat, name
}

func ReplayLinks(b *gotgbot.Bot, button []tg_md2html.Button) []tg_md2html.Button {
	var rtButton []tg_md2html.Button
	for _, v := range button {
		currentButton := tg_md2html.Button{
			Name:     v.Name,
			Content:  v.Content,
			SameLine: v.SameLine,
		}

		if strings.HasPrefix(v.Content, "*") {
			idStr := strings.Split(v.Content, "*")[1]
			container := rangeValues.ParseIntContainer[int64](idStr)
			if container == nil {
				continue
			}

			chat, _ := b.GetChat(container.Value, nil)
			if container.HasFlag("invite") {
				myLink, _ := b.CreateChatInviteLink(chat.Id, &gotgbot.CreateChatInviteLinkOpts{
					Name:               "@kaizoku's generated link",
					CreatesJoinRequest: true,
				})
				currentButton.Content = myLink.InviteLink
			} else {
				if chat.InviteLink == "" {
					chat.InviteLink, _ = chat.ExportInviteLink(b, nil)
				}
				currentButton.Content = chat.InviteLink
			}
		}
		rtButton = append(rtButton, currentButton)
	}
	return rtButton
}

func SendHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.EffectiveMessage
	user := ctx.EffectiveUser
	currentChat := ctx.EffectiveChat

	postLink := false
	approvalLink := false
	// var chat, group int64
	var sentId, groupId int64
	var text string
	//message = u.EffectiveMessage.Text
	if currentChat.Id != user.Id {
		return ext.ContinueGroups
	}

	if !utils.IsUserOwner(user.Id) && !utils.IsUserSudo(user.Id) {
		return nil
	}

	htmlText, btn := tg_md2html.MD2HTMLButtons(strings.SplitAfter(message.Text, utils.GetSendCommand())[1])
	text = htmlText

	if strings.HasSuffix(message.Text, "}") {
		postLink = false
		mssg := strings.Split(message.Text, "{")
		label := (strings.Split(mssg[len(mssg)-2], "}"))[0]
		groupStr := strings.Split((strings.Split(mssg[len(mssg)-1], "{"))[0], "}")[0]

		if strings.HasPrefix(groupStr, "*") {
			container := rangeValues.ParseIntContainer[int64](strings.Split(groupStr, "*")[1])
			if container == nil {
				postLink = false
			} else {
				postLink = true
				approvalLink = container.HasFlag("invite")
			}
			groupId = container.Value
		} else {
			groupId = ssg.ToInt64(strings.Split((strings.Split(mssg[len(mssg)-1], "{"))[0], "}")[0])
			postLink = false
		}

		chat, err := b.GetChat(groupId, nil)

		if err != nil {
			txt := mdparser.GetBold("Can't get the information of the respected chat you mentioned.\nPlease double check your chat/channel id you mentioned:\n")
			txt.Mono(err.Error())
			_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
				ParseMode: gotgbot.ParseModeMarkdownV2,
			})
			return ext.EndGroups
		}

		if postLink {
			if approvalLink {
				myLink, _ := b.CreateChatInviteLink(chat.Id, &gotgbot.CreateChatInviteLinkOpts{
					Name:               "@kaizoku's generated link",
					CreatesJoinRequest: true,
				})
				chat.InviteLink = myLink.InviteLink
			} else if len(chat.InviteLink) == 0 {
				_, err = chat.ExportInviteLink(b, nil)
				if err != nil {
					txt := mdparser.GetBold("Can't get the invite link of the chat you mentioned.\nPlease give me permission to generate invite links\n")
					txt.Mono(err.Error())
					_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
						ParseMode: gotgbot.ParseModeMarkdownV2,
					})
					return ext.EndGroups
				}
			}
			add := fmt.Sprintf("[%v](buttonurl:%v)", label, chat.InviteLink)
			htmlText, btn = tg_md2html.MD2HTMLButtons(strings.Split(strings.SplitAfter(message.Text, utils.GetSendCommand())[1], "{")[0] + "\n" + add)
			btn = ReplayLinks(b, btn)
			text = "<b>" + html.EscapeString(chat.Title) + "</b>" + "\n\n" + fmt.Sprintf("<code>%v</code>", chat.Id) + "\n\n" + htmlText
		} else {
			htmlText, btn = tg_md2html.MD2HTMLButtons(strings.Split(strings.SplitAfter(message.Text, utils.GetSendCommand())[1], "{")[0])
			text = "<b>" + html.EscapeString(chat.Title) + "</b>" + "\n\n" + fmt.Sprintf("<code>%v</code>", chat.Id) + "\n\n" + htmlText
			btn = ReplayLinks(b, btn)
		}
	}

	torep := ""

	if p, f := utils.GetStringInBetweenTwoString(message.Text, "{image:", ":imageend}"); f {
		msg, _ := b.SendPhoto(utils.GetChannelId(), p, &gotgbot.SendPhotoOpts{
			ParseMode: gotgbot.ParseModeHTML,
		})
		sentId = msg.MessageId
		torep = "####" + p + "$$$$"
	}

	if torep != "" {
		text = strings.ReplaceAll(text, torep, "")
	}

	if message.ReplyToMessage != nil && message.ReplyToMessage.Photo != nil {
		image, _ := GetImage(b, ctx)
		msg, _ := b.SendPhoto(utils.GetChannelId(), bytes.NewReader(image), &gotgbot.SendPhotoOpts{
			ParseMode: gotgbot.ParseModeHTML,
		})
		sentId = msg.MessageId
	}

	buttons := MakeKeyboards(btn)
	if sentId != 0 {
		_, _, err := b.EditMessageCaption(&gotgbot.EditMessageCaptionOpts{
			ChatId:      utils.GetChannelId(),
			MessageId:   sentId,
			Caption:     text,
			ParseMode:   gotgbot.ParseModeHTML,
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons},
		})

		if err != nil {
			txt := mdparser.GetBold("Error on sending message: \n")
			txt.Mono(err.Error())
			_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
				ParseMode: gotgbot.ParseModeMarkdownV2,
			})
			return ext.EndGroups
		}

		txt := mdparser.GetNormal("Message sent successfully!")
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
		return ext.EndGroups
	} else {
		_, err := b.SendMessage(utils.GetChannelId(), text, &gotgbot.SendMessageOpts{
			ReplyMarkup:           gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons},
			ParseMode:             gotgbot.ParseModeHTML,
			DisableWebPagePreview: true,
		})

		if err != nil {
			txt := mdparser.GetBold("Error on sending message: \n")
			txt.Mono(err.Error())
			_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
				ParseMode: gotgbot.ParseModeMarkdownV2,
			})
			return ext.EndGroups
		}

		txt := mdparser.GetNormal("Message sent successfully!")
		_, _ = message.Reply(b, txt.ToString(), &gotgbot.SendMessageOpts{
			ParseMode: gotgbot.ParseModeMarkdownV2,
		})
	}

	return ext.EndGroups
}

func LoadSendHandler(d *ext.Dispatcher, t []rune) {
	sendCommand := handlers.NewCommand(utils.GetSendCommand(), SendHandler)

	sendCommand.Triggers = t

	d.AddHandler(sendCommand)
}
