package main

import (
	"embed"
	"fmt"
	"log"
	"path"
	"runtime"
	"strings"
	"text/template"

	lp "github.com/wangii/littlepool"
)

var (
	//go:embed prompts/*.md
	_promptsFS embed.FS
	prompts    map[string]*template.Template
)

func getTemplate() *template.Template {
	_, fn, _, _ := runtime.Caller(1)
	fn = strings.TrimSuffix(path.Base(fn), ".go")
	if prompts[fn] == nil {
		log.Fatal("No prompt template found for " + fn)
	}

	return prompts[fn]
}

func init() {

	files, err := _promptsFS.ReadDir("prompts")
	if err != nil {
		log.Fatal(err)
	}

	prompts = make(map[string]*template.Template)
	for _, file := range files {
		fileName := file.Name()
		log.Println(fileName)
		if !strings.HasSuffix(fileName, ".md") {
			continue
		}

		tplName := strings.TrimSuffix(fileName, ".md")
		fileContent, err := _promptsFS.ReadFile("prompts/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
		tpl, err := template.New(tplName).Parse(string(fileContent))
		if err != nil {
			log.Fatal(err)
		}
		prompts[tplName] = tpl
	}
	log.Println(prompts)
}

func main() {
	c := lp.NewController[*MyTask2](lp.PoolConfig{ConcurrencyLimit: 5, ID: "pool1"},
		lp.PoolConfig{ConcurrencyLimit: 2, ID: "pool2"},
	)

	for i := range 10 {
		c.Add(newMyTask1(fmt.Sprintf("task1-%d", i)))
	}

	c.Start()

	// c.IterateFinished(func(t *MyTask2) {
	// 	fmt.Println(t.id)
	// })

	for _, t := range c.GetFinished() {
		fmt.Println(t.ID())
	}
}
