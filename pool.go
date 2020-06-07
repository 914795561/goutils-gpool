package gpool

type F func()

type Pool struct {
	work  chan F
	queue chan struct{}
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
		<-p.queue
	}()

	for {
		task()
		task = <-p.work
	}
}
