package ussd_handler

import (
	"encoding/json"
	ussd2 "github.com/alainmucyo/ussd-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"survey-ussd/core/environment"
	"survey-ussd/core/service"
	"survey-ussd/core/ussd"
	"survey-ussd/logs"
	"survey-ussd/store/kafka/producer"
)

type Handler struct {
	ussd     *ussd.USSD
	producer *producer.Producer
	env      *environment.Environment
	service  *service.Service
}

func New(ussd *ussd.USSD, producer *producer.Producer, env *environment.Environment, surveyService *service.Service) *Handler {
	return &Handler{ussd: ussd, producer: producer, env: env, service: surveyService}
}

func (h *Handler) HandleUssdRequestsKafka(message []byte, topic string) {
	var value USSDRequest
	err := json.Unmarshal(message, &value)
	go logs.AppLog("Received USSD request", value.TrackId, "HandleUssdRequestsKafka()", value)

	if err != nil {
		println(err.Error())
		return
	}

	if value.Sequence == 1 {
		value.Text = "*" + value.ServiceCode + "#"
	}

	data := ussd2.Data{"trackId": value.TrackId, "service": h.service}
	res := Response{}
	defer func() {
		if err := recover(); err != nil {
			println("Survived a panic:==============================")
			value.Response = "Something went wrong"
		}
	}()
	h.ussd.Ussd.Process(h.ussd.Store, data, value, &res)
	if res.Message != "" {
		value.Response = res.Message

	}
	value.Action = res.Action
	value.Tag = "ussd response"
	go logs.AppLog("Sending response back to gateway ", value.TrackId, "Menu()", value)
	h.producer.Produce(value, h.env.OutgoingResponse)
}

func (h *Handler) HandleUSSDRequests(c *gin.Context) {
	var ussdRqDTO USSDRequest
	if err := c.BindJSON(&ussdRqDTO); err != nil {
		println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON object",
		})
		return
	}
	if ussdRqDTO.Sequence == 1 {
		ussdRqDTO.Text = "*" + ussdRqDTO.ServiceCode + "#"
	}
	ussdRqDTO.Tag = "ussd response"
	data := ussd2.Data{"trackId": ussdRqDTO.TrackId, "service": h.service}
	res := Response{}
	h.ussd.Ussd.Process(h.ussd.Store, data, ussdRqDTO, &res)
	ussdRqDTO.Response = res.Message
	ussdRqDTO.Action = res.Action
	ussdRqDTO.Text = ""
	c.JSON(http.StatusOK, ussdRqDTO)
}

type USSDRequest struct {
	SessionId   string            `json:"sessionId"`
	TrackId     string            `json:"trackId,omitempty"`
	Text        string            `json:"text"`
	PhoneNumber string            `json:"phoneNumber"`
	ServiceCode string            `json:"serviceCode" binding:"-"`
	Action      string            `json:"action" binding:"-"`
	Response    string            `json:"response,omitempty" binding:"-"`
	Tag         string            `json:"tag,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Sequence    int               `json:"sequence"`
}

type Response struct {
	Message string
	Action  string
}

func (r *Response) SetResponse(response ussd2.Response) {
	if response.Release {
		r.Action = "END"
	} else {
		r.Action = "CON"
	}
	r.Message = response.Message
}

func (s USSDRequest) GetRequest() *ussd2.Request {
	return &ussd2.Request{
		Action:      s.Action,
		Text:        s.Text,
		PhoneNumber: s.PhoneNumber,
		SessionId:   s.SessionId,
		ServiceCode: s.ServiceCode,
	}
}
