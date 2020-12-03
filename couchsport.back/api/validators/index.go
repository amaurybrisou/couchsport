package validators

import (
	"github.com/asaskevich/govalidator"
	"regexp"
)

//Init validators
func Init() {
	govalidator.SetFieldsRequiredByDefault(false)

	govalidator.TagMap["name"] = govalidator.Validator(func(str string) bool {
		re := regexp.MustCompile("^[a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,.'-]+$")
		return re.MatchString(str)
	})

	govalidator.TagMap["zipcode"] = govalidator.Validator(func(str string) bool {
		re := regexp.MustCompile("^[a-zA-Z0-9- ]+$")
		return re.MatchString(str)
	})

	govalidator.TagMap["text"] = govalidator.Validator(func(str string) bool {
		re := regexp.MustCompile("^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]+$")
		return re.MatchString(str)
	})
}
