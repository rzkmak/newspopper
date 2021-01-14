package output

import tb "github.com/demget/telebot"

type Telegram struct {
	token  string
	chatId string
	parser string
}

const (
	ParserMd   = "markdown"
	ParserHtml = "html"
)

func NewTelegram(token string, chatId string) *Telegram {
	return &Telegram{token: token, chatId: chatId}
}

func (t Telegram) Write(p []byte) (n int, err error) {
	bot := tb.Bot{
		Token: t.token,
	}
	channel, err := bot.ChatByID(t.chatId)
	if err != nil {
		return 0, err
	}
	parseMode := ""
	if t.parser == ParserHtml {
		parseMode = ParserHtml
	}

	if t.parser == ParserMd {
		parseMode = ParserMd
	}

	if _, err := bot.Send(channel, p, tb.SendOptions{
		DisableWebPagePreview: false,
		DisableNotification:   false,
		ParseMode:             parseMode,
	}); err != nil {
		return 0, err
	}

	return len(p), nil
}
