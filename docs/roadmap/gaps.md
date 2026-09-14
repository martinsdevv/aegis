# Lacunas e roadmap (Aegis como gateway)

Gap analysis entre o estado da `main` e o contrato descrito em [gateway-contract.md](../architecture/gateway-contract.md). Escopo: **somente o Aegis**.

## Prioridade essencial

Itens necessários para o Aegis cumprir o papel de gateway de borda da plataforma.

| # | Lacuna | Direção de solução | Áreas de código (orientação) |
|---|--------|--------------------|------------------------------|
| 1 | Roteamento `/api/{modulo}/**` com destino por ambiente | Parser de mapa `modulo=url`; match por prefixo; rewrite conforme política acordada | `internal/config`, novo `internal/gateway/routing`, `gtwhttp/router.go`, `proxy/handler.go` |
| 2 | Validação JWT via JWKS | Middleware Bearer; cache de chaves; `iss`/`aud`/`exp`; sem emissão de token | novo pacote ex. `internal/gateway/jwt` |
| 3 | Rotas públicas `/public/**` e health sem credencial | Cadeia de middleware com bypass por path | `middleware` + router |
| 4 | CORS centralizado | Middleware + origens/métodos/headers por env | `middleware/cors.go`, config |
| 5 | Timeout de proxy (~3s) e `503` se módulo cair | `http.Transport` / client com deadline; `ErrorHandler` do reverse proxy | `proxy/handler.go` |
| 6 | Rate limit por IP em `/public/**` | Store existente adaptado a chave = IP; separado do RL por API key | `middleware/ratelimit.go` ou equivalente |
| 7 | Tenant apenas do JWT | Remover/ignorar headers e não confiar em body/query de tenant na borda | middleware JWT + Director do proxy |
| 8 | Desacoplar boot do caminho “plataforma” de Postgres/Redis obrigatórios | DB/Redis opcionais ou restritos ao modo legado API key | `cmd/gateway/main.go`, config |

## Prioridade desejável

| # | Lacuna | Nota |
|---|--------|------|
| 9 | Encaminhar claims já validadas como headers internos (`X-User-Id`, `X-Tenant-Id`) | Conveniência; módulo continua obrigado a validar JWT |
| 10 | Envelope de erro padronizado (`401`/`403`/`503`) | Alinhamento com contrato de integração da plataforma |
| 11 | Métricas Prometheus / tracing | Roadmap de observabilidade |
| 12 | Testes de integração: rota, JWT, timeout→503, CORS preflight, public RL | Evitar regressão no caminho de borda |
| 13 | Documentação operacional (env vars, exemplos curl, health) | Atualizar README + este `docs/` |

## Itens a não implementar no Aegis

Ver [boundaries.md](../architecture/boundaries.md). Em particular:

- Login, refresh, usuários, perfis, tenants.
- UI da casca / registro visual de módulos.
- Quota mensal por API key como auth do browser.
- Obrigar tráfego inter-módulos a passar pelo gateway.

## Ordem sugerida de entrega

1. Config de módulos + roteamento `/api/{modulo}` + timeout/503 (proxy “certo”).
2. JWT/JWKS + bypass `/public` e `/healthz`.
3. CORS + rate limit IP em público.
4. Higiene de tenant/headers + testes.
5. Decisão explícita sobre o legado API key (manter atrás de flag, mover ou descontinuar no caminho default).

## Critérios de aceite (gateway)

- Requisição autenticada a `/api/{modulo}/...` com JWT válido chega ao container configurado por env.
- Sem JWT em rota protegida → `401` no Aegis.
- JWKS indisponível não pode ser contornado com token forjado (falha fechada).
- `/public/**` não exige Bearer; excesso de chamadas por IP → `429`.
- Upstream parado ou lento além do timeout → `503`.
- Nenhum módulo precisa configurar CORS para o browser falar com a origem do gateway.
- Documentação em `docs/` descreve o papel do Aegis sem confundir com Identity ou casca.
