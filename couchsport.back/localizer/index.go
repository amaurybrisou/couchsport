package localizer

import (
	"bytes"
	"html/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

//Localizer definition
type Localizer struct {
	bundle *i18n.Bundle
}

//NewLocalizer return the app localizer
func NewLocalizer(languageFiles []string) *Localizer {
	bundle := i18n.NewBundle(language.French)

	for _, l := range languageFiles {
		bundle.LoadMessageFile(l)
	}

	return &Localizer{
		bundle: bundle,
	}
}

func (me Localizer) t(locale string) func(message string) string {
	localizer := i18n.NewLocalizer(me.bundle, locale)
	return func(message string) string {
		mess, err := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: message,
			},
		})

		if err != nil {
			return message
		}
		return mess
	}
}

//Translate a string into the locale specified with variables required to compile the message
func (me *Localizer) Translate(message, locale string, variables map[string]string) string {
	localizer := i18n.NewLocalizer(me.bundle, locale)
	mess, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: message,
		},
		TemplateData: variables,
	})

	if err != nil {
		log.Errorf("localize error : %s", err)
		return message
	}

	return mess
}

//ParseTemplateI18n translate template in the specified locale
func (me *Localizer) ParseTemplateI18n(fileName, filePath, locale string, templateVars map[string]string) (string, error) {
	t, err := template.New(fileName).Funcs(template.FuncMap{
		"T": me.t(locale),
	}).ParseFiles(filePath)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, templateVars); err != nil {
		return "", err
	}

	return buffer.String(), nil

}
