package signer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultSigner(t *testing.T) {
	expected := &DefaultSigner{}

	assert.Equal(t, expected, NewDefaultSigner())
}

func TestSignatureHeaderName(t *testing.T) {
	expected := "X-Signature"
	actual := NewDefaultSigner().SignatureHeaderName()
	assert.Equal(t, expected, actual)
}

func TestCalculateSignature(t *testing.T) {
	url := "http://webhook.com/webhook"
	payload := []byte("Hello Gotipath")
	secret := "gotipath ErrSecretNotSet"

	expected := "c569541beaaaed1f917fe993635cc8c027d29d2ff73a6f5f7c28843e24a3aead"

	assert.Equal(t, expected, NewDefaultSigner().CalculateSignature(url, payload, secret))
}
