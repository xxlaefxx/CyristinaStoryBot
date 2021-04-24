package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	story "github.com/xxlaefxx/CyristinaStoryBot/internal/story"
)

func GenerateKeyboard(titles []string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton
	for _, title := range titles {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(title, title))
	}
	var storiesKeyboard = tgbotapi.NewInlineKeyboardMarkup(buttons)
	return storiesKeyboard
}

func main() {
	var stories = story.ReadAllStories()
	var titles = story.GetTitles(stories)

	var allStories map[string]story.Story = make(map[string]story.Story)

	for _, t := range titles {
		for _, s := range stories {
			if t == s.Title {
				allStories[t] = s
			}
		}
	}

	bot, err := tgbotapi.NewBotAPI("1781364855:AAGJJqx0pjhCWG_GqpPTuQgZKmZIhwc9Yh4")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Bot started : %s", bot.Self.UserName)
	log.Printf("Stories titles: %#v", titles)

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {
			log.Printf("Callback Query: %s", update.CallbackQuery.Data)
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))

			st, ok := allStories[update.CallbackQuery.Data]
			if !ok {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Я не знаю такой сказки. Вот что я знаю:")
				msg.ReplyMarkup = GenerateKeyboard(titles)
				bot.Send(msg)
				continue
			} else {
				storyMessages := story.GenerateMessagesForStory(update.CallbackQuery.Message.Chat.ID, st)
				for _, sm := range storyMessages {
					_, err = bot.Send(sm)
					if err != nil {
						log.Print(err)
					}
				}
				continue
			}
		}

		if update.Message.Command() != "" {
			log.Printf("Command: %s", update.Message.Command())
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Выбирай сказку:")
				msg.ReplyMarkup = GenerateKeyboard(titles)
				bot.Send(msg)
				continue
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите /start чтобы получить список сказок")
				bot.Send(msg)
				continue
			}
		}

		if update.Message == nil {
			log.Printf("Message is nil. Skipping")
			continue
		}

		log.Printf("Just a message: [%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите /start чтобы получить список сказок")
		bot.Send(msg)
	}
}
