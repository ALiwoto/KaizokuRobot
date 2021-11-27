package send

import (
	"TGChannelGo/utils"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	tg_md2html "github.com/PaulSonOfLars/gotg_md2html"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func MakeKeyboards(button []tg_md2html.Button) [][]ext.InlineKeyboardButton {
	var allbuttons [][]ext.InlineKeyboardButton
	for _, v := range button {
		var button1 []ext.InlineKeyboardButton
		tempbtn := ext.InlineKeyboardButton{
			Text: v.Name,
			Url:  v.Content,
		}
		if !v.SameLine {
			button1 = append(button1, tempbtn)
			allbuttons = append(allbuttons, button1)
		} else {
			if allbuttons == nil {
				button1 = append(button1, tempbtn)
				allbuttons = append(allbuttons, button1)
			} else {
				button1 = append(allbuttons[len(allbuttons)-1], tempbtn)
				allbuttons = allbuttons[:len(allbuttons)-1]
				allbuttons = append(allbuttons, button1)
			}
		}
	}
	return allbuttons
}

func GetImage(b ext.Bot, u *gotgbot.Update) ([]byte, string) {
	imageLast := u.EffectiveMessage.ReplyToMessage.Photo[len(u.EffectiveMessage.ReplyToMessage.Photo)-1]
	imageFile, _ := b.GetFile(imageLast.FileId)
	file, err := utils.DownloadFile(strings.Split(imageFile.FilePath, "/")[1], imageFile.FilePath)
	if err != nil {
		_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Can't download the image file", u.Message.MessageId)
		fmt.Println(err)
	}
	dat, _ := ioutil.ReadFile(file)
	_ = os.Remove(strings.Split(imageFile.FilePath, "/")[1])
	name := strings.Split(imageFile.FilePath, "/")[1]
	return dat, name
}

func ReplaeLinks(b ext.Bot, button []tg_md2html.Button) []tg_md2html.Button {
	var rtbutton []tg_md2html.Button
	for _, v := range button {
		lol := tg_md2html.Button{
			Name:     v.Name,
			Content:  v.Content,
			SameLine: v.SameLine,
		}
		if strings.HasPrefix(v.Content, "*") {
			id := strings.Split(v.Content, "*")[1]
			intId, err := strconv.Atoi(id)
			if err != nil {
				continue
			}
			chat, _ := b.GetChat(intId)
			if len(chat.InviteLink) < 0 {
				_, _ = chat.ExportInviteLink()
			}
			lol.Content = chat.InviteLink
		}
		rtbutton = append(rtbutton, lol)
	}
	return rtbutton
}

func SendHandler(b ext.Bot, u *gotgbot.Update) error {
	var postlink bool
	postlink = false
	var chat, group int
	var text string
	message := u.EffectiveMessage.Text
	if u.EffectiveChat.Id != u.EffectiveUser.Id {
		return nil
	}
	if !utils.IsUserOwner(u.EffectiveUser.Id) {
		if !utils.IsUserSudo(u.EffectiveUser.Id) {
			return nil
		}
	}

	html, btn := tg_md2html.MD2HTMLButtons(strings.SplitAfter(message, utils.GetSendCommand())[1])
	text = html

	if strings.HasSuffix(message, "}") {
		postlink = false
		mssg := strings.Split(message, "{")
		label := (strings.Split(mssg[len(mssg)-2], "}"))[0]
		groupStr := strings.Split((strings.Split(mssg[len(mssg)-1], "{"))[0], "}")[0]

		if strings.HasPrefix(groupStr, "*") {
			lol, err := strconv.Atoi(strings.Split(groupStr, "*")[1])
			postlink = true
			if err != nil {
				postlink = false
			}
			group = lol
		} else {
			group, _ = strconv.Atoi(strings.Split((strings.Split(mssg[len(mssg)-1], "{"))[0], "}")[0])
			postlink = false
		}

		chat, err := b.GetChat(group)

		if err != nil {
			_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Can't get the information of the respected chat you mentioned, Please double check your chat/channel id you mentioned", u.Message.MessageId)
			return err
		}

		if postlink {
			if len(chat.InviteLink) == 0 {
				_, err = chat.ExportInviteLink()
				if err != nil {
					_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Can't get the invite link of the respected chat you mentioned, Please give me permission to generate invite links", u.Message.MessageId)
				}
			}
			add := fmt.Sprintf("[%v](buttonurl:%v)", label, chat.InviteLink)
			html, btn = tg_md2html.MD2HTMLButtons(strings.Split(strings.SplitAfter(message, utils.GetSendCommand())[1], "{")[0] + "\n" + add)
			btn = ReplaeLinks(b, btn)
			text = "<b>" + chat.Title + "</b>" + "\n\n" + fmt.Sprintf("<code>%v</code>", chat.Id) + "\n\n" + html
		} else {
			html, btn = tg_md2html.MD2HTMLButtons(strings.Split(strings.SplitAfter(message, utils.GetSendCommand())[1], "{")[0])
			text = "<b>" + chat.Title + "</b>" + "\n\n" + fmt.Sprintf("<code>%v</code>", chat.Id) + "\n\n" + html
			btn = ReplaeLinks(b, btn)
		}
	}

	torep := ""

	if p, f := utils.GetStringInBetweenTwoString(message, "{image:", ":imageend}"); f {
		rq := b.NewSendablePhoto(utils.GetChannelId(), "")
		rq.Photo = b.NewFileURL(p)
		rq.ParseMode = "HTML"
		msg, _ := rq.Send()
		chat = msg.MessageId
		torep = "####" + p + "$$$$"
	}

	if torep != "" {
		text = strings.ReplaceAll(text, torep, "")
	}
	if u.EffectiveMessage.ReplyToMessage != nil && u.EffectiveMessage.ReplyToMessage.Photo != nil {
		image, name := GetImage(b, u)
		lol := b.NewSendablePhoto(utils.GetChannelId(), "")
		lol.Photo = b.NewFileReader(name, bytes.NewReader(image))
		lol.ParseMode = "HTML"
		msg, _ := lol.Send()
		chat = msg.MessageId
	}
	buttons := MakeKeyboards(btn)
	if chat != 0 {
		msg := b.NewSendableEditMessageCaption(utils.GetChannelId(), chat, text)
		msg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &buttons}
		msg.ParseMode = "HTML"
		_, err := msg.Send()
		if err == nil {
			_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Message sent successfully!", u.EffectiveMessage.MessageId)
		} else {
			_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, fmt.Sprintf("Error sending message : %v", err.Error()), u.EffectiveMessage.MessageId)
		}
	} else {
		msg := b.NewSendableMessage(utils.GetChannelId(), text)
		msg.ReplyMarkup = &ext.InlineKeyboardMarkup{InlineKeyboard: &buttons}
		msg.ParseMode = "HTML"
		msg.DisableWebPreview = true
		_, err := msg.Send()
		if err == nil {
			_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, "Message sent successfully!", u.EffectiveMessage.MessageId)
		} else {
			_, _ = b.ReplyMarkdownV2(u.EffectiveChat.Id, fmt.Sprintf("Error sending message : %v", err.Error()), u.EffectiveMessage.MessageId)
		}
	}
	return nil
}

func LoadSendHandler(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	defer l.Info("Send Module Loaded.")
	updater.Dispatcher.AddHandler(handlers.NewCommand(utils.GetSendCommand(), SendHandler))
}
