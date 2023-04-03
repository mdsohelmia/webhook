package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/mdsohelmia/webhook/signer"
)

var (
	ErrWebhookUrlNotSet = errors.New("could not call the webhook because the URL has not been set")
	ErrSecretNotSet     = errors.New("could not call the webhook because no secret has been set")
	ErrPayloadNotSet    = errors.New("could not call the webhook because no payload has been set")
	ErrPayloadNotJson   = errors.New("could not call the webhook because the payload is not a valid JSON")
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
	RetryMax int
	//Signer signer.Signer
	Signer              signer.Signer
	SignWebhook         bool
	SignatureHeaderName string
	// Timeout for the request
	Timeout time.Duration
	//Debug bool
	//Default false
	Debug bool
	Url   string
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

func NewWebhook(config *Config) *Webhook {
	webhook := &Webhook{
		config:     config,
		signer:     config.Signer,
		httpMethod: MethodPost,
	}
	webhook.client = webhook.newClient()
	return webhook
}

//DefaultWebhook returns a new webhook with default configuration
//Default http method is POST
//Default signer is DefaultSigner
//Default retry wait min is 1 second
// Default retry wait max is 30 second
// Default retry max is 10
// Default timeout is 30 second

func DefaultWebhook() *Webhook {
	webhook := &Webhook{
		httpMethod: MethodPost,
		signer:     signer.NewDefaultSigner(),
		config: &Config{
			RetryWaitMin: 1 * time.Second,
			RetryWaitMax: 30 * time.Second,
			Timeout:      30 * time.Second,
			RetryMax:     10,
			SignWebhook:  false,
			Debug:        false,
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
	retryableClient.HTTPClient.Timeout = receiver.config.Timeout

	if !receiver.config.Debug {
		retryableClient.Logger = nil
	}
	// Return the underlying http.Client
	return retryableClient.StandardClient()
}

func (receiver *Webhook) SetUrl(url string) *Webhook {
	receiver.url = url
	return receiver
}

func (receiver *Webhook) Payload(payload interface{}) *Webhook {
	if payload == nil {
		receiver.err = ErrPayloadNotSet
		return receiver
	}
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

	if err := receiver.prepareRequest(); err != nil {
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

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return &Response{
		response: response,
		body:     body,
	}, nil
}

func (receiver *Webhook) prepareRequest() error {
	if receiver.url == "" {
		return ErrWebhookUrlNotSet
	}

	if receiver.config.SignWebhook && receiver.secret == "" {
		return ErrSecretNotSet
	}

	if receiver.payload == nil {
		return ErrPayloadNotSet
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

func (receiver *Webhook) SetDebug(debug bool) *Webhook {
	receiver.config.Debug = debug
	return receiver
}
