# Estado atual (`main`)

Descrição objetiva do que o repositório implementa **hoje**, antes das adaptações ao contrato de gateway da plataforma.

## Modelo atual

O Aegis na `main` opera como **gateway orientado a API Key de consumidor**:

1. Cliente envia `X-API-Key`.
2. Chave é validada (PostgreSQL + cache Redis).
3. Rate limit e quota mensal são aplicados por chave.
4. O request é proxyado para o `upstream_host` associado àquela chave (`/proxy/*`).

Esse modelo é útil como laboratório de gateway e metering. **Não** corresponde ao papel de borda da casca (JWT + roteamento por módulo).

## Capacidades presentes

| Área | Implementação |
|------|----------------|
| Reverse proxy | `httputil.ReverseProxy` em `/proxy` e `/proxy/*` |
| Autenticação | `X-API-Key` obrigatória (middleware global) |
| Persistência de chaves | PostgreSQL; migrations e seeds no boot |
| Cache de chave | Redis com TTL |
| Rate limiting | Token bucket em memória por ID de API key |
| Quota | Contador mensal por API key (Redis + fallback) |
| Request ID | `X-Request-ID` gerado/propagado |
| Resiliência de processo | Recover de panic; graceful shutdown |
| Saúde | `/healthz` |
| Observabilidade | Logs estruturados (`slog`); stream de usage no Redis |
| Configuração | `AEGIS_LISTEN_PORT`, `AEGIS_DATABASE_URL`, `AEGIS_REDIS_ADDR` |

## Limitações em relação ao contrato-alvo

| Necessidade da plataforma | Situação na `main` |
|---------------------------|--------------------|
| Roteamento `/api/{modulo}/**` por env | Ausente (upstream por API key) |
| Validação JWT / JWKS | Ausente |
| Rotas `/public/**` sem auth | Ausente (API key global) |
| CORS centralizado | Ausente |
| Timeout de proxy + `503` | Ausente (comportamento default do transport) |
| Rate limit por IP em rotas públicas | Ausente |
| Tenant exclusivamente do JWT | Ausente |
| Boot sem depender de Postgres/Redis para o caminho JWT | Postgres/Redis são centrais ao modelo atual |

## Estrutura de código relevante

```text
cmd/gateway/                 # processo do gateway
cmd/upstream-mock/           # upstream de teste
internal/config/             # env
internal/gateway/gtwhttp/    # rotas HTTP
internal/gateway/middleware/ # API key, RL, quota, logger, request id
internal/gateway/proxy/      # reverse proxy
internal/db/migrations/      # schema de api_keys
internal/seed/               # seeds de chaves
```

## Conclusão

A base de **proxy, middleware chain, request id e operação** é reaproveitável. O produto de autenticação/roteamento precisa ser realinhado: de **API key + upstream por chave** para **JWT + upstream por prefixo de módulo**.
