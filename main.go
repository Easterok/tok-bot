package main

import (
	"fmt"
	"log"
	"os"
	"telegram-onboarding-kit-bot/db"
	"telegram-onboarding-kit-bot/utils"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/joho/godotenv"
)

var client = db.Dbconnect()

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		panic("TOKEN environment variable is empty")
	}

	b, err := gotgbot.NewBot(token, nil)

	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling update:", err.Error())

			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	success, err := b.SetMyCommands(
		[]gotgbot.BotCommand{
			{Command: "/new", Description: "Create a new onboarding"},
			{Command: "/examples", Description: "Show examples of onboarding"},
			{Command: "/list", Description: "Show your onboardings"},
		},
		nil,
	)

	if !success || err != nil {
		log.Fatal("Unable to set commands: %w", err)
	}

	dispatcher.AddHandler(handlers.NewCommand("start", startCommandHandler))
	dispatcher.AddHandler(handlers.NewCommand("new", newCommandHandler))
	dispatcher.AddHandler(handlers.NewCommand("examples", examplesCommandHandler))
	dispatcher.AddHandler(handlers.NewCommand("list", listCommandHandler))

	dispatcher.AddHandler(handlers.NewMessage(webappData, onMessage))

	err = updater.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})

	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func registerUserIfNotExist(b *gotgbot.Bot, ctx *ext.Context, user *gotgbot.User) {
	if !client.CheckIfUserExist(user.Id) {
		client.AddNewUser(user.Id, ctx.EffectiveChat.Id, user.Username, user.FirstName, user.LastName)
	}
}

var START_MESSAGE = "♥️ Hi! I'm bot for <a href='https://github.com/Easterok/telegram-onboarding-kit'>Telegram Onboarding Kit</a>\n\nYou can create your onboarding using this bot by using the command /new"

func getUserData(user *gotgbot.User) map[string]string {
	return map[string]string{
		"language_code": user.LanguageCode,
	}
}

func startCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	registerUserIfNotExist(b, ctx, ctx.EffectiveUser)

	_, err := ctx.EffectiveChat.SendMessage(b, START_MESSAGE, &gotgbot.SendMessageOpts{
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	})

	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	return examplesHandler(b, ctx)
}

func examplesCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	registerUserIfNotExist(b, ctx, ctx.EffectiveUser)

	return examplesHandler(b, ctx)
}

var EXAMPLES_TEXT = "Below you can see demo onboardings <b>created with our kit</b>.\nIt's better to you watch them from 📱 mobile device"

func onboardingsKeyboard(userData map[string]string) [][]gotgbot.KeyboardButton {
	return [][]gotgbot.KeyboardButton{
		{
			{Text: "🌈 Base onboarding", WebApp: &gotgbot.WebAppInfo{Url: utils.AddQueryToUrl("https://easterok.github.io/telegram-onboarding-kit", userData)}},
		},
		{
			{Text: "💃 Fashion AI", WebApp: &gotgbot.WebAppInfo{Url: utils.AddQueryToUrl("https://tok-ai.netlify.app", userData)}},
		},
		{
			{Text: "🔐 VPN", WebApp: &gotgbot.WebAppInfo{Url: utils.AddQueryToUrl("https://tok-vpn.netlify.app", userData)}},
		},
		{
			{Text: "🧠 ChatGPT", WebApp: &gotgbot.WebAppInfo{Url: utils.AddQueryToUrl("https://tok-chatgpt.netlify.app", userData)}},
		},
	}
}

func examplesHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	userData := getUserData(ctx.EffectiveUser)

	_, err := ctx.EffectiveChat.SendMessage(b, EXAMPLES_TEXT, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			Keyboard: onboardingsKeyboard(userData),
		},
	})

	if err != nil {
		return fmt.Errorf("failed to send examples message: %w", err)
	}

	return nil
}

func listCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	registerUserIfNotExist(b, ctx, ctx.EffectiveUser)

	return nil
}

var NEW_COMMAND_TEXT = "To create a new onboarding tap on a button below"

func newCommandHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	registerUserIfNotExist(b, ctx, ctx.EffectiveUser)

	_, err := ctx.EffectiveChat.SendMessage(b, NEW_COMMAND_TEXT, &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{
				{Text: "New onboarding", WebApp: &gotgbot.WebAppInfo{Url: "https://resonant-cactus-0fefb5.netlify.app?edit=1"}},
			}},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to send new message: %w", err)
	}

	return nil
}

func webappData(msg *gotgbot.Message) bool {
	return msg.WebAppData != nil
}

func onMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage.WebAppData.Data

	_, err := ctx.EffectiveChat.SendMessage(b, fmt.Sprintf("Got message from miniapp %s", msg), nil)

	if err != nil {
		return fmt.Errorf("failed to send start message: %w", err)
	}

	return nil
}
