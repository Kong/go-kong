package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//https://docs.konghq.com/enterprise/2.5.x/admin-api/event-hooks/reference/#test-an-event-hook
func TestEventHook(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	webhook := "webhook"
	crud := "crud"
	consumers := "consumers"
	cfg := &EventHooks{
		Config: map[string]interface{}{
			"url":        "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90",
			"ssl_verify": false,
			"secret":     " ",
		},
		Handler: &webhook,
		Source:  &crud,
		Event:   &consumers,
	}

	createdWebHook, err := client.EventHooks.AddWebhook(defaultCtx, cfg)
	if err != nil {
		panic(err)
	}
	if createdWebHook == nil {
		panic(err)
	}
}

func TestCustomWebHook(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	webhookCustom := "webhook-custom"
	crud := "crud"
	admins := "admins"
	cfg := &EventHooks{
		Config: map[string]interface{}{
			"url":            "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90",
			"body":           nil,
			"body_format":    nil,
			"method":         "POST",
			"payload_format": true,
			"secret":         nil,
			"ssl_verify":     false,
		},
		Handler: &webhookCustom,
		Source:  &crud,
		Event:   &admins,
	}

	createdCustomWebHook, err := client.EventHooks.AddWebhook(defaultCtx, cfg)
	if err != nil {
		panic(err)
	}
	if createdCustomWebHook == nil {
		panic(err)
	}
}
