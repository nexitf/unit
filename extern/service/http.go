package service

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/nexitf/unit/internal/errors"
	"github.com/nexitf/unit/plugin"
)

var (
	ErrEndpointNotFound = errors.New("endpoint not found")
)

// WithHTTPHeader binds the http client with the header.
//
// Support:
//   - HTTPClient
func WithHTTPHeader(key, value string, replace bool) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			if replace {
				client.headers.Set(key, value)
			} else {
				client.headers.Add(key, value)
			}
		}
		return
	}
}

// WithHTTPHost binds the http client with the host.
//
// Support:
//   - HTTPClient
func WithHTTPHost(host string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.host = host
		}
		return
	}
}

// WithHTTPS binds the http client with the https.
//
// Support:
//   - HTTPClient
func WithHTTPS() plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.https = true
		}
		return
	}
}

// WithHTTPTimeout binds the http client with the timeout.
//
// Support:
//   - HTTPClient
func WithHTTPTimeout(timeout time.Duration) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.timeout = timeout
		}
		return
	}
}

// WithHTTPRetryCount binds the http client with the count.
//
// Support:
//   - HTTPClient
func WithHTTPRetryCount(count int) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.retry = count
		}
		return
	}
}

// WithHTTPProxy binds the http client with the url.
//
// Support:
//   - HTTPClient
func WithHTTPProxy(url string) plugin.BindOption {
	return func(varp plugin.Resource) (used bool) {
		client, used := varp.(*HTTPClient)
		if used {
			client.proxy = url
		}
		return
	}
}

type HTTPClient struct {
	Resource
	balancer
	client  *httpclient.Client
	retry   int
	timeout time.Duration
	proxy   string
	host    string
	https   bool
	headers http.Header
}

// init
func (client *HTTPClient) init() {
	client.headers = make(http.Header)
}

// bind
func (client *HTTPClient) bind(opts ...plugin.BindOption) (unused []plugin.BindOption) {
	for _, setOpt := range opts {
		if !setOpt(client) {
			unused = append(unused, setOpt)
		}
	}

	var (
		newOpts []httpclient.Option
	)

	// Option: timeout
	if client.timeout <= 0 {
		client.timeout = 3 * time.Second
	}
	newOpts = append(newOpts, httpclient.WithHTTPTimeout(client.timeout))
	// Option: retry
	if client.retry <= 0 {
		client.retry = 0
	} else {
		newOpts = append(newOpts, httpclient.WithRetryCount(client.retry))
	}
	// Option: proxy
	if client.proxy != "" {
		newOpts = append(newOpts, httpclient.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) { return url.Parse(client.proxy) },
			},
			Timeout: client.timeout,
		}))
	}
	// Option: balancer
	if client.balancer == nil {
		client.balancer = NewRoundRobinBalancer()
	}

	client.client = httpclient.NewClient(newOpts...)
	return
}

// makeHeaders
func (client *HTTPClient) makeHeaders(h http.Header) (headers http.Header) {
	headers = make(http.Header)
	// Option: headers
	for key, values := range client.headers {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	// Copy custom header
	for key, values := range h {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	return
}

// makeURL
func (client *HTTPClient) makeURL(uri string) (url string, err error) {
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	addr, found := client.pick()
	if !found {
		return "", ErrEndpointNotFound
	}
	url = addr + uri
	// Option: https
	if !client.https {
		url = "http://" + url
	} else {
		url = "https://" + url
	}
	return
}

// Get
func (client *HTTPClient) Get(uri string, headers http.Header) (resp *http.Response, err error) {
	url, err := client.makeURL(uri)
	if err != nil {
		return nil, errors.Wrap(err, "GET - url creation failed")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET - request creation failed")
	}
	// Headers
	req.Header = client.makeHeaders(headers)
	// Option: host
	if client.host != "" {
		req.Host = client.host
	}

	return client.client.Do(req)
}

// Post
func (client *HTTPClient) Post(uri string, params url.Values, headers http.Header) (resp *http.Response, err error) {
	return client.PostWithBody(uri, strings.NewReader(params.Encode()), headers)
}

// PostWithBody
func (client *HTTPClient) PostWithBody(uri string, body io.Reader, headers http.Header) (resp *http.Response, err error) {
	url, err := client.makeURL(uri)
	if err != nil {
		return nil, errors.Wrap(err, "POST - url creation failed")
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "POST - request creation failed")
	}
	// Headers
	req.Header = client.makeHeaders(headers)
	// Option: host
	if client.host != "" {
		req.Host = client.host
	}

	return client.client.Do(req)
}

// Put
func (client *HTTPClient) Put(uri string, params url.Values, headers http.Header) (resp *http.Response, err error) {
	return client.PutWithBody(uri, strings.NewReader(params.Encode()), headers)
}

// PutWithBody
func (client *HTTPClient) PutWithBody(uri string, body io.Reader, headers http.Header) (resp *http.Response, err error) {
	url, err := client.makeURL(uri)
	if err != nil {
		return nil, errors.Wrap(err, "PUT - url creation failed")
	}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PUT - request creation failed")
	}
	// Headers
	req.Header = client.makeHeaders(headers)
	// Option: host
	if client.host != "" {
		req.Host = client.host
	}

	return client.client.Do(req)
}

// Patch
func (client *HTTPClient) Patch(uri string, params url.Values, headers http.Header) (resp *http.Response, err error) {
	return client.PatchWithBody(uri, strings.NewReader(params.Encode()), headers)
}

// PatchWithBody
func (client *HTTPClient) PatchWithBody(uri string, body io.Reader, headers http.Header) (resp *http.Response, err error) {
	url, err := client.makeURL(uri)
	if err != nil {
		return nil, errors.Wrap(err, "PATCH - url creation failed")
	}

	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PATCH - request creation failed")
	}
	// Headers
	req.Header = client.makeHeaders(headers)
	// Option: host
	if client.host != "" {
		req.Host = client.host
	}

	return client.client.Do(req)
}

// Delete
func (client *HTTPClient) Delete(uri string, headers http.Header) (resp *http.Response, err error) {
	url, err := client.makeURL(uri)
	if err != nil {
		return nil, errors.Wrap(err, "DELETE - url creation failed")
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DELETE - request creation failed")
	}
	// Headers
	req.Header = client.makeHeaders(headers)
	// Option: host
	if client.host != "" {
		req.Host = client.host
	}

	return client.client.Do(req)
}
