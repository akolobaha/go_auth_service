package service

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/gomail.v2"
)

type ConfigEmail struct {
	SmtpHost string `env:"SMTP_HOST" env-default:"smtp.yandex.ru"`
	SmtpPort int    `env:"SMTP_PORT" env-default:"465"`
	SmtpUser string `env:"SMTP_USER"`
	SmtpPass string `env:"SMTP_PASS"`
}

func sendEmail(to string, subject string, body string) {
	var cfg ConfigEmail

	err := cleanenv.ReadConfig("../.env", &cfg)
	if err != nil {
		return
	}

	fmt.Println(cfg)

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SmtpUser)
	m.SetHeader("To", to) // Укажите email получателя
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUser, cfg.SmtpPass)
	d.SSL = true

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
