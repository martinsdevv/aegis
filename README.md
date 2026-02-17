# Aegis

Aegis é um API Gateway escrito em Go que atua como um reverse proxy seguro, oferecendo autenticação por API Key, rate limiting por consumidor, health checks com controle de readiness e uma pipeline própria de middlewares HTTP.

O objetivo do projeto é demonstrar arquitetura modular em Go, controle explícito de fluxo de requisições e fundamentos de segurança em APIs HTTP.

## ✨ Funcionalidades Implementadas

* **Reverse Proxy:** Com rewrite de path automático.
* **Healthcheck:** Endpoint de saúde do gateway (`/healthz`) com controle de readiness.
* **Middleware Chain:** Implementação customizada (Chain Pattern) para processamento de requisições.
* **Segurança:** 
  * Autenticação via header `X-API-Key`.
  * Validação de chave por keyring configurável.
  * Sanitização de headers sensíveis.
  * Recovery middleware para tratamento de panics.
* **Rate Limiting:** Controle por consumidor (Token Bucket) usando `golang.org/x/time/rate`.
* **Observabilidade:** Logging estruturado de requisições e enriquecimento de resposta com headers customizados (ex: `X-Content-Id`).

---

## 🏗 Arquitetura

O projeto segue uma organização idiomática em Go:

* `cmd/`: Entrypoints da aplicação (Gateway e Upstream Mock).
* `internal/`: Implementação do domínio, middlewares e configurações.
* **Middleware Chain:** Composta manualmente para controle total da ordem de execução.
* **Rate Limit Store:** Em memória com TTL e rotina de cleanup automática.

---

## ⚙️ Variáveis de Ambiente

Antes de rodar o gateway, configure as seguintes variáveis obrigatórias:

| Variável | Descrição | Exemplo |
| :--- | :--- | :--- |
| `AegisListenPort` | Porta onde o gateway será executado | `8000` |
| `AegisUpstreamURL` | URL base do serviço de upstream | `http://localhost:9000` |
| `AegisUpstreamPort` | Porta do serviço de upstream | `9000` |
| `AegisAPIKeys` | Lista de API Keys válidas (separadas por vírgula) | `K1,K2,K3` |

### Como configurar:

**Linux / macOS (Bash/Zsh)**
```bash
export AegisListenPort=8000
export AegisUpstreamURL=http://localhost:9000
export AegisUpstreamPort="9000"
export AegisAPIKeys=K1,K2,K3
```

**Windows (PowerShell)**
```bash
$env:AegisListenPort="8000"
$env:AegisUpstreamURL="http://localhost:9000"
$env:AegisUpstreamPort="9000"
$env:AegisAPIKeys="K1,K2,K3"
```
---

## 🚀 Como Rodar

1. **Verifique a instalação do Go:**
```bash
go version
```
 
2. **Configure as variáveis de ambiente** (conforme seção acima).

3. **Rode o Upstream Mock (Serviço de teste):**

```bash
go run ./cmd/upstream-mock/main.go
```

4. **Rode o Gateway:**
```bash
go run ./cmd/gateway/main.go
```

---

## 🔎 Endpoints e Testes

### Upstream (:9000)
* `GET /ping` → `{"pong": true}`
* `GET /healthz` → `{"ok": true}`
* `POST /echo` → Retorna o mesmo body enviado.

### Gateway (:8000)
> ⚠️ Todas as rotas protegidas exigem o header `X-API-Key`.

**Teste de Healthcheck do Gateway:**
```bash
curl -i -H "X-API-Key: K1" http://localhost:8000/healthz
```

**Teste de Proxy (Exemplo Echo):**

```bash
curl -i -X POST http://localhost:8000/proxy/echo \
        -H "X-API-Key: K1" \
        -H "Content-Type: application/json" \
        -d '{"name":"Joao","age":21}'
```

---

## 🔐 Segurança e Rate Limiting

Aegis implementa as seguintes camadas de proteção:

* **Bloqueio de requisições sem Key:** Retorna `401 Unauthorized`.
* **Chaves Inválidas:** Retorna `403 Forbidden`.
* **Rate Limiting:** Atualmente configurado para **5 req/s** com burst de **10**. Caso excedido, retorna `429 Too Many Requests`.
* **Privacidade:** O header `X-API-Key` é removido antes da requisição ser encaminhada ao Upstream.
* **Resiliência:** Recovery middleware contra panics inesperados.

---

## 📌 Próximos Passos
* [ ] Quotas por consumidor
* [ ] Observabilidade (metrics / tracing)
* [ ] Persistência de API keys
* [ ] Dockerfile
* [ ] Configuração via arquivo `.env`
* [ ] Testes unitários para middlewares
* [ ] Circuit breaker
