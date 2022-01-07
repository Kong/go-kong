package kong

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInfoService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	info, err := client.Info.Get(defaultCtx)
	assert.Nil(err)
	assert.NotNil(info)
	assert.NotNil(info.Version)
	assert.NotNil(info.Configuration)
	assert.NotNil(info.Configuration.Database)
}

func TestConvert(T *testing.T) {
	assert := assert.New(T)
	information := map[string]interface{}{
		"version": "2.3.3.2-enterprise-edition",
		"configuration": map[string]interface{}{
			"portal":   true,
			"rbac":     "on",
			"database": "postgres",
		},
	}
	expected := &Info{
		Version: "2.3.3.2-enterprise-edition",
		Configuration: &RuntimeConfiguration{
			Portal:   true,
			RBAC:     "on",
			Database: "postgres",
		},
	}
	var actual Info
	err := convert(information, &actual)
	assert.NoError(err)
	assert.True(reflect.DeepEqual(expected, &actual))
	assert.False(actual.Configuration.IsInMemory())
	assert.True(actual.Configuration.IsRBACEnabled())
}

func TestIsConfigReady(T *testing.T) {
	for _, tt := range []struct {
		name      string
		status    int
		wantReady bool
		wantErr   bool
	}{
		{
			name:      "ready",
			status:    http.StatusOK,
			wantReady: true,
			wantErr:   false,
		},
		{
			name:      "not ready",
			status:    http.StatusServiceUnavailable,
			wantReady: false,
			wantErr:   false,
		},
		{
			name:      "unexpected status",
			status:    http.StatusNotFound,
			wantReady: false,
			wantErr:   true,
		},
	} {
		T.Run(tt.name, func(T *testing.T) {
			assert := assert.New(T)
			wait := &sync.WaitGroup{}
			server := mockConfigReady(T, tt.status, wait)
			mockURL := "http://localhost:32997"
			client, err := NewTestClient(&mockURL, nil)
			assert.NoError(err)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			ready, err := client.Info.IsConfigReady(ctx)
			assert.Equal(tt.wantReady, ready)
			if (err == nil) == tt.wantErr {
				T.Errorf("processCredential() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NoError(server.Shutdown(ctx))
			wait.Wait()
		})
	}
}

func mockConfigReady(T *testing.T, status int, wait *sync.WaitGroup) *http.Server {
	assert := assert.New(T)
	mux := http.NewServeMux()
	mux.Handle("/config/ready/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(status)
	}))
	server := &http.Server{Addr: ":32997", Handler: mux}
	go func() {
		defer wait.Done()
		wait.Add(1)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			assert.NoError(err)
		}
	}()
	return server
}
