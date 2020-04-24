package bot

import (
	"anipokev2/config"
	"fmt"
	tb "github.com/demget/telebot"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	StartCommand       = "/start"
	HelpCommand        = "/help"
	SubscribeCommand   = "/subscribe"
	UnsubscribeCommand = "/unsubscribe"
	HelpMessage        = "available command: \n" +
		"1. /subscribe\n" +
		"2. /unsubscribe\n" +
		"3./help to display all command\n" +
		"Remember that I didn't want to reply your group :p, just reach me through personal chat okay"

	UnknownMessageReply = "%v is unknown command, please see /help to view available one"
)

type Telegram struct {
	B   *tb.Bot
	C   *config.Config
	Rds *redis.Client
}

func NewTelegram(b *tb.Bot, c *config.Config, rds *redis.Client) *Telegram {
	return &Telegram{B: b, C: c, Rds: rds}
}

func (t *Telegram) OnCreate() {
	log.Println("starting telegram bot..")
}

func (t *Telegram) OnHelp() {
	t.B.Handle(StartCommand, func(m *tb.Message) {
		t.send(m, HelpMessage)
	})
	t.B.Handle(HelpCommand, func(m *tb.Message) {
		t.send(m, HelpMessage)
	})
}

func (t *Telegram) OnUnknown() {
	t.B.Handle(tb.OnText, func(m *tb.Message) {
		t.send(m, fmt.Sprintf(UnknownMessageReply, m.Text))
	})
}

func (t *Telegram) OnSubscribe() {
	t.B.Handle(SubscribeCommand, func(m *tb.Message) {
		if m.FromGroup() {
			return
		}
		if m.FromChannel() {
			return
		}

		t.Rds.HSet("anipoke:subscriber", strconv.Itoa(m.Sender.ID), "1")
		t.sendToUser(m.Sender, "subscribe success")
	})
}

func (t *Telegram) OnUnsubscribe() {
	t.B.Handle(UnsubscribeCommand, func(m *tb.Message) {
		t.Rds.HDel("anipoke:subscriber", strconv.Itoa(m.Sender.ID))
		t.sendToUser(m.Sender, "unsubscribe success")
	})
}

func (t *Telegram) send(m *tb.Message, message string) {
	m, err := t.B.Send(m.Sender, message)
	if err != nil {
		log.Println(err)
	}
}

func (t *Telegram) sendToUser(m *tb.User, message string) {
	_, err := t.B.Send(m, message)
	if err != nil {
		log.Println(err)
	}
}

func (t *Telegram) Broadcast(fansub, anime, link string) {
	subs, err := t.Rds.HGetAll("anipoke:subscriber").Result()
	if err != nil {
		log.Println(err)
	}
	for sub, _ := range subs {
		sid, err := strconv.Atoi(sub)
		if err != nil {
			log.Println(err)
		}
		t.sendToUser(&tb.User{
			ID: sid,
		}, "Update: "+fansub+"\n"+
			anime+"\n"+
			"Download now: "+link)
	}
}

func (t *Telegram) Run() {
	t.OnCreate()
	t.OnHelp()
	t.OnSubscribe()
	t.OnUnsubscribe()

	go t.B.Start()
}
