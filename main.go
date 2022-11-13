package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ALiwoto/mdparser/mdparser"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/add"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/getchats"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/help"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/remove"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/send"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/start"
	"github.com/AnimeKaizoku/KaizokuRobot/handlers/sudo"
	"github.com/AnimeKaizoku/KaizokuRobot/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func RegisterAllHandlers(d *ext.Dispatcher, triggers []rune) {
	sudo.LoadSudoHandler(d, triggers)
	start.LoadStartHandler(d, triggers)
	send.LoadSendHandler(d, triggers)
	getchats.LoadGetChatsHandler(d, triggers)
	add.LoadAddHandler(d, triggers)
	remove.LoadRemoveHandler(d, triggers)
	help.LoadHelpHandler(d, triggers)
}

func StartBot() error {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zap.InfoLevel))
	defer logger.Sync() // flushes buffer, if any
	l := logger.Sugar()
	token := utils.GetBotToken()

	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: 6 * gotgbot.DefaultTimeout,
		},
	})
	if err != nil {
		return err
	}

	mdparser.AddSecret(token, "$TOKEN")
	uTmp := ext.NewUpdater(nil)
	updater := &uTmp
	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
	})
	if err != nil {
		return err
	}

	l.Info(fmt.Sprintf("%s has started | ID: %d", b.Username, b.Id))

	RegisterAllHandlers(updater.Dispatcher, []rune{'/', '!'})

	updater.Idle()
	return nil
}

func main() {
	err := StartBot()
	if err != nil {
		log.Fatal(err)
	}

}
