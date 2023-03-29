package webhook

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

func TestWebhookTest(t *testing.T) {

	res, err := New().
		Url("https://api.gotipath.test/v1/webhook").
		UseHttpVerb(MethodPut).
		Payload(map[string]string{
			"id":      "100",
			"event":   "customer.created",
			"event1":  "customer.created",
			"event2":  "customer.created",
			"event3":  "customer.created",
			"event4":  "customer.created",
			"event5":  "customer.created",
			"event6":  "customer.created",
			"event7":  "customer.created",
			"event8":  "customer.created",
			"event9":  "customer.created",
			"event10": "customer.created",
		}).
		UseSignatureHeaderName("X-Gotipath-Signature").
		UseSecret("2Y4Zeqi2foh70p8xiDbpDjGP8kjttCmmSIXLi7EUVeREwapHaZH1gSDuQLicbyzbOUP2kwidQlhZ40moXTY2mB9iKYtkRfjrWsL5h8RQUMs4eQjK2sL97QsXKtw6Qexogtl5cVVeVyElk1JgxjalI0gUWkivjHXl1haxdkHoTCFkuH3J2sRlRGDSlhYEkETeC7eU3TOxSAtwY1vlm2AvlyaUy8rlCoccO0tZOKiOoHWWkx6UTIzSIKIgX86hV4WCUKyKUmvGGhuKL9ci3puLtLhzS1sEbC5RpsShYBPXYHGBL2BGD38BpzABve1r1Kg8k8AdGSfJ3yBSGGo4Crj0HGyF1OX6ij00K4qtzID8A1iOisklgI3IwdPt4yplCF3HCDJOrxAQexKj3kJZJ1x54JA4MU5hI75A6zIpCskp0oraKsg70BOxZ8wh").
		Dispatch()

	if err != nil {
		pp.Println(err.Error())
	}

	fmt.Println(res.response.Status)

}
