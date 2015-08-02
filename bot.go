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

var msg_map map[string]string

func parse(text string) (string, []string) {
	arr := strings.Split(text, " ")
	if len(arr) == 0 || arr[0][0] != '/' {
		return "", nil
	}

	return arr[0][1:], arr[1:]
}

func msgfmt(key string, args ...interface{}) string {
	return fmt.Sprintf(msg_map[key], args...)
}

func handler_about(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("about"), nil)
}

func handler_start(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("start"), nil)
}

func handler_help(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("help"), nil)
}

func handler_unknown(bot *telebot.Bot, msg telebot.Message, args []string) {
	bot.SendMessage(msg.Chat, msgfmt("unknown"), nil)
}

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
	default:
		handler_unknown(bot, msg, args)
	}
}

func main() {
	bot, err := telebot.NewBot(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		fmt.Println("Fatal(0x0): Unable to start bot")
		return
	}

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
