package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mdsohelmia/webhook/signer"
)

var (
	ErrWebhookUrlNotSet = errors.New("could not call the webhook because the URL has not been set")
	ErrSecretNotSet     = errors.New("could not call the webhook because no secret has been set")
)

type HttpMethod string

const (
	MethodPost HttpMethod = "POST"
	MethodPut  HttpMethod = "PUT"
	MethodGet  HttpMethod = "GET"
)

type Config struct {
	// Minimum time to wait
	RetryWaitMin time.Duration
	// Maximum time to wait
	RetryWaitMax time.Duration
	// Maximum number of retries
	RetryMax            int
	Signer              signer.Signer
	SignWebhook         bool
	SignatureHeaderName string
}
type Webhook struct {
	client     *http.Client
	config     *Config
	headers    map[string]string
	signer     signer.Signer
	payload    *bytes.Buffer
	err        error
	url        string
	secret     string
	httpMethod HttpMethod
}

func NewWithConfig(config *Config) *Webhook {
	webhook := &Webhook{
		config:     config,
		signer:     config.Signer,
		httpMethod: MethodPost,
	}
	webhook.client = webhook.newClient()
	return webhook
}

func New() *Webhook {
	webhook := &Webhook{
		httpMethod: MethodPost,
		signer:     signer.NewDefaultSigner(),
		config: &Config{
			RetryWaitMin: 1 * time.Second,
			RetryWaitMax: 30 * time.Second,
			RetryMax:     3,
			SignWebhook:  true,
		},
	}
	webhook.client = webhook.newClient()
	return webhook
}

func (receiver *Webhook) newClient() *http.Client {
	// Create a new RetryableClient
	retryableClient := retryablehttp.NewClient()
	// Configure the client to automatically retry failed requests with exponential backoff
	retryableClient.RetryWaitMin = receiver.config.RetryWaitMin
	retryableClient.RetryWaitMax = receiver.config.RetryWaitMax
	retryableClient.RetryMax = receiver.config.RetryMax

	return retryableClient.StandardClient()
}

func (receiver *Webhook) Url(url string) *Webhook {
	receiver.url = url
	return receiver
}
func (receiver *Webhook) Payload(payload interface{}) *Webhook {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		receiver.err = err
		return receiver
	}
	receiver.payload = bytes.NewBuffer(payloadBytes)

	return receiver
}

func (receiver *Webhook) UseSecret(secret string) *Webhook {
	receiver.secret = secret
	return receiver
}

func (receiver *Webhook) UseHttpVerb(method HttpMethod) *Webhook {
	receiver.httpMethod = method
	return receiver
}
func (receiver *Webhook) SignUsing(signer signer.Signer) *Webhook {
	receiver.signer = signer
	return receiver
}

func (receiver *Webhook) DoNotSign() *Webhook {
	receiver.config.SignWebhook = false
	return receiver
}

func (receiver *Webhook) UseSignatureHeaderName(name string) *Webhook {
	receiver.config.SignatureHeaderName = name
	return receiver
}

func (receiver *Webhook) WithHeaders(headers map[string]string) *Webhook {
	receiver.headers = headers
	return receiver
}

func (receiver *Webhook) getHeaders() map[string]string {
	headers := make(map[string]string)

	if receiver.config.SignatureHeaderName == "" {
		headers[receiver.signer.SignatureHeaderName()] = receiver.signer.CalculateSignature(receiver.url, receiver.payload.Bytes(), receiver.secret)
	}

	if receiver.config.SignatureHeaderName != "" {
		headers[receiver.config.SignatureHeaderName] = receiver.signer.CalculateSignature(receiver.url, receiver.payload.Bytes(), receiver.secret)
	}

	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["User-Agent"] = "Webhook"

	return headers
}

// Webhook Dispatch
func (receiver *Webhook) Dispatch() (*Response, error) {
	//If exist error
	if receiver.err != nil {
		return nil, receiver.err
	}

	if err := receiver.prepareForDispatch(); err != nil {
		return nil, err
	}

	request, err := receiver.makeRequest()

	if err != nil {
		return nil, err
	}

	response, err := receiver.client.Do(request)

	if err != nil {
		return nil, err
	}

	return &Response{
		response: response,
	}, nil
}

func (receiver *Webhook) prepareForDispatch() error {
	if receiver.url == "" {
		return ErrWebhookUrlNotSet
	}

	if receiver.config.SignWebhook && receiver.secret == "" {
		return ErrSecretNotSet
	}

	return nil
}
func (receiver *Webhook) makeRequest() (*http.Request, error) {
	request, err := http.NewRequest(string(receiver.httpMethod), receiver.url, receiver.payload)
	if err != nil {
		return nil, err
	}
	for key, value := range receiver.getHeaders() {
		request.Header.Set(key, value)
	}
	return request, nil
}
