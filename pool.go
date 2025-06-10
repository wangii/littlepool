package littlepool

import (
	"sync"
)

type PoolConfig struct {
	ID               string
	ConcurrencyLimit int
}

type Pool struct {
	config PoolConfig
	mu     sync.Mutex

	tasks []Task
}

func NewPool(config PoolConfig) *Pool {
	return &Pool{
		config: config,
	}
}

func (p *Pool) addTask(task Task) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.tasks = append(p.tasks, task)
}

func (p *Pool) takeTask() Task {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.tasks) > 0 {
		task := p.tasks[0]
		p.tasks = p.tasks[1:]
		return task
	}
	return nil
}

func (p *Pool) hasMore() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.tasks) > 0
}
