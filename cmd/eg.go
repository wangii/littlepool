package main

import (
	"fmt"

	lp "github.com/wangii/littlepool"
)

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
