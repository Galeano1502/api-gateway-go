package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := []byte("minha_chave_secreta_super_segura_123")

	claims := jwt.MapClaims{
		"sub":  "user_12345",
		"name": "Joao Ferraz",
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		fmt.Println("Erro ao gerar token:", err)
		return
	}

	fmt.Println("--- TOKEN JWT DE TESTE ---")
	fmt.Println(tokenString)
	fmt.Println("\nComo usar no PowerShell:")
	fmt.Printf("curl.exe -i -H \"Authorization: Bearer %s\" http://localhost:8080/orders/list\n", tokenString)
}