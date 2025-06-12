package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"log"
	"time"

	lp "github.com/wangii/littlepool"
)

type MyTask1 struct {
	id string
}

//go:embed prompts/task1.md
var templateText string
var tpTask1 *template.Template

func init() {
	log.Print(templateText)
	tpTask1 = template.Must(template.New("task1").Parse(templateText))
}

func newMyTask1(id string) *MyTask1 {
	return &MyTask1{id: id}
}

func (t *MyTask1) ID() string {
	return t.id
}

func (t *MyTask1) GetPoolID() string {
	return "pool1"
}

func (t *MyTask1) Run() lp.TaskResult {
	time.Sleep(time.Second * 5)

	a := struct {
		Name string
	}{
		Name: "hu jing",
	}

	w := bytes.NewBufferString("")
	_ = tpTask1.Execute(w, a)
	log.Print(w.String())

	return lp.TaskResultSuccess
}

func (t *MyTask1) Next() []lp.Task {
	return []lp.Task{newMyTask2(t.id + "-S2")}
}
