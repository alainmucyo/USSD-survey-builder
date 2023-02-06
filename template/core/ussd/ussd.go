package ussd

import (
	"github.com/alainmucyo/ussd-go"
	"github.com/alainmucyo/ussd-go/sessionstores"
	"survey-ussd/core/survey"
)

type USSD struct {
	Ussd  *ussd.Ussd
	Store *sessionstores.Redis
}

func New(store *sessionstores.Redis) *USSD {
	u := &ussd.Ussd{}
	u = ussd.New("Survey", "Menu")
	u.Middleware(addData("global", "i'm here"))
	u.Ctrl(new(survey.Survey))
	return &USSD{Store: store, Ussd: u}
}

func addData(key string, value interface{}) ussd.Middleware {
	return func(c *ussd.Context) {
		c.Data[key] = value
	}
}
