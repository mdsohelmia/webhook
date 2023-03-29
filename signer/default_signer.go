package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// DefaultSigner
type DefaultSigner struct{}

var _ Signer = (*DefaultSigner)(nil)

func NewDefaultSigner() *DefaultSigner {
	return &DefaultSigner{}
}

func (ds *DefaultSigner) CalculateSignature(webhookUrl string, payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ds *DefaultSigner) SignatureHeaderName() string {
	return "X-Signature"
}
