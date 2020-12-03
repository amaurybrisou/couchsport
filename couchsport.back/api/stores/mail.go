package stores

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"path/filepath"

	"github.com/amaurybrisou/couchsport.back/api/models"
	"github.com/amaurybrisou/couchsport.back/localizer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type mailStore struct {
	Email, Password, Server string
	Port                    int
	Localizer               *localizer.Localizer
}

func (me *mailStore) sendMail(mail models.Mail) error {
	smtpServer := fmt.Sprintf("%s:%d", me.Server, me.Port)

	if err := smtp.SendMail(smtpServer, smtp.PlainAuth("", me.Email, me.Password, me.Server), me.Email, mail.To, []byte(mail.Body)); err != nil {
		return err
	}

	return nil
}

func (me *mailStore) sendMailTLS(mail models.Mail) error {

	smtpServer := fmt.Sprintf("%s:%d", me.Server, me.Port)

	auth := smtp.PlainAuth("", me.Email, me.Password, me.Server)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         me.Server,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", smtpServer, tlsconfig)
	if err != nil {
		log.Error(err)
		return err
	}

	c, err := smtp.NewClient(conn, me.Server)
	if err != nil {
		log.Error(err)
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Error(err)
		return err
	}

	// To && From
	if err = c.Mail(mail.From); err != nil {
		log.Error(err)
		return err
	}

	if err = c.Rcpt(mail.To[0]); err != nil {
		log.Error(err)
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = w.Write([]byte(mail.Body))
	if err != nil {
		log.Error(err)
	}

	err = w.Close()
	if err != nil {
		log.Error(err)
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	log.Printf("email sent to %s", mail.To[0])

	return nil
}

func (me *mailStore) send(mail models.Mail, tls bool) error {

	if tls {
		if err := me.sendMailTLS(mail); err != nil {
			return err
		}
	} else {
		if err := me.sendMail(mail); err != nil {
			return err
		}
	}

	return nil
}

//AccountAutoCreated send the password
func (me *mailStore) AccountAutoCreated(email, password, locale string) {
	log.Printf("sending 'AccountAutoCreated' email to %s", email)

	fileName := "account_auto_created.html"
	template := filepath.Join(viper.GetString("MAIL_TEMPLATE_PATH"), fileName)

	mail := models.NewMail(
		me.Email,
		[]string{email},
		me.Localizer.Translate("account_auto_created.title", locale, nil),
	)

	body, err := me.Localizer.ParseTemplateI18n(fileName, template, locale, map[string]string{"email": email, "password": password})
	if err != nil {
		log.Error(err)
		return
	}

	headers := mail.GetHeaders()
	mail.Body = headers + body

	if err := me.send(*mail, true); err != nil {
		log.Error(err)
	}
}
