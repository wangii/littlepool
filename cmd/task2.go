package main

import (
	"log"
	"time"

	lp "github.com/wangii/littlepool"
)

type MyTask2 struct {
	id string
}

func newMyTask2(id string) *MyTask2 {
	return &MyTask2{id: id}
}

func (t *MyTask2) ID() string {
	return t.id
}

func (t *MyTask2) GetPoolID() string {
	return "pool2"
}

func (t *MyTask2) Run() lp.TaskResult {
	time.Sleep(time.Second * 5)
	log.Printf("task2: %s", t.id)
	return lp.TaskResultSuccess
}

func (t *MyTask2) Next() []lp.Task {
	return []lp.Task{}
}
