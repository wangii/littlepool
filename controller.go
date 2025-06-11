package littlepool

import "sync"

type kWorkerStatus int

const (
	kWorkerStatusIdle kWorkerStatus = iota
	kWorkerStatusBusy
)

type Controller[T Task] struct {
	pools    []*Pool
	status   []kWorkerStatus
	muFinish sync.Mutex
	finished []T
}

func NewController[T Task](cfgs ...PoolConfig) *Controller[T] {
	ret := &Controller[T]{}

	nWorker := 0
	for _, cfg := range cfgs {
		ret.pools = append(ret.pools, NewPool(cfg))
		nWorker += cfg.ConcurrencyLimit
	}
	ret.status = make([]kWorkerStatus, nWorker)

	for idx := range nWorker {
		ret.status[idx] = kWorkerStatusBusy
	}

	return ret
}

func (c *Controller[T]) getPool(id string) *Pool {
	for _, p := range c.pools {
		if p.config.ID == id {
			return p
		}
	}
	return nil
}

func (c *Controller[T]) Add(task Task) {
	pool := c.getPool(task.GetPoolID())
	if pool == nil {
		panic("No pool found for task")
	}
	pool.addTask(task)
}

func (c *Controller[T]) Start() {
	busyChan := make(chan int)
	idleChan := make(chan int)

	for idx, p := range c.pools {
		for range p.config.ConcurrencyLimit {
			go newWorker(c, p, idx, busyChan, idleChan)
		}
	}

main:
	for {
		select {
		case idx := <-idleChan:
			{
				c.status[idx] = kWorkerStatusIdle
				for i := range len(c.status) {
					if c.status[i] == kWorkerStatusBusy {
						continue main
					}
				}

				for _, p := range c.pools {
					if p.hasMore() {
						continue main
					}
				}

				break main
			}
		case idx := <-busyChan:
			{
				c.status[idx] = kWorkerStatusBusy
			}
		}
	}
}

func (c *Controller[T]) appendFinished(task T) {
	c.muFinish.Lock()
	defer c.muFinish.Unlock()

	c.finished = append(c.finished, task)
}

func (c *Controller[T]) IterateFinished(f func(task T)) {
	for _, task := range c.finished {
		f(task)
	}
}

func (c *Controller[T]) GetFinished() []T {
	return c.finished
}
