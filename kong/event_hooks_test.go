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

	url := "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90"
	config := Config{
		URL: &url,
	}

	webhook := "webhook"
	crud := "crud"
	consumers := "consumers"
	cfg := &EventHooks{
		Config:  &config,
		Handler: &webhook,
		Source:  &crud,
		Event:   &consumers,
	}

	createdWebHook, err := client.EventHooks.AddWebhook(defaultCtx, cfg)
	assert.Nil(t, err)
	assert.NotNil(t, createdWebHook)
}

func TestCustomWebHook(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	url := "https://webhook.site/ec707ef0-ab91-4693-8dd2-114471ff6f90"
	method := "POST"
	config := Config{
		URL:    &url,
		Method: &method,
	}

	webhookCustom := "webhook-custom"
	crud := "crud"
	admins := "admins"
	cfg := &EventHooks{
		Config:  &config,
		Handler: &webhookCustom,
		Source:  &crud,
		Event:   &admins,
	}

	createdCustomWebHook, err := client.EventHooks.AddWebhook(defaultCtx, cfg)
	assert.Nil(t, err)
	assert.NotNil(t, createdCustomWebHook)
}
