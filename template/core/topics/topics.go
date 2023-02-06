package topics

import (
	"survey-ussd/core/environment"
	ussdhandler "survey-ussd/handlers/ussd-handler"
)

type Topics struct {
	env         *environment.Environment
	ussdRequest *ussdhandler.Handler
	List        map[string]func([]byte, string)
}

func New(
	env *environment.Environment,
	ussdRequest *ussdhandler.Handler,

) *Topics {

	var topics = map[string]func([]byte, string){
		env.IncomingRequest: ussdRequest.HandleUssdRequestsKafka,
	}
	return &Topics{env: env, ussdRequest: ussdRequest, List: topics}
}
