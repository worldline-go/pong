package rest

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/registry"
)

// regex for space
var rgxSpace = regexp.MustCompile(`\s+`)

func Request(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, args *model.RestCheck, concurrent int, reg *registry.ClientReg) {
	// remove trailing spaces and multiple spaces
	urlX := strings.TrimSpace(args.URL)

	urlX = rgxSpace.ReplaceAllString(urlX, " ")

	urls := strings.Split(urlX, " ")
	timeout := time.Duration(args.Timeout) * time.Millisecond

	gData := GeneralData{}

	// create a new client holderFunc
	newClientHolderFn := NewClientHolder(gData)

	for i, url := range urls {
		selectedClient := i % concurrent
		// open new client if not exist
		reg.SetClient(ctx, wg, cancel, selectedClient, newClientHolderFn)

		if err := reg.SendMessage(ctx, &Msg{
			URL:     url,
			Method:  args.Method,
			Status:  args.Status,
			Timeout: timeout,
		}); err != nil {
			break
		}
	}
}
