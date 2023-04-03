package main

import (
	"github.com/k0kubun/pp/v3"
	"github.com/mdsohelmia/webhook"
)

func main() {
	webhook := webhook.DefaultWebhook()

	response, err := webhook.
		SetUrl("https://api.gotipath.test/v1/webhook").
		SetDebug(true).
		Payload("hello playlod").
		Dispatch()

	if err != nil {
		pp.Println(err)
		return
	}

	pp.Println("body", string(response.Body()))
	pp.Println(response.StatusCode())
	pp.Println(response.Headers())
	pp.Println(response.Status())
	pp.Println(response.Ok())
}
