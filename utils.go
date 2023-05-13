package main

import (
	"encoding/json"
	"time"
)

type TelegramRequests struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateId int `json:"update_id"`
		Message  struct {
			MessageId int `json:"message_id"`
			From      struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date     int    `json:"date"`
			Text     string `json:"text"`
			Entities []struct {
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				Type   string `json:"type"`
			} `json:"entities,omitempty"`
		} `json:"message"`
	} `json:"result"`
}

const (
	TEXT_BOT_DESCRIPTION string = "С помощью этого бота вы можете следить за своими эмоциями в удобном формате."
	NO_SUCH_COMMAND             = "Такой команды не существует :("
	MOOD_CHECK                  = "Как вы себя чувствуете?"
	INFO_SAVED                  = "Данные были записаны."
	NEW_ENTRY                   = "Новая запись"
	SHOW_ENTRIES                = "Показать записи"
	DEV_INF                     = "Спасибо за ваш интерес к проекту!" +
		"\nGithub: https://github.com/SukharevaSofia" +
		"\nTelegram: @Topinamburka" +
		"\nMail: work@sssofya.ru"
	ABOUT_DEV    = "О разработчике"
	COMMAND_LIST = "Список команд"
	NO_DATA      = "В таблице пока нет данных!"
	DATABASE_URL = "DATABASE_URL"
	HELP         = "Доступные команды:\n" +
		"\n`Новая запись` — добавляет новую запись в дневник наблюдений" +
		"\n`Показать записи` — выводит ранее введённые в дневник данные" +
		"\n`О разработчике` — контактная информация" +
		"\n`Список команд` — выводит список доступных команд"
)

const (
	HAPPY  string = "Отлично 😄"
	GOOD          = "Хорошо 😌"
	MEH           = "Так себе 😐"
	UNSURE        = "Не уверен(а) 🧐"
	SAD           = "Грустинка 😔"
	BIGSAD        = "Большая грустинка 😭"
	DEAD          = "потихоньку 💀"
	CRAZY         = "хихихаха 🤪"
)

func isMood(s string) bool {
	if s == HAPPY || s == GOOD || s == MEH || s == UNSURE || s == SAD || s == BIGSAD || s == DEAD || s == CRAZY {
		return true
	}
	return false
}

func menuKboard() string {
	var keyboard = NewReplyKeyboard(
		NewKeyboardButtonRow(
			NewKeyboardButton(NEW_ENTRY),
			NewKeyboardButton(SHOW_ENTRIES),
		),
		NewKeyboardButtonRow(
			NewKeyboardButton(ABOUT_DEV),
			NewKeyboardButton(COMMAND_LIST),
		),
	)
	kboard, _ := json.Marshal(keyboard)
	return string(kboard)
}

func moodKboard() string {
	var keyboard = NewReplyKeyboard(
		NewKeyboardButtonRow(
			NewKeyboardButton(HAPPY),
			NewKeyboardButton(GOOD),
		),
		NewKeyboardButtonRow(
			NewKeyboardButton(MEH),
			NewKeyboardButton(UNSURE),
		),
		NewKeyboardButtonRow(
			NewKeyboardButton(SAD),
			NewKeyboardButton(BIGSAD),
		),
		NewKeyboardButtonRow(
			NewKeyboardButton(CRAZY),
			NewKeyboardButton(DEAD),
		),
	)

	kboard, _ := json.Marshal(keyboard)
	return string(kboard)
}

type dateMoodPair struct {
	date time.Time
	mood string
}

// Source of the keyboard structs and their functions https://github.com/go-telegram-bot-api/

func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	var keyboard [][]KeyboardButton
	keyboard = append(keyboard, rows...)
	return ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact,omitempty"`
	RequestLocation bool                    `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
}

type KeyboardButtonPollType struct {
	Type string `json:"type"`
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

func NewKeyboardButtonRow(buttons ...KeyboardButton) []KeyboardButton {
	var row []KeyboardButton
	row = append(row, buttons...)
	return row
}

func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}
