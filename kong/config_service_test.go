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

			if err, _ := client.Configs.ReloadDeclarativeRawConfig(ctx, bytes.NewBuffer(b), true); (err != nil) != tt.wantErr {
				t.Errorf("Client.SendConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
