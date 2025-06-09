package littlepool

import (
	"log"
	"time"
)

func newWorker(controller *Controller, pool *Pool, poolIdx int, busyChan chan int, idleChan chan int) {
	for {
		task := pool.takeTask()
		if task == nil {
			idleChan <- poolIdx
			time.Sleep(time.Millisecond * 100)
			continue
		}

		busyChan <- poolIdx
		log.Printf("start task %s, from pool %s", task.ID(), pool.config.ID)

		ret := task.Run()

		if TaskResultFailedAbort == ret {
			log.Printf("task %s in pool %s failed, abort", task.ID(), pool.config.ID)
		}

		if TaskResultFailedRetry == ret {
			log.Printf("task %s in pool %s failed, retry", task.ID(), pool.config.ID)
			pool.addTask(task)
		}

		if TaskResultSuccess == ret {
			log.Printf("task %s in pool %s success", task.ID(), pool.config.ID)
			n := task.Next()

			if n != nil {
				p := controller.getPool(n.GetPoolID())
				p.addTask(n)
			} else {
				controller.finished = append(controller.finished, task)
			}
		}
	}
}
