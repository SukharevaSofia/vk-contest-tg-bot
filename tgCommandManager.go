package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const telegramPollTimeout = 3600

var telegramBaseUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_API_KEY")
var postSendResponseUrl = telegramBaseUrl + "/sendMessage"

func tgManageRequests(db *sql.DB) {
	log.Default().Println("tgManageRequests has been called")
	getUpdatesURL := fmt.Sprintf(
		"%s/getUpdates?timeout=%d",
		telegramBaseUrl,
		telegramPollTimeout)

	userRequestsResponse, err := http.Get(getUpdatesURL)
	log.Println("Getting the requests from users")
	if err != nil {
		fmt.Println("Failed to get users' requests")
		os.Exit(1)
	}

	userRequestsBody, _ := io.ReadAll(userRequestsResponse.Body)

	var requestList TelegramRequests
	var sendMessageFields url.Values

	_ = json.Unmarshal(userRequestsBody, &requestList)

	var offset = 0
	for _, value := range requestList.Result {
		if value.UpdateId > offset {
			offset = value.UpdateId
		}
		textMessage := NO_SUCH_COMMAND
		switch value.Message.Text {
		case "/start":
			kBoard := menuKboard()
			textMessage = TEXT_BOT_DESCRIPTION
			sendMessageFields = url.Values{
				"chat_id":      {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":         {textMessage},
				"reply_markup": {kBoard}}
		case NEW_ENTRY:
			textMessage = MOOD_CHECK
			kBoard := moodKboard()
			sendMessageFields = url.Values{
				"chat_id":      {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":         {textMessage},
				"reply_markup": {kBoard}}
		case SHOW_ENTRIES:
			createTable(db, value.Message.Chat.Id)
			textMessage = getDataFromDb(value.Message.Chat.Id)
			sendMessageFields = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		case ABOUT_DEV:
			textMessage = DEV_INF
			sendMessageFields = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		case COMMAND_LIST:
			textMessage = HELP
			sendMessageFields = url.Values{
				"chat_id":    {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":       {textMessage},
				"parse_mode": {"MarkdownV2"}}
		default:
			if isMood(value.Message.Text) {
				log.Println("СОЗДАНИЕ ТАБЛИЦЫ ", value.Message.Chat.Id)
				createTable(db, value.Message.Chat.Id)
				addToDb(value.Message.Chat.Id, value.Message.Text)
				textMessage = INFO_SAVED
				kBoard := menuKboard()
				sendMessageFields = url.Values{
					"chat_id":      {fmt.Sprintf("%d", value.Message.From.Id)},
					"text":         {textMessage},
					"reply_markup": {kBoard}}
				break
			}
			sendMessageFields = url.Values{
				"chat_id": {fmt.Sprintf("%d", value.Message.From.Id)},
				"text":    {textMessage}}
		}
		response, _ := http.PostForm(postSendResponseUrl, sendMessageFields)
		log.Println("request ", value.Message.Text)
		log.Println("response :", sendMessageFields)
		log.Println("has sent with status " + string(response.Status))
	}

	// drop processed messages
	if offset > 0 {
		_, _ = http.Get(fmt.Sprintf("%s/getUpdates?offset=%d&limit=1", telegramBaseUrl, offset+1))
	}
}
