package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"api-gateway-go/resilience"
)

type ServiceProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	cb     *resilience.CircuitBreaker
	prefix string
}

func NewServiceProxy(targetURL string, cbName string, prefix string) (*ServiceProxy, error) {
	parsedURL, err := errURL(targetURL)
	if err != nil {
		return nil, err
	}

	sp := &ServiceProxy{
		target: parsedURL,
		proxy:  httputil.NewSingleHostReverseProxy(parsedURL),
		cb:     resilience.NewCircuitBreaker(cbName),
		prefix: prefix,
	}

	return sp, nil
}

func errURL(targetURL string) (*url.URL, error) {
	return url.Parse(targetURL)
}

func (sp *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Host = sp.target.Host

	// Remove o prefixo (/users ou /orders) ao repassar ao microserviço downstream
	if sp.prefix != "" && strings.HasPrefix(r.URL.Path, sp.prefix) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, sp.prefix)
		if r.URL.Path == "" || r.URL.Path[0] != '/' {
			r.URL.Path = "/" + r.URL.Path
		}
	}

	_, err := sp.cb.Execute(func() (*http.Response, error) {
		recorder := &responseCapture{ResponseWriter: w, statusCode: http.StatusOK}
		sp.proxy.ServeHTTP(recorder, r)

		if recorder.statusCode >= 500 {
			return &http.Response{StatusCode: recorder.statusCode}, http.ErrHandlerTimeout
		}
		return &http.Response{StatusCode: recorder.statusCode}, nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{
			"error": "Service Unavailable",
			"message": "Servico temporariamente indisponivel (Circuit Breaker ativado)",
			"fallback": true
		}`))
	}
}

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (rc *responseCapture) WriteHeader(code int) {
	rc.statusCode = code
	rc.ResponseWriter.WriteHeader(code)
}