package rest

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/registry"
)

func Request(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, args *model.RestCheck, concurrent int, reg *registry.ClientReg) {
	// remove trailing spaces and multiple spaces
	urlX := strings.TrimSpace(args.Request.URL)
	urls := strings.Fields(urlX)

	var timeout time.Duration
	if args.Request.Timeout != 0 {
		timeout = time.Duration(args.Request.Timeout) * time.Millisecond
	}

	gData := GeneralData{}
	cData := reg.ClientData

	// create a new client holderFunc
	newClientHolderFn := NewClientHolder(gData, cData)

	for i, url := range urls {
		selectedClient := i % concurrent
		// open new client if not exist
		reg.SetClient(ctx, wg, cancel, selectedClient, newClientHolderFn)

		if err := reg.SendMessage(ctx, &Msg{
			URL:     url,
			Args:    *args,
			Timeout: timeout,
		}); err != nil {
			break
		}
	}
}
