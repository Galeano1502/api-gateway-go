package config

import "os"

type Config struct {
	Port        string
	JWTSecret   []byte
	UserService string
	OrderService string
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		JWTSecret:    []byte(getEnv("JWT_SECRET", "minha_chave_secreta_super_segura_123")),
		UserService:  getEnv("USER_SERVICE_URL", "http://localhost:8081"),
		OrderService: getEnv("ORDER_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}