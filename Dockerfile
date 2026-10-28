# Etapa 1: Build da aplicação
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia os arquivos do módulo
COPY go.mod ./
# Se o go.sum existir, copia ele também (sem falhar caso não exista)
COPY go.mod go.sum* ./

# Copia todo o código-fonte antes de resolver as dependências
COPY . .

# Baixa e sincroniza todas as dependências automaticamente
RUN go mod tidy

# Compila o binário estático
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api-gateway .

# Etapa 2: Imagem final ultraleve
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]