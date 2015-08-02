package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var cmd_file_map map[string]string

func get_file_content(cmd string) string {
	elem, ok := cmd_file_map[cmd]

	if ok {
		return elem
	}

	data, err := ioutil.ReadFile("msg_" + cmd + ".txt")
	if err != nil {
		fmt.Println("Fatal(0x1): Unable to load a content file of " + cmd)
		return "Fatal(0x1)"
	}

	str := string(data)
	cmd_file_map[cmd] = str
	return str
}

func parse(text string) (string, []string) {
	arr := strings.Split(text, " ")
	if len(arr) == 0 || arr[0][0] != '/' {
		return "", nil
	}

	return arr[0][1:], arr[1:]
}

func handler_about(bot *telebot.Bot, msg telebot.Message, args []string) {
	str := get_file_content("about")
	bot.SendMessage(msg.Chat, str, nil)
}

func handler_start(bot *telebot.Bot, msg telebot.Message, args []string) {
	str := get_file_content("start")
	bot.SendMessage(msg.Chat, str, nil)
}

func handler_help(bot *telebot.Bot, msg telebot.Message, args []string) {
	str := get_file_content("help")
	bot.SendMessage(msg.Chat, str, nil)
}

func handler_unknown(bot *telebot.Bot, msg telebot.Message, args []string) {
	str := get_file_content("unknown")
	bot.SendMessage(msg.Chat, str, nil)
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

	cmd_file_map = make(map[string]string)

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	fmt.Println("Info(0x0): Now listening")

	for msg := range messages {
		handler(bot, msg)
	}
}
