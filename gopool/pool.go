package gopool

import (
	"errors"
	"time"

	"go.uber.org/atomic"
)

var (
	ErrSubmitTimeout    = errors.New("submit error: timed out")
	ErrOptionsDeadQueue = errors.New("options error:dead queue configuration")
	ErrOptionsSpawn     = errors.New("options error:spawn > workers")
)

// Pool goroutine pool
type Pool struct {
	sem            chan struct{}
	work           chan func()
	maxIdleTimeout time.Duration
	panicHandler   func(err any)
	pendingTaskNum atomic.Int64
	successTotal   atomic.Uint64
	failedTotal    atomic.Uint64
}

func NewPool(maxWorkerNum int, opts ...Option) (*Pool, error) {
	options := loadOptions(opts...)
	options.size = maxWorkerNum

	if options.spawn <= 0 && options.queue > 0 {
		return nil, ErrOptionsDeadQueue
	}

	if options.spawn > options.size {
		return nil, ErrOptionsSpawn
	}

	p := &Pool{
		sem:            make(chan struct{}, options.size),
		work:           make(chan func(), options.queue),
		panicHandler:   options.panicHandler,
		maxIdleTimeout: time.Second * time.Duration(options.maxIdleTimeoutSeconds),
	}

	for i := 0; i < options.spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(nil)
	}

	return p, nil
}

// Submit _
func (p *Pool) Submit(task func()) error {
	return p.schedule(task, nil)
}

// SubmitTimeout _
func (p *Pool) SubmitTimeout(timeout time.Duration, task func()) error {
	return p.schedule(task, time.After(timeout))
}

// ActiveWorkerNum get active worker number
func (p *Pool) ActiveWorkerNum() int64 {
	return int64(len(p.sem))
}

// QueueLength get queue item number
func (p *Pool) QueueLength() int64 {
	return int64(len(p.work))
}

// PendingTaskTotal get pending task num
func (p *Pool) PendingTaskTotal() int64 {
	return p.pendingTaskNum.Load() + p.QueueLength()
}

// SuccessTotal _
func (p *Pool) SuccessTotal() uint64 {
	return p.successTotal.Load()
}

// FailedTotal _
func (p *Pool) FailedTotal() uint64 {
	return p.failedTotal.Load()
}

// HandleTotal _
func (p *Pool) HandleTotal() uint64 {
	return p.SuccessTotal() + p.FailedTotal()
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrSubmitTimeout
	case p.work <- task:
		p.pendingTaskNum.Inc()
		return nil
	case p.sem <- struct{}{}:
		p.pendingTaskNum.Inc()
		go p.worker(task)
		return nil
	}
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()

	p.runTask(task)

	if p.maxIdleTimeout < time.Second {
		for {
			select {
			case t, ok := <-p.work:
				p.runTask(t)
				if !ok {
					return
				}
			default:
			}
		}
	}

	var (
		ticker  = time.NewTicker(p.maxIdleTimeout)
		working = false
	)

	defer func() {
		ticker.Stop()
		ticker = nil
	}()

	for {
		select {
		case t, ok := <-p.work:
			working = true
			p.runTask(t)
			if !ok {
				return
			}
		case <-ticker.C:
			if !working {
				return
			}
			working = false
		default:
		}
	}
}

func (p *Pool) runTask(task func()) {
	if task == nil {
		return
	}

	defer func() {
		p.pendingTaskNum.Dec()
		if err := recover(); err != nil {
			p.failedTotal.Inc()
			p.panicHandler(err)
		} else {
			p.successTotal.Inc()
		}
	}()

	task()
}
