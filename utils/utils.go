package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AnimeKaizoku/ssg/ssg"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

const ConfigJsonPath string = "config.json"
const CommandConfigPath = "commands.json"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ConfigJson struct {
	BOT_TOKEN  string  `json:"bot_token"`
	SUDO_USERS []int64 `json:"sudo_users"`
	OWNER_ID   int64   `json:"owner_id"`
	CHANNEL_ID int64   `json:"channel_id"`
	Monitor    []int64 `json:"chats_to_monitor"`
}

type CommandJson struct {
	START    string `json:"start"`
	HELP     string `json:"help"`
	GETCHATS string `json:"getchats"`
	ADD      string `json:"add"`
	REMOVE   string `json:"remove"`
	SEND     string `json:"send"`
	SUDO     string `json:"sudo"`
}

var CommandConfig *CommandJson = InitCommandConfig()
var Config *ConfigJson = InitConfig()

func InitCommandConfig() *CommandJson {
	file, err := os.ReadFile(CommandConfigPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}

	var Config CommandJson
	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		log.Fatal(err)
	}
	return &Config
}

func InitConfig() *ConfigJson {
	file, err := os.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}

	var Config ConfigJson
	err = json.Unmarshal([]byte(file), &Config)
	if err != nil {
		log.Fatal(err)
	}
	Config.SUDO_USERS = append(Config.SUDO_USERS, Config.OWNER_ID)
	log.Println(Config.SUDO_USERS)
	return &Config
}

func GetBotToken() string {
	return Config.BOT_TOKEN
}

func AddId(id int64) error {
	file, err := os.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	if !Exists(ConfigJsonPath + ".bak") {
		_ = os.WriteFile(ConfigJsonPath+".bak", file, 0644)
	}
	var TempConfig ConfigJson
	err = json.Unmarshal(file, &TempConfig)
	if err != nil {
		log.Fatal(err)
	}
	TempConfig.Monitor = append(TempConfig.Monitor, id)
	Config.Monitor = TempConfig.Monitor
	newFile, _ := json.MarshalIndent(&TempConfig, "", "   ")
	err = os.WriteFile(ConfigJsonPath, newFile, 0644)
	if err != nil {
		file, err := os.ReadFile(ConfigJsonPath + ".bak")
		if err != nil {
			log.Println("Backup Config File Bad, exiting!")

		}
		os.WriteFile(ConfigJsonPath, file, 0644)
		return err
	} else {
		return nil
	}
}

func DelId(id int64) error {
	var index int
	file, err := os.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	if !Exists(ConfigJsonPath + ".bak") {
		_ = os.WriteFile(ConfigJsonPath+".bak", file, 0644)
	}
	var TempConfig ConfigJson
	err = json.Unmarshal(file, &TempConfig)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range TempConfig.Monitor {
		if v == id {
			index = i
		}
	}
	fmt.Println(index)
	fmt.Println(TempConfig.Monitor)
	TempConfig.Monitor = removeElement(TempConfig.Monitor, index)
	Config.Monitor = TempConfig.Monitor
	newFile, _ := json.MarshalIndent(&TempConfig, "", "   ")
	err = os.WriteFile(ConfigJsonPath, newFile, 0644)
	if err != nil {
		file, err := os.ReadFile(ConfigJsonPath + ".bak")
		if err != nil {
			log.Println("Backup Config File Bad, exiting!")

		}
		os.WriteFile(ConfigJsonPath, file, 0644)
		return err
	} else {
		return nil
	}
}

func removeElement(s []int64, i int) []int64 {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func IsUserOwner(userId int64) bool {
	return Config.OWNER_ID == userId
}

func GetAllChats() []int64 {
	return Config.Monitor
}

func IsUserSudo(userId int64) bool {
	for _, i := range Config.SUDO_USERS {
		if i == userId {
			return true
		}
	}
	return false
}

func GetChannelId() int64 {
	return Config.CHANNEL_ID
}

func FormatTGFileLink(sub string, token string) string {
	return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, sub)
}

func DownloadFile(fileName string, sub string) (string, error) {
	link := FormatTGFileLink(sub, GetBotToken())
	log.Printf("[Download] %s", link)
	resp, err := http.Get(link)
	rand.Seed(time.Now().UnixNano())
	if err != nil {
		return fileName, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		return fileName, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return fileName, err
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsChatInJson(chatID int64) bool {
	for _, i := range Config.Monitor {
		if i == chatID {
			return true
		}
	}

	return false
}

func GetStartCommand() string {
	return CommandConfig.START
}

func GetHelpCommand() string {
	return CommandConfig.HELP
}
func GetGetChatsCommand() string {
	return CommandConfig.GETCHATS
}
func GetSendCommand() string {
	return CommandConfig.SEND
}

func GetAddCommand() string {
	return CommandConfig.ADD
}

func GetRemoveCommand() string {
	return CommandConfig.REMOVE
}

func GetSudoCommand() string {
	return CommandConfig.SUDO
}

func GetStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

func GetChat(chatId string, b *gotgbot.Bot) (*gotgbot.Chat, error) {
	v := map[string]string{}
	theId := ssg.ToInt64(chatId)
	if theId == 0 {
		if !strings.HasPrefix(chatId, "@") {
			chatId = "@" + chatId
		}
	}

	v["chat_id"] = chatId

	r, err := b.Request("getChat", v, nil, nil)
	if err != nil {
		return nil, err
	}

	c := gotgbot.Chat{}
	return &c, json.Unmarshal(r, &c)
}
