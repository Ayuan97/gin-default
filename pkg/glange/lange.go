package glange

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)
var langeBundle *i18n.Bundle
func Setup() {
	langeBundle = i18n.NewBundle(language.English)
	langeBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	//_, _ = bundle.LoadMessageFile("pkg/lange/active.yh.toml")
	lange := []string{
		"zh-Hans","zh-Hant","en","th","fr","es","fil","ms","pt","ja","id","af","am","bg","ca","hr","cs","da","nl","et","fi","de","el","he","hi","hu","is","it","ko","lv","lt","nb","pl","ro","ru","sr","sk","sl","sw","sv","tr","uk","vi","zu",
	}
	for _,v := range lange{
		langeBundle.MustLoadMessageFile("pkg/langefile/active."+v+".toml")
	}
}

func GetlangeMessage(lange string,messageKey string)string{
	localizer := i18n.NewLocalizer(langeBundle, lange)
	return localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageKey})
}