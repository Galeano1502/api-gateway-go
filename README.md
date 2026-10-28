# 🛡️ Resilient Go API Gateway

Um API Gateway ultraleve, de altíssima performance e baixa latência desenvolvido em **Go**, projetado para atuar como ponto único de entrada (*Entrypoint*) em arquiteturas de microserviços.

---

## 🎯 Principais Recursos

- **⚡ Proxy Reverso Nativo**: Roteamento dinâmico e transparente usando `net/http/httputil`.
- **🚦 Rate Limiting (Token Bucket)**: Proteção por IP contra abusos/DDoS com limpeza automática de clientes inativos (`golang.org/x/time/rate`).
- **⚡ Circuit Breaker & Fallback**: Prevenção de falhas em cascata utilizando o padrão Circuit Breaker (`gobreaker`). Retorna fallback gracioso (`503 Service Unavailable`) quando o serviço de destino apresenta falhas contínuas.
- **🔑 Validação Transparente de JWT**: Interceptação e validação de tokens JWT no gateway com propagação de contexto (`X-User-Id`, `X-User-Role`) via cabeçalhos HTTP para os serviços downstream.

---

## 🏛️ Arquitetura

```text
[ Cliente ] 
     │
     ▼
┌─────────────────────────────────────────────────────────┐
│                     API GATEWAY                         │
│                                                         │
│  1. Token Bucket (Rate Limiter por IP)                  │
│  2. JWT Validation & Header Injection                   │
│  3. Router (Proxy Reverso)                              │
│  4. Circuit Breaker (CLOSED / OPEN / HALF-OPEN)         │
└────────────────┬───────────────────┬────────────────────┘
                 │                   │
                 ▼                   ▼
       [ User-Service ]     [ Order-Service ]