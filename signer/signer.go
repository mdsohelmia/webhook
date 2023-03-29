package signer

type Signer interface {
	CalculateSignature(webhookUrl string, payload []byte, secret string) string
	SignatureHeaderName() string
}
