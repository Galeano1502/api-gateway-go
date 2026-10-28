package resilience

import (
	"errors"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(name string) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                 // Maximo de requisicoes permitidas no estado HALF-OPEN
		Interval:    5 * time.Second,   // Janela de contagem de erros no estado CLOSED
		Timeout:     10 * time.Second,  // Tempo que permanece OPEN antes de mudar para HALF-OPEN
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 5 && failureRatio >= 0.5 // Abre o circuito se 50%+ falharem de no minimo 5 reqs
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			println("Circuit Breaker [" + name + "] mudou de estado: " + from.String() + " -> " + to.String())
		},
	}

	return &CircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *CircuitBreaker) Execute(reqFunc func() (*http.Response, error)) (*http.Response, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		resp, err := reqFunc()
		if err != nil {
			return nil, err
		}
		// Considera status HTTP 5xx como falhas do servico de destino
		if resp.StatusCode >= 500 {
			return resp, errors.New("resposta de erro no servico downstream")
		}
		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*http.Response), nil
}