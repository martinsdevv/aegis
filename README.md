# Aegis

Aegis é um **API Gateway modular escrito em Go** que atua como um reverse proxy seguro, oferecendo:

* Autenticação via **API Key**
* Rate limiting por consumidor
* Controle de **quota mensal**
* Cache distribuído via Redis
* Persistência de API Keys em PostgreSQL
* Logging estruturado (`slog`) e observabilidade

O projeto demonstra **arquitetura modular em Go**, fluxo de requisições explícito, fundamentos de segurança em APIs HTTP e preparação para **escalabilidade horizontal**.

---

# ✨ Funcionalidades

## Segurança

* Header `X-API-Key` obrigatório
* Validação de API Keys no PostgreSQL
* Cache em Redis com TTL para performance
* Remoção do header antes do envio ao upstream
* Middleware de recovery para panics
* Sanitização de headers sensíveis

## Rate Limiting

* Token Bucket por API Key
* Store em memória com TTL e cleanup automático
* Status `429 Too Many Requests` quando excedido

## Quota Mensal

* Controle de consumo mensal por API Key
* Redis primário e fallback in-memory
* Chave formatada: `quota:<api_key_id>:<YYYY-MM>`
* Retorno `403 Forbidden` quando excedido

## Reverse Proxy

* Encaminha `/proxy/*` → `/` do upstream
* Upstream configurável por API Key
* Enriquecimento de headers para rastreabilidade

## Health & Readiness

* Endpoint `/healthz` para checagem
* Simulação de readiness para orquestradores

## Observabilidade

* Logging estruturado (`slog`)
* Inclui método, path, host, status, latência e API Key
* Propagação de contexto interno

---

# Arquitetura

```
cmd/
    gateway/           # Entrada principal do gateway
    upstream-mock/     # Mock de upstream para testes
internal/
    gateway/           # Router e handlers
    middleware/        # Rate limiting, quota, logging, auth
    proxy/             # Reverse proxy dinâmico
    db/                # Migrations e seeds
```

* Middleware chain manual (Chain Pattern)
* Rate limit em memória com TTL
* Quota mensal via Redis + fallback
* API Key Store PostgreSQL + cache Redis
* Reverse proxy customizado com `httputil.ReverseProxy`

Arquitetura **modular e stateless**, pronta para escalabilidade horizontal.

---

# Dependências

* Go 1.21+
* PostgreSQL
* Redis

---

# ⚙️ Variáveis de Ambiente

| Variável             | Descrição                    | Exemplo                                     |
| -------------------- | ---------------------------- | ------------------------------------------- |
| `AEGIS_LISTEN_PORT`  | Porta do gateway             | `8000`                                      |
| `AEGIS_DATABASE_URL` | String de conexão PostgreSQL | `postgres://user:pass@localhost:5432/aegis` |
| `AEGIS_REDIS_ADDR`   | Endereço Redis               | `localhost:6379`                            |

---

# Banco de Dados

* Migrations e seeds **executados automaticamente** no startup do gateway.
* Tabelas principais:

```sql
CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    key TEXT UNIQUE NOT NULL,         -- SHA256 da chave
    upstream_host TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    monthly_quota INTEGER NOT NULL DEFAULT 10000,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

* Seeds padrão criadas no startup:

| Nome          | Raw Key        | Upstream                                             | Quota |
| ------------- | -------------- | ---------------------------------------------------- | ----- |
| default-dev   | `DEV_KEY_123`  | [https://httpbin.org](https://httpbin.org)           | 10000 |
| internal-test | `TEST_KEY_456` | [https://postman-echo.com](https://postman-echo.com) | 5000  |

> **Nota:** O gateway já cria a migration e insere estas chaves automaticamente ao iniciar.

* **Inserir novas chaves manualmente:**

```sql
INSERT INTO api_keys (name, key, upstream_host, monthly_quota, is_active)
VALUES (
    'minha-chave',
    '<SHA256 da chave>',
    'https://meu-upstream.com',
    10000,
    TRUE
);
```

* Para gerar SHA256 de uma key em Go:

```go
import (
    "crypto/sha256"
    "encoding/hex"
)

raw := "NOVA_KEY_123"
sum := sha256.Sum256([]byte(raw))
fmt.Println(hex.EncodeToString(sum[:]))
```

---

# Como Rodar

## 1. Subir dependências

* PostgreSQL e Redis localmente.

## 2. Configurar variáveis

Linux/macOS:

```bash
export AEGIS_LISTEN_PORT=8000
export AEGIS_DATABASE_URL=postgres://user:pass@localhost:5432/aegis
export AEGIS_REDIS_ADDR=localhost:6379
```

Windows PowerShell:

```powershell
$env:AEGIS_LISTEN_PORT="8000"
$env:AEGIS_DATABASE_URL="postgres://user:pass@localhost:5432/aegis"
$env:AEGIS_REDIS_ADDR="localhost:6379"
```

## 3. Rodar Upstream Mock (opcional)

```bash
go run ./cmd/upstream-mock/main.go
```

## 4. Rodar Gateway

```bash
go run ./cmd/gateway/main.go
```

> O gateway aplica migrations, insere seeds e inicia listeners automaticamente.

---

# Testes

## Healthcheck

```bash
curl -H "X-API-Key: DEV_KEY_123" http://localhost:8000/healthz
```

## Proxy GET

```bash
curl -i -H "X-API-Key: DEV_KEY_123" http://localhost:8000/proxy/get
```

---

# Modelo de Segurança

* Sem API Key → `401 Unauthorized`
* API Key inválida → `403 Forbidden`
* Rate limit excedido → `429 Too Many Requests`
* Quota mensal excedida → `403 Forbidden`
* Headers sensíveis removidos antes do upstream

---

# Observabilidade

* Logs estruturados via `slog`
* Informações logadas: método, path, host/upstream, status, duração, API Key, quota
* Middleware central de logging evita duplicidade

---

# Roadmap

* [ ] Rate limit distribuído via Redis
* [ ] Métricas Prometheus
* [ ] Request ID global
* [ ] Circuit breaker
* [ ] Graceful shutdown
* [ ] Dockerfile
* [ ] Testes unitários e integração

---

# Objetivo do Projeto

Demonstrar construção de **um API Gateway modular em Go** com:

* Separação clara de responsabilidades
* Controle explícito de fluxo
* Segurança aplicada
* Preparação para ambientes escaláveis
* Integração com Redis e PostgreSQL
