package output

import (
	tb "github.com/demget/telebot"
)

type TelegramChannel struct {
	chatId string
	token  string
}

func NewTelegramChannel(token string, chatId string) *TelegramChannel {
	return &TelegramChannel{token: token, chatId: chatId}
}

func (t TelegramChannel) Write(p []byte) (n int, err error) {
	botInstance, err := tb.NewBot(tb.Settings{
		Token: t.token,
	})

	if err != nil {
		return 0, err
	}

	channel, err := botInstance.ChatByID(t.chatId)
	if err != nil {
		return 0, err
	}

	if _, err := botInstance.Send(channel, string(p), &tb.SendOptions{
		DisableWebPagePreview: false,
		DisableNotification:   false,
	}); err != nil {
		return 0, err
	}

	return len(p), nil
}
