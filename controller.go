package littlepool

type kPoolStatus int

const (
	kPoolStatusIdle kPoolStatus = iota
	kPoolStatusBusy
)

type Controller struct {
	pools  []*Pool
	status []kPoolStatus
}

func NewController(cfgs ...PoolConfig) *Controller {
	ret := &Controller{}
	for _, cfg := range cfgs {
		ret.pools = append(ret.pools, NewPool(cfg))
	}
	ret.status = make([]kPoolStatus, len(cfgs))

	for idx := range ret.status {
		ret.status[idx] = kPoolStatusBusy
	}

	return ret
}

func (c *Controller) getPool(id string) *Pool {
	for _, p := range c.pools {
		if p.config.ID == id {
			return p
		}
	}
	return nil
}

func (c *Controller) Add(task Task) {
	pool := c.getPool(task.GetPoolID())
	if pool == nil {
		panic("No pool found for task")
	}
	pool.addTask(task)
}

func (c *Controller) Start() {
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
				c.status[idx] = kPoolStatusIdle
				for i := range len(c.pools) {
					if c.status[i] == kPoolStatusBusy {
						continue main
					}
				}
				break main
			}
		case idx := <-busyChan:
			{
				c.status[idx] = kPoolStatusBusy
			}
		}
	}
}
