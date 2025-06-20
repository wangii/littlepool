package littlepool

import (
	"log"
	"time"
)

func newWorker[T Task](controller *Controller[T], pool *Pool, workerIdx int, busyChan chan int, idleChan chan int) {
	for {
		task := pool.takeTask()
		if task == nil {
			time.Sleep(time.Millisecond * 100)
			idleChan <- workerIdx
			continue
		}

		busyChan <- workerIdx
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
				f, ok := task.(T)
				if ok {
					controller.appendFinished(f)
				} else {
					log.Printf("Not able to add %s to finished.", task.ID())
				}
			}
		}
	}
}
