package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/pong/internal/registry"
)

type GeneralData struct{}

type Msg struct {
	URL     string
	Method  string
	Status  int
	Timeout time.Duration
}

type ClientHolder struct {
	Ctx         context.Context
	Client      *http.Client
	MsgChan     <-chan interface{}
	Reg         *registry.ClientReg
	CtxCancel   context.CancelFunc
	GeneralData GeneralData
}

func NewClientHolder(gData GeneralData) func(ctx context.Context, ctxCancel context.CancelFunc, r *registry.ClientReg) registry.ClientHolder {
	return func(ctx context.Context, ctxCancel context.CancelFunc, r *registry.ClientReg) registry.ClientHolder {
		return &ClientHolder{
			Client:      &http.Client{},
			MsgChan:     r.GetMsgChan(),
			GeneralData: gData,
			Reg:         r,
			Ctx:         ctx,
			CtxCancel:   ctxCancel,
		}
	}
}

func (c *ClientHolder) DoRequest(ctx context.Context, timeout time.Duration, method, url string, status int) error {
	method = cleanMethod(method)

	ctxT := ctx
	if timeout != 0 {
		var cancel context.CancelFunc
		ctxT, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctxT, method, url, nil)
	if err != nil {
		return fmt.Errorf("%s, creating request: %w", url, err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%s, doing request: %w", url, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Info().Msgf("%s", body)

	if resp.StatusCode != status {
		return fmt.Errorf("%s, status code: %d; want: %d", url, resp.StatusCode, status)
	}

	return nil
}

func (c *ClientHolder) Work() {
	for {
		select {
		case <-c.Ctx.Done():
			return
		case msg := <-c.MsgChan:
			// check channel is closed
			if msg == nil {
				return
			}

			// check message type
			m, ok := msg.(*Msg)
			if !ok {
				log.Error().Msgf("wrong message type: %T", msg)
				continue
			}

			log.Debug().Msgf("Sending request to %s", m.URL)

			if err := c.DoRequest(c.Ctx, m.Timeout, m.Method, m.URL, m.Status); err != nil {
				c.close()
				// record error
				c.Reg.AddError(err)

				return
			}
		}
	}
}

func (c *ClientHolder) close() {
	// stop to redirect new messages
	c.CtxCancel()
}
