package main

import (
	"bytes"
	"log"
	"time"

	lp "github.com/wangii/littlepool"
)

type MyTask1 struct {
	id string
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
	_ = getTemplate().Execute(w, a)
	log.Print(w.String())

	return lp.TaskResultSuccess
}

func (t *MyTask1) Next() lp.Task {
	return newMyTask2(t.id + "-S2")
}
