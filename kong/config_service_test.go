package kong

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigService(t *testing.T) {
	RunWhenDBMode(t, "off")

	tests := []struct {
		name    string
		config  Configuration
		wantErr bool
	}{
		{
			name: "basic config works",
			config: Configuration{
				"_format_version": "1.1",
				"services": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
						"routes": []Configuration{
							{"paths": []string{"/"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing _format_version fails",
			config: Configuration{
				"services": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
						"routes": []Configuration{
							{"paths": []string{"/"}},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config fails",
			config: Configuration{
				"dummy_key": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		client, err := NewTestClient(nil, nil)
		require.NoError(t, err)
		require.NotNil(t, client)

		tt := tt
		t.Run("with_schema/"+tt.name, func(t *testing.T) {
			ctx := context.Background()
			b, err := json.Marshal(tt.config)
			require.NoError(t, err)

			body, err := client.Configs.ReloadDeclarativeRawConfig(ctx, bytes.NewBuffer(b), true)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SendConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			// this is somewhat untrue: network or HTTP-level failures _can_ result in a nil response body. however,
			// none of our test cases should cause network or HTTP-level failures, so fail if they do occur. if this
			// _does_ encounter such a failure, we need to investigate and either update tests or fix some upstream bug
			// if it's not some transient issue with the testing environment
			require.NotNilf(t, body, "body was nil; should never be nil")
		})
	}
}
