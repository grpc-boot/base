package gopool

import (
	"errors"
	"time"
)

var (
	ErrSubmitTimeout    = errors.New("submit error: timed out")
	ErrOptionsDeadQueue = errors.New("options error:dead queue configuration")
	ErrOptionsSpawn     = errors.New("options error:spawn > workers")
)

// Pool goroutine pool
type Pool struct {
	sem          chan struct{}
	work         chan func()
	panicHandler func(err interface{})
}

func NewPool(size int, opts ...Option) (*Pool, error) {
	options := loadOptions(opts...)
	options.size = size

	if options.spawn <= 0 && options.queue > 0 {
		return nil, ErrOptionsDeadQueue
	}

	if options.spawn > options.size {
		return nil, ErrOptionsSpawn
	}

	p := &Pool{
		sem:  make(chan struct{}, options.size),
		work: make(chan func(), options.queue),
	}

	for i := 0; i < options.spawn; i++ {
		p.sem <- struct{}{}
		go p.worker(func() {})
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
func (p *Pool) ActiveWorkerNum() int {
	return len(p.sem)
}

// QueueLength get queue item number
func (p *Pool) QueueLength() int {
	return len(p.work)
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrSubmitTimeout
	case p.work <- task:
		return nil
	case p.sem <- struct{}{}:
		go p.worker(task)
		return nil
	}
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.sem
	}()

	p.runTask(task)

	for t := range p.work {
		p.runTask(t)
	}
}

func (p *Pool) runTask(task func()) {
	defer func() {
		if err := recover(); err != nil {
			p.panicHandler(err)
		}
	}()
	task()
}
