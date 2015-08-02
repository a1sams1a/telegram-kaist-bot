package main

import (
	"encoding/json"
	"fmt"
	"github.com/tucnak/telebot"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type BotData struct {
	data    string
	created time.Time
}

// SECTION - GLOBAL VARIABLE
// store message formats of each command
var msg_map map[string]string

// store data for some commands
var data_map map[string]BotData

// SECTION - HELPER FUNCTIONS
// parse given text to cmd and args
func parse(text string) (string, []string) {
	arr := strings.Split(text, " ")
	if len(arr) == 0 || arr[0][0] != '/' {
		return "", nil
	}

	return arr[0][1:], arr[1:]
}

// get message using given cmd key and args
func msgfmt(key string, args ...interface{}) string {
	return fmt.Sprintf(msg_map[key], args...)
}

// SECTION - DATA FUNCTION
// update river temp
func update_river() BotData {
	data, err := ioutil.ReadFile("extdata_river.txt")
	if err != nil {
		fmt.Println("Fatal(0x4): Unable to read external data file")
		return BotData{data: "Fatal(0x4)", created: time.Now()}
	}

	arr := strings.Split(string(data), ",")
	if len(arr) < 3 {
		fmt.Println("Fatal(0x5): Format of external data file is not expected")
		return BotData{data: "Fatal(0x5)", created: time.Now()}
	}

	rdata := BotData{data: msgfmt("river", arr[0], arr[1], arr[2]), created: time.Now()}
	data_map["river"] = rdata
	return rdata
}

// SECTION - HANDLER FUNCTION
// about - display bot infos
func handler_about(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("about"), nil)
}

// start - display init message
func handler_start(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("start"), nil)
}

// help - display help message
func handler_help(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("help"), nil)
}

// unknown - default handler for unknown command
func handler_unknown(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("unknown"), nil)
}

// rand - get random integer in [0, n)
func handler_rand(bot *telebot.Bot, msg telebot.Message, args []string) {
	if len(args) == 0 {
		bot.SendMessage(msg.Chat, msgfmt("rand_noarg"), nil)
		return
	}

	num, err := strconv.Atoi(args[0])
	if err != nil || num <= 0 {
		bot.SendMessage(msg.Chat, msgfmt("rand_invarg"), nil)
		return
	}

	result := rand.Intn(num)
	bot.SendMessage(msg.Chat, msgfmt("rand", num, result), nil)
}

// river - get temp information of Hangang and Gapchun
func handler_river(bot *telebot.Bot, msg telebot.Message, args []string) {
	rdata, ok := data_map["river"]
	if !ok || time.Since(rdata.created).Hours() >= 2 {
		rdata = update_river()
	}
	bot.SendMessage(msg.Chat, rdata.data, nil)
}

// SECTION - MAIN
// top-level handler for KAIST BOT
func handler(bot *telebot.Bot, msg telebot.Message) {
	cmd, args := parse(msg.Text)
	switch cmd {
	case "start":
		handler_start(bot, msg, args)
	case "about":
		handler_about(bot, msg, args)
	case "help":
		handler_help(bot, msg, args)
	case "rand":
		handler_rand(bot, msg, args)
	case "river":
		handler_river(bot, msg, args)
	default:
		handler_unknown(bot, msg, args)
	}
}

// bot execution logic, init logic
func main() {
	bot, err := telebot.NewBot(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		fmt.Println("Fatal(0x0): Unable to start bot")
		return
	}

	data_map = make(map[string]BotData)
	msg_map = make(map[string]string)
	data, err := ioutil.ReadFile("msg.txt")
	if err != nil {
		fmt.Println("Fatal(0x1): Unable to read msg file")
		return
	}

	err = json.Unmarshal(data, &msg_map)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Fatal(0x2): Unable to parse msg file")
		return
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	fmt.Println("Info(0x0): Now listening")

	for msg := range messages {
		handler(bot, msg)
	}
}
