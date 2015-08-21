package main

import (
	"encoding/json"
	"fmt"
	"github.com/tucnak/telebot"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
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

// split data first and then fmt
func datafmt(key string, args string) string {
	arr := strings.Split(string(args), ",")

	dt := make([]interface{}, len(arr))
	for i := range arr {
		dt[i] = arr[i]
	}
	return msgfmt(key, dt...)
}

// SECTION - DATA FUNCTION
// check update time
func should_update(o time.Time, expire string) bool {
	switch expire {
	case "d":
		return time.Now().Day() != o.Day()
	default:
		num, err := strconv.Atoi(expire)
		return err == nil && time.Since(o).Hours() >= float64(num)
	}
	return false
}

// read data file and get the content
func get_data(name string) string {
	rdata, ok := data_map[name]
	if ok {
		return rdata.data
	}

	fmt.Println("Info(0x2): Try to read data of " + name)
	data, err := ioutil.ReadFile("data_" + name + ".txt")
	if err != nil {
		fmt.Println("Fatal(0x7): Unable to read data file")
		return ""
	}

	data_map[name] = BotData{data: string(data), created: time.Now()}
	return data_map[name].data
}

// check update status of external data
func check_external(name string, expire string) bool {
	rdata, ok := data_map[name]
	if ok && !should_update(rdata.created, expire) {
		return true
	}

	script_name := strings.Split(name, "_")[0]
	fmt.Println("Info(0x1): Execute external update script of " + script_name)
	_, err := exec.Command("./update_" + script_name + ".py").Output()
	if err != nil {
		fmt.Println("Fatal(0x4): Unable to execute update script")
		return false
	}
	return true
}

// SECTION - HANDLER FUNCTION
// about - display bot infos
func handler_about(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, get_data("about"), nil)
}

// start - display init message
func handler_start(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, get_data("start"), nil)
}

// help - display help message
func handler_help(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, get_data("help"), nil)
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

// food - get menu of north, east, west cafe
func handler_food(bot *telebot.Bot, msg telebot.Message, args []string) {
	if len(args) == 0 {
		bot.SendMessage(msg.Chat, msgfmt("food_noarg"), nil)
		return
	}

	if args[0] != "n" && args[0] != "w" && args[0] != "e" {
		bot.SendMessage(msg.Chat, msgfmt("food_invarg"), nil)
		return
	}

	iname := "food_" + args[0]
	if !check_external(iname, "d") {
		bot.SendMessage(msg.Chat, "Fatal(0x4)", nil)
		return
	}

	bot.SendMessage(msg.Chat, msgfmt(iname, get_data(iname)), nil)
}

// river - get temp information of Hangang and Gapchun
func handler_river(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendChatAction(msg.Chat, "typing")

	if !check_external("river", "2") {
		bot.SendMessage(msg.Chat, "Fatal(0x4)", nil)
		return
	}

	bot.SendMessage(msg.Chat, datafmt("river", get_data("river")), nil)
}

// weather - get weather information of Daejon
func handler_weather(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendChatAction(msg.Chat, "typing")

	if !check_external("weather", "2") {
		bot.SendMessage(msg.Chat, "Fatal(0x4)", nil)
		return
	}

	bot.SendMessage(msg.Chat, datafmt("weather", get_data("weather")), nil)
}

// store - get store opening time
func handler_store(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, get_data("store"), nil)
}

// loc - search building number by name
func handler_loc(bot *telebot.Bot, msg telebot.Message, args []string) {
	if len(args) == 0 {
		bot.SendMessage(msg.Chat, msgfmt("loc_noarg"), nil)
		return
	}

	data := get_data("loc")
	arr := strings.Split(data, ",")
	result := ""
	for _, v := range arr {
		if strings.Contains(v, args[0]) {
			result += v + "\n"
		}
	}
	bot.SendMessage(msg.Chat, msgfmt("loc", result), nil)
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
	case "weather":
		handler_weather(bot, msg, args)
	case "store":
		handler_store(bot, msg, args)
	case "loc":
		handler_loc(bot, msg, args)
	case "food":
		handler_food(bot, msg, args)
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
