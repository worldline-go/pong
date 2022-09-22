package rest

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/registry"
	"github.com/worldline-go/pong/pkg/compare"
	"github.com/worldline-go/pong/pkg/template"
	"gopkg.in/yaml.v3"
)

type GeneralData struct{}

type Msg struct {
	URL     string
	Args    model.RestCheck
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

func NewClientHolder(gData GeneralData, cData model.RestSetting) func(ctx context.Context, ctxCancel context.CancelFunc, r *registry.ClientReg) registry.ClientHolder {
	return func(ctx context.Context, ctxCancel context.CancelFunc, r *registry.ClientReg) registry.ClientHolder {
		var client *http.Client

		if cData.InsecureSkipVerify {
			customTransport := http.DefaultTransport.(*http.Transport).Clone()
			customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			client = &http.Client{Transport: customTransport}
		} else {
			client = &http.Client{}
		}

		return &ClientHolder{
			Client:      client,
			MsgChan:     r.GetMsgChan(),
			GeneralData: gData,
			Reg:         r,
			Ctx:         ctx,
			CtxCancel:   ctxCancel,
		}
	}
}

func (c *ClientHolder) DoRequest(ctx context.Context, timeout time.Duration, urlV string, m model.RestCheck) error {
	method := cleanMethod(m.Method)

	ctxT := ctx
	if timeout != 0 {
		var cancel context.CancelFunc
		ctxT, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctxT, method, urlV, nil)
	if err != nil {
		return fmt.Errorf("%s, creating request: %w", urlV, err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%s, doing request: %w", urlV, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Info().Msgf("%s", body)

	if resp.StatusCode != m.Status {
		return fmt.Errorf("%s, status code: %d; want: %d", urlV, resp.StatusCode, m.Status)
	}

	if m.Body != nil {
		if m.Body.Map != nil {
			var bodyMap interface{}
			if err := yaml.Unmarshal(body, &bodyMap); err != nil {
				return fmt.Errorf("%s, unmarshaling body: %w", body, err)
			}

			mapValues := m.Body.Variable.Set
			if m.Body.Variable.Set == nil {
				mapValues = make(map[string]interface{})
			}

			urlP, err := url.Parse(urlV)
			if err != nil {
				return fmt.Errorf("%s, parsing url: %w", urlV, err)
			}

			urlValues := urlP.Query()
			for _, queryValue := range m.Body.Variable.From.Query {
				mapValues[queryValue] = urlValues.Get(queryValue)
			}

			rendered, err := template.Ext(mapValues, *m.Body.Map)
			if err != nil {
				return fmt.Errorf("%s, rendering template: %w", urlV, err)
			}

			// check body
			var checkBody interface{}
			if err := yaml.Unmarshal([]byte(rendered), &checkBody); err != nil {
				return fmt.Errorf("%s, unmarshaling body: %w", rendered, err)
			}

			if err := compare.IsSubset(bodyMap, checkBody); err != nil {
				return fmt.Errorf("%s, comparing body: %w", urlV, err)
			}
		}

		if m.Body.Raw != nil {
			if *m.Body.Raw != string(body) {
				return fmt.Errorf("%s, comparing body: %s; want: %s", urlV, body, *m.Body.Raw)
			}
		}
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

			if err := c.DoRequest(c.Ctx, m.Timeout, m.URL, m.Args); err != nil {
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
