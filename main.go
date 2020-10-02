package main

import (
	"TGChannelGo/handlers/add"
	"TGChannelGo/handlers/getchats"
	"TGChannelGo/handlers/help"
	"TGChannelGo/handlers/remove"
	"TGChannelGo/handlers/send"
	"TGChannelGo/handlers/start"
	"TGChannelGo/utils"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

func RegisterAllHandlers(updater *gotgbot.Updater, l *zap.SugaredLogger) {
	start.LoadStartHandler(updater, l)
	send.LoadSendHandler(updater, l)
	getchats.LoadGetChatsHandler(updater, l)
	add.LoadAddHandler(updater, l)
	remove.LoadRemoveHandler(updater, l)
	help.LoadHelpHandler(updater, l)
}

func main() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel))
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()
	token := utils.GetBotToken()
	l.Info("Starting Bot.")
	l.Info("token: ", token)
	updater, err := gotgbot.NewUpdater(logger, token)
	l.Info("Got Updater")
	updater.UpdateGetter = ext.BaseRequester{
		Client: http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Second * 65,
		},
		ApiUrl: ext.ApiUrl,
	}
	updater.Bot.Requester = ext.BaseRequester{Client: http.Client{Timeout: time.Second * 65}}
	if err != nil {
		l.Fatalw("failed to start updater", zap.Error(err))
	}
	l.Info("Starting updater")
	RegisterAllHandlers(updater, l)
	_ = updater.StartPolling()
	l.Info("Started Updater.")
	updater.Idle()
}
