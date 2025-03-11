package server

import (
	"context"
	"sync"
)

type Notifier struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewNotifier() *Notifier {
	ctx, cancel := context.WithCancel(context.Background())
	return &Notifier{
		wg:     &sync.WaitGroup{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (n *Notifier) Register() (doneChan <-chan struct{}, finish func()) {
	n.wg.Add(1)
	return n.ctx.Done(), func() {
		n.wg.Done()
	}
}

func (n *Notifier) NotifyAndWait() {
	n.cancel()
	n.wg.Wait()
}
