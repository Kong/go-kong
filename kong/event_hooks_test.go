package kong

import (
	"fmt"
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
			"url":        "https://webhook.site/a1b2c3-d4e5-g6h7-i8j9-k1l2m3n4o5p6",
			"ssl_verify": false,
			"secret":     " ",
		},
		Handler: &webhook,
		Source:  &crud,
		Event:   &consumers,
	}

	createdWebHook, err := client.EventHook.AddWebhook(defaultCtx, cfg)
	if err != nil {
		panic(err)
	}
	if createdWebHook == nil {
		panic(err)
	}
}

func TestCustomWebHook(t *testing.T) {
	webhook_custom := "webhook-custom"
	crud := "crud"
	admins := "admins"
	cfg := &EventHooks{
		Config: map[string]interface{}{
			"url":        "https://webhook.site/a1b2c3-d4e5-g6h7-i8j9-k1l2m3n4o5p6",
			"body":       nil,
			"body_format": nil,
			"headers":{
				"content-type": "application/json",
			},
			"method": "POST",
        "payload": {
            "text": "Admin account `` d; email address set to ``"
        },
        "payload_format": true,
        "secret": null,
        "ssl_verify": false,
		},
		Handler: &webhook_custom,
		Source:  &crud,
		Event:   &admins,
	}

	createdCustomWebHook, err := client.EventHook.AddWebhook(defaultCtx, cfg)
	if err != nil {
		panic(err)
	}
	if createdCustomWebHook == nil {
		panic(err)
	}
}
