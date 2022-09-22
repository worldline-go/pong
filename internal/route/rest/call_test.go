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
					Status:  200,
					Timeout: 2,
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
					URL:     "/abc /xyz",
					Method:  "GET",
					Status:  200,
					Timeout: 2,
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
					URL:     "/abc /xyz",
					Method:  "GET",
					Status:  200,
					Timeout: 2,
				},
			},
			concurrent: 2,
			want:       `status code: 502; want: 200`,
			wantIn:     true,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			},
		},
	}

	var handlerFunc = func(w http.ResponseWriter, r *http.Request) {}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handlerFunc != nil {
			handlerFunc(w, r)
		}
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := &sync.WaitGroup{}
			errs := &registry.Errors{}
			reg := registry.NewClientReg(errs, nil)

			handlerFunc = tt.handler
			urls := strings.Split(tt.args.check.URL, " ")
			for i, url := range urls {
				urls[i] = fmt.Sprintf("%s%s", srv.URL, url)
			}

			tt.args.check.URL = strings.Join(urls, " ")
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
