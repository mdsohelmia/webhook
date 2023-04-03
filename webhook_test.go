package webhook

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp/v3"
)

func TestWebhook(t *testing.T) {
	webhook := DefaultWebhook()
	pp.Println(webhook)
	res, err := webhook.
		SetUrl("https://api.gotipath.test/v1/webhook").
		SetDebug(true).
		Payload("hello playlod").
		UseHttpVerb(MethodPost).
		Dispatch()

	if err != nil {
		pp.Println(err)
		return
	}

	fmt.Println("status:", res.Status())

}

// 240p=  181 kb/s = 25 fps
// 360p = 370 kb/s = 25 fps
// 480p = 1000 kb/s = 25 fps
// 540p.mp4 = 719 kb/s = 25 fps
// 720p =  1033 kb/s = 25 fps
// 1080p = 2133 kb/s = 25 fps
// 1440p = 6000 kb/s = 25 fps
// 2160p = 10000 kb/s = 25 fps
