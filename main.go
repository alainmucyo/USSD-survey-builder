package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

type data struct {
	ProjectName string
	SurveyName  string
	Questions   string
}
type Question struct {
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Answers []Answer `json:"answers"`
}
type Answer struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func main() {
	q := []Question{
		Question{
			Id:    1,
			Title: "Question 1",
			Answers: []Answer{
				Answer{
					Id:    1,
					Title: "Answer 1.1",
				},
				Answer{
					Id:    2,
					Title: "Answer 1.2",
				},
			},
		},
		Question{
			Id:    2,
			Title: "Question 2",
			Answers: []Answer{
				Answer{
					Id:    3,
					Title: "Answer 2.1",
				},
				Answer{
					Id:    4,
					Title: "Answer 2.2",
				},
			},
		},
	}
	// questions to string
	jsonData, _ := json.Marshal(q)
	// jsondata to string
	jsonString := string(jsonData)
	d := data{
		ProjectName: "food-survey",
		SurveyName:  "Food survey",
		Questions:   jsonString,
	}
	t := template.Must(template.ParseGlob("template/main.go"))

	os.Mkdir(d.ProjectName, os.ModePerm)
	f, _ := os.Create(d.ProjectName + "/main.go")
	defer f.Close()
	t.ExecuteTemplate(f, "main.go", d)
	type Answer struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	}

	copyFiles("template/core", d.ProjectName+"/core")
	copyFiles("template/handlers", d.ProjectName+"/handlers")
	copyFiles("template/logs", d.ProjectName+"/logs")
	copyFiles("template/store", d.ProjectName+"/store")
	copyFile("template/go.mod", d.ProjectName+"/go.mod")
	copyFile("template/go.sum", d.ProjectName+"/go.sum")
	copyFile("template/.env", d.ProjectName+"/.env")

	template.Must(t.ParseGlob("template/core/survey/survey.go"))
	f, _ = os.Create(d.ProjectName + "/core/survey/survey.go")
	defer f.Close()
	t.ExecuteTemplate(f, "survey.go", d)

	println("hello world")
}

func copyFiles(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyFiles(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
