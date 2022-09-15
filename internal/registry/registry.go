package registry

import (
	"context"
	"fmt"
	"sync"
)

type ClientReg struct {
	msgChan chan interface{}
	clients []interface{}
	*Errors
	// mutexErr sync.RWMutex
}

func NewClientReg(errs *Errors) *ClientReg {
	msgChan := make(chan interface{})
	return &ClientReg{
		msgChan: msgChan,
		Errors:  errs,
	}
}

type ClientHolder interface {
	Work()
}

type GetNewClientHolder func(ctx context.Context, ctxCancel context.CancelFunc, r *ClientReg) ClientHolder

func (r *ClientReg) GetMsgChan() <-chan interface{} {
	return r.msgChan
}

func (r *ClientReg) CloseChan() {
	close(r.msgChan)
}

func (r *ClientReg) SetClient(ctx context.Context, wg *sync.WaitGroup, cancel context.CancelFunc,
	i int, getNewClientHolder GetNewClientHolder,
) *ClientReg {
	if i < 0 {
		i = 0
	}

	if i+1 > len(r.clients) {
		instance := getNewClientHolder(ctx, cancel, r)

		r.clients = append(r.clients, instance)

		wg.Add(1)

		go func(client ClientHolder) {
			defer wg.Done()

			client.Work()
		}(instance)
	}

	return r
}

func (r *ClientReg) SendMessage(ctx context.Context, msg interface{}) error {
	select {
	case r.msgChan <- msg:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context canceled")
	}
}
