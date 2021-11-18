package util

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var tgAlerter TgAlerter

var mailAlerter MailAlerter

type TgAlerter struct {
	BotId  string
	ChatId string
}

type MailAlerter struct {
	Sender          string
	SenderName      string
	SenderAuthToken string
	Receiver        string
	SMTPServer      string
	SMTPPort        int
}

func InitTgAlerter(cfg AlertConfig) {
	tgAlerter = TgAlerter{
		BotId:  cfg.TelegramBotId,
		ChatId: cfg.TelegramChatId,
	}
}

func InitMailAlerter(cfg AlertConfig) {
	senderAuthToken := cfg.MailSenderAuthToken
	exist := true
	if senderAuthToken == "" {
		senderAuthToken, exist = os.LookupEnv("SENDER_AUTH_TOKEN")
	}
	if !exist {
		Logger.Fatalf("can not find SENDER_AUTH_TOKEN/mail_sender_auth_token in config file or environment")
	}
	mailAlerter = MailAlerter{
		Sender:          cfg.MailSender,
		SenderName:      cfg.MailSender,
		SenderAuthToken: senderAuthToken,
		Receiver:        cfg.MailReceiver,
		SMTPServer:      cfg.MailSMTPServer,
		SMTPPort:        cfg.MailSMTPPort,
	}
}

func SendTelegramMessage(msg string) {
	if tgAlerter.BotId == "" || tgAlerter.ChatId == "" || msg == "" {
		return
	}
	msg = fmt.Sprintf("bsc-eth-swap-backend alert: %s", msg)
	endPoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tgAlerter.BotId)
	formData := url.Values{
		"chat_id":    {tgAlerter.ChatId},
		"parse_mode": {"html"},
		"text":       {msg},
	}
	Logger.Infof("send tg message, bot_id=%s, chat_id=%s, msg=%s", tgAlerter.BotId, tgAlerter.ChatId, msg)
	res, err := http.PostForm(endPoint, formData)
	if err != nil {
		Logger.Errorf("send telegram message error, bot_id=%s, chat_id=%s, msg=%s, err=%s", tgAlerter.BotId, tgAlerter.ChatId, msg, err.Error())
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		Logger.Errorf("read http response error, err=%s", err.Error())
		return
	}
	Logger.Infof("tg response: %s", string(bodyBytes))
}

func SendMailMessage(subject string, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mailAlerter.Sender)
	m.SetHeader("To", mailAlerter.Receiver)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	Logger.Warningf("please uncomment code to send mail...")
	//d := gomail.NewDialer(mailAlerter.SMTPServer, mailAlerter.SMTPPort, mailAlerter.SenderName, mailAlerter.SenderAuthToken)
	//if err := d.DialAndSend(m); err != nil {
	//	Logger.Errorf("dial and send mail error, err=%s", err.Error())
	//}
}
