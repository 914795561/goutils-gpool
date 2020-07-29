package gpool

import (
	"log"
	"runtime"
)

type F func()

type Pool struct {
	work  chan F
	queue chan struct{}
	PanicHandler func(...interface{})
}

func NewPool(cap int) *Pool {
	return &Pool{
		work:  make(chan F),
		queue: make(chan struct{}, cap),
	}
}

func (p *Pool) Add(task F) {
	select {
	case p.work <- task:
	case p.queue <- struct{}{}:
		go p.worker(task)
	}
}

func (p *Pool) worker(task F) {
	defer func() {
		if r := recover(); r != nil {
			pc, file, line, ok := runtime.Caller(3)
			if ok{
				if p.PanicHandler != nil {
					p.PanicHandler(r)
				} else {
					funcName := runtime.FuncForPC(pc).Name()
					log.Println("func", funcName, "file", file, "line", line, " ", r)
				}
			}
		}
		<-p.queue
	}()

	for {
		task()
		task = <-p.work
	}
}
