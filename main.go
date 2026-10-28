package main

import (
	"fmt"
	"log"
	"net/http"

	"api-gateway-go/config"
	"api-gateway-go/middleware"
	"api-gateway-go/proxy"

	"golang.org/x/time/rate"
)

func main() {
	cfg := config.LoadConfig()

	ipRateLimiter := middleware.NewIPRateLimiter(rate.Limit(5), 10)

	userProxy, err := proxy.NewServiceProxy(cfg.UserService, "UserServiceBreaker", "/users")
	if err != nil {
		log.Fatalf("Erro ao configurar proxy de usuarios: %v", err)
	}

	orderProxy, err := proxy.NewServiceProxy(cfg.OrderService, "OrderServiceBreaker", "/orders")
	if err != nil {
		log.Fatalf("Erro ao configurar proxy de pedidos: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "UP"}`))
	})

	jwtMiddleware := middleware.JWTMiddleware(cfg.JWTSecret)

	mux.Handle("/users/", userProxy)
	mux.Handle("/orders/", jwtMiddleware(orderProxy))

	handler := middleware.RateLimitMiddleware(ipRateLimiter)(mux)

	fmt.Printf("🚀 API Gateway rodando na porta %s...\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}