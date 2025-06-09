package main

import (
	"fmt"

	lp "github.com/wangii/littlepool"
)

func main() {
	c := lp.NewController(lp.PoolConfig{ConcurrencyLimit: 5, ID: "pool1"},
		lp.PoolConfig{ConcurrencyLimit: 2, ID: "pool2"},
	)

	for i := 0; i < 10; i++ {
		c.Add(newMyTask1(fmt.Sprintf("task1-%d", i)))
	}

	c.Start()
}
