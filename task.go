package littlepool

type TaskResult int

const (
	TaskResultSuccess TaskResult = iota
	TaskResultFailedAbort
	TaskResultFailedRetry
)

type Task interface {
	ID() string
	GetPoolID() string

	Run() TaskResult
	Next() []Task
}

type DependentTask interface {
	HasPendingDependency() bool
}
