package signer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureHeaderName(t *testing.T) {
	expected := "X-Signature"
	actual := NewDefaultSigner().SignatureHeaderName()
	assert.Equal(t, expected, actual)
}

func TestCalculateSignature(t *testing.T) {
	url := "http://webhook.com/webhook"
	payload := []byte("Hello Gotipath")
	secret := "gotipath ErrSecretNotSet"

	expected := NewDefaultSigner().CalculateSignature(url, payload, secret)

	assert.Equal(t, expected, NewDefaultSigner().CalculateSignature(url, payload, secret))
}
