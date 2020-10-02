package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const ConfigJsonPath string = "config.json"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ConfigJson struct {
	BOT_TOKEN  string `json:"bot_token"`
	SUDO_USERS []int  `json:"sudo_users"`
	OWNER_ID   int    `json:"owner_id"`
	CHANNEL_ID int    `json:"channel_id"`
	Monitor    []int  `json:"chats_to_monitor"`
}

var Config *ConfigJson = InitConfig()

func InitConfig() *ConfigJson {
	file, err := ioutil.ReadFile(ConfigJsonPath)
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

func AddId(id int) error {
	file, err := ioutil.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	if !Exists(ConfigJsonPath + ".bak") {
		err = ioutil.WriteFile(ConfigJsonPath+".bak", file, 0644)
	}
	var TempConfig ConfigJson
	err = json.Unmarshal(file, &TempConfig)
	if err != nil {
		log.Fatal(err)
	}
	TempConfig.Monitor = append(TempConfig.Monitor, id)
	Config.Monitor = TempConfig.Monitor
	newfile, _ := json.MarshalIndent(&TempConfig, "", "   ")
	err = ioutil.WriteFile(ConfigJsonPath, newfile, 0644)
	if err != nil {
		file, err := ioutil.ReadFile(ConfigJsonPath + ".bak")
		if err != nil {
			log.Println("Backup Config File Bad, exiting!")

		}
		ioutil.WriteFile(ConfigJsonPath, file, 0644)
		return err
	} else {
		return nil
	}
}

func DelId(id int) error {
	var index int
	file, err := ioutil.ReadFile(ConfigJsonPath)
	if err != nil {
		log.Fatal("Config File Bad, exiting!")
	}
	if !Exists(ConfigJsonPath + ".bak") {
		err = ioutil.WriteFile(ConfigJsonPath+".bak", file, 0644)
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
	TempConfig.Monitor = removeelement(TempConfig.Monitor, index)
	Config.Monitor = TempConfig.Monitor
	newfile, _ := json.MarshalIndent(&TempConfig, "", "   ")
	err = ioutil.WriteFile(ConfigJsonPath, newfile, 0644)
	if err != nil {
		file, err := ioutil.ReadFile(ConfigJsonPath + ".bak")
		if err != nil {
			log.Println("Backup Config File Bad, exiting!")

		}
		ioutil.WriteFile(ConfigJsonPath, file, 0644)
		return err
	} else {
		return nil
	}
}

func removeelement(s []int, i int) []int {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func IsUserOwner(userId int) bool {
	return Config.OWNER_ID == userId
}

func GetAllChats() []int {
	return Config.Monitor
}

func IsUserSudo(userId int) bool {
	for _, i := range Config.SUDO_USERS {
		if i == userId {
			return true
		}
	}
	return false
}

func GetChannelId() int {
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

func randSeq(n int) string {
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

func IsChatInJson(chatID int) bool {
	for _, i := range Config.Monitor {
		if i == chatID {
			return true
		}
	}
	InitConfig()
	return false
}
