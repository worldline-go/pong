package route

import (
	"context"
	"sync"

	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/registry"
	"github.com/worldline-go/pong/internal/route/rest"
)

func Request(ctx context.Context, args *model.ModuleArgs) []error {
	wg := &sync.WaitGroup{}

	errs := &registry.Errors{}

	// call rest requests
	wg.Add(1)
	go RestRequest(ctx, wg, errs, args.Client.Rest)

	wg.Wait()

	return errs.Errs
}

func RestRequest(ctx context.Context, wg *sync.WaitGroup, errs *registry.Errors, args []model.RestClient) {
	defer wg.Done()

	wgRest := &sync.WaitGroup{}
	for _, client := range args {
		// new area for each client
		reg := registry.NewClientReg(errs)
		wgRest.Add(1)
		// fast pass
		go RestCheck(ctx, wgRest, client.Check, client.Concurrent, reg)
	}

	wgRest.Wait()
}

func RestCheck(ctx context.Context, wg *sync.WaitGroup, args []model.RestCheck, concurrent int, reg *registry.ClientReg) {
	defer wg.Done()
	wgRest := &sync.WaitGroup{}

	ctxRest, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range args {
		rest.Request(ctxRest, cancel, wgRest, &args[i], concurrent, reg)
	}

	// all messages sent close channel for registry
	reg.CloseChan()
	wgRest.Wait()
}
