package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

var telegramBaseUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_KEY")
var postSendResponseUrl = telegramBaseUrl + "/sendMessage"

const telegramPollTimeout = 3600

func tgGetRequests() {
	log.Default().Println("tgGetRequests has been called")

	postUserUrl := fmt.Sprintf(
		"%s/getUpdates?timeout=%d",
		telegramBaseUrl,
		telegramPollTimeout)
	userRequestsRes, err := http.Get(postUserUrl)
	log.Println("Getting the requests from users")
	if err != nil {
		fmt.Println("Failed to get users' requests")
		os.Exit(1)
	}

	userRequestsBody, _ := io.ReadAll(userRequestsRes.Body)

	var requestList TelegramRequests
	var postSendMessage url.Values

	_ = json.Unmarshal(userRequestsBody, &requestList)

	var offset = 0
	for _, value := range requestList.Result {
		if value.UpdateId > offset {
			offset = value.UpdateId
		}

		textMessage := "Такой команды не существует :("

		switch value.Message.Text {
		case "/start":
			kBoard := startMenu()
			textMessage = "С помощью этого бота вы можете следить за своими эмоциями в удобном формате."
			postSendMessage = url.Values{
				"chat_id":      {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":         {textMessage},
				"reply_markup": {kBoard}}
		case "Новая запись":
			//TODO
			textMessage = "Hi " + value.Message.From.Username + "!"
			postSendMessage = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		case "Показать записи":
			//TODO
			textMessage = "Hello " + value.Message.From.Username + "!"
			postSendMessage = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		case "Больше о разработчике":
			//TODO
			textMessage = "Dev info"
			postSendMessage = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		default:
			postSendMessage = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		}
		a, _ := http.PostForm(postSendResponseUrl, postSendMessage)
		log.Println("request ", value.Message.Text)
		log.Println("response :", postSendMessage)
		log.Println("has sent with status " + string(a.Status))
	}

	// drop processed messages
	if offset > 0 {
		_, _ = http.Get(fmt.Sprintf("%s/getUpdates?offset=%d&limit=1", telegramBaseUrl, offset+1))
	}
}

func startMenu() string {
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новая запись"),
			tgbotapi.NewKeyboardButton("Показать записи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Больше о разработчике"),
		),
	)
	kboard, _ := json.Marshal(keyboard)
	return string(kboard)
}
