package survey

import (
	"encoding/json"
	"github.com/alainmucyo/ussd-go"
	"survey-ussd/logs"
)

type Survey struct {
}

type Answer struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Question struct {
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Answers []Answer `json:"answers"`
}

var questionIndexes = map[string]int{}

var questions []Question

func (s Survey) Menu(c *ussd.Context) ussd.Response {
	go logs.AppLog("showing home menu ", c.Data["trackId"].(string), "Menu()", c.Request)
	_ = json.Unmarshal([]byte(`{{.Questions}}`), &questions)
	questionIndexes[c.Request.SessionId] = 0
	menu := ussd.NewMenu()
	menu.Add("Welcome to {{.SurveyName}}!\n\n1. Start answering", "Survey", "AnswerQuestion")
	menu.Add("2. Language", "Survey", "LanguageMenu")
	menu.AddZero("Exit", "Survey", "Exit")

	return c.RenderMenu(menu)
}

func (s Survey) AnswerQuestion(c *ussd.Context) ussd.Response {
	// Check if it is the latest question, then shows the end menu
	if questionIndexes[c.Request.SessionId] == len(questions) {
		menu := ussd.NewMenu()
		menu.Add("Thank you for answering the survey", "Survey", "Menu")
		menu.AddZero("Exit", "Survey", "Exit")
		return c.RenderMenu(menu)
	}
	question := questions[questionIndexes[c.Request.SessionId]]
	menu := ussd.NewMenu()
	menu.Add(question.Title, "Survey", "AnswerQuestion")
	for _, answer := range question.Answers {
		menu.Add(answer.Title, "Survey", "AnswerQuestion")
	}
	questionIndexes[c.Request.SessionId]++
	return c.RenderMenu(menu)

}
