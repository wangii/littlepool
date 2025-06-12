package littlepool

import (
	"log"
	"time"
)

func newWorker[T Task](controller *Controller[T], pool *Pool, poolIdx int, busyChan chan int, idleChan chan int) {
	for {
		task := pool.takeTask()
		if task == nil {
			time.Sleep(time.Millisecond * 100)
			idleChan <- poolIdx
			continue
		}

		busyChan <- poolIdx
		if dt, ok := task.(DependentTask); ok && dt.HasPendingDependency() {
			log.Printf("task %s in pool %s has pending dependency, wait", task.ID(), pool.config.ID)
			pool.addTask(task)
			continue
		}

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
			ns := task.Next()

			if len(ns) > 0 {
				for _, n := range ns {
					p := controller.getPool(n.GetPoolID())
					p.addTask(n)
				}
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
