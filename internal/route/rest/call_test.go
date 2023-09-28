package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/worldline-go/pong/internal/model"
	"github.com/worldline-go/pong/internal/registry"
)

func TestRequest(t *testing.T) {
	type args struct {
		check *model.RestCheck
	}
	tests := []struct {
		name       string
		args       args
		concurrent int
		want       string
		wantIn     bool
		handler    func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "simple one test",
			args: args{
				check: &model.RestCheck{
					Respond: model.RestRespond{
						Status: 200,
					},
					Request: model.RestRequest{
						Timeout: 2,
					},
				},
			},
			concurrent: 1,
			want:       "[]",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name: "multi test",
			args: args{
				check: &model.RestCheck{
					Request: model.RestRequest{
						Timeout: 2,
						Method:  "GET",
						URL:     "/abc /xyz /def /a /b /c /d /e /f /g /h /i /j /k /l /m /n /o /p /q /r /s /t /u /v /w /x /y /z",
					},
					Respond: model.RestRespond{
						Status: 200,
					},
				},
			},
			concurrent: 1,
			want:       `status code: 502; want: 200`,
			wantIn:     true,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			},
		},
		{
			name: "multi test with concurrent",
			args: args{
				check: &model.RestCheck{
					Request: model.RestRequest{
						URL:     "/abc /xyz /def /a /b /c /d /e /f /g /h /i /j /k /l /m /n /o /p /q /r /s /t /u /v /w /x /y /z",
						Method:  "GET",
						Timeout: 2,
					},
					Respond: model.RestRespond{
						Status: 200,
					},
				},
			},
			concurrent: 20,
			want:       `status code: 502; want: 200`,
			wantIn:     true,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			},
		},
		{
			name: "test header",
			args: args{
				check: &model.RestCheck{
					Request: model.RestRequest{
						URL:     "/abc /xyz",
						Method:  "GET",
						Timeout: 2,
						Headers: map[string]string{
							"X-Test":  "test",
							"X-Test2": "test2",
						},
					},
					Respond: model.RestRespond{
						Status: 200,
					},
				},
			},
			concurrent: 2,
			want:       `[]`,
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Test") != "test" && r.Header.Get("X-Test2") != "test2" {
					w.WriteHeader(http.StatusBadGateway)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		},
	}

	var handlerFunc = func(w http.ResponseWriter, r *http.Request) {}
	mx := sync.RWMutex{}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mx.RLock()
		defer mx.RUnlock()

		if handlerFunc != nil {
			handlerFunc(w, r)
		}
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			errs := &registry.Errors{}
			reg := registry.NewClientReg(errs, model.RestSetting{}, false)

			mx.Lock()
			handlerFunc = tt.handler
			mx.Unlock()

			urls := strings.Split(tt.args.check.Request.URL, " ")
			for i, url := range urls {
				urls[i] = fmt.Sprintf("%s%s", srv.URL, url)
			}

			tt.args.check.Request.URL = strings.Join(urls, " ")
			ctx, cancel := context.WithCancel(context.Background())
			Request(ctx, cancel, wg, tt.args.check, tt.concurrent, reg)
			reg.CloseChan()
			wg.Wait()
			cancel()

			if !tt.wantIn && fmt.Sprintf("%v", reg.Errors.Errs) != tt.want {
				t.Errorf("Request() = %v, want %v", reg.Errors.Errs, tt.want)
			}

			if tt.wantIn && !strings.Contains(fmt.Sprintf("%v", reg.Errors.Errs), tt.want) {
				t.Errorf("Request() = %v, want %v", reg.Errors, tt.want)
			}
		})
	}
}
