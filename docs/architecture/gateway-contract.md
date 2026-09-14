# Contrato do gateway

Especificação do comportamento esperado do Aegis como gateway de borda. Destina-se a orientar implementação e integração com casca, Identity e módulos.

## Superfície HTTP

| Prefixo | Autenticação | Destino |
|---------|--------------|---------|
| `/api/{modulo}/**` | JWT Bearer obrigatório (exceto regras explícitas) | Container do módulo (`{modulo}` resolvido por configuração) |
| `/api/identity/**` | Conforme rota (login/JWKS públicos; demais protegidas) | Módulo Identity |
| `/public/**` | Sem token | Upstreams públicos configurados |
| `/healthz` | Sem token | Saúde do próprio gateway |

O mapeamento `{modulo} → URL base` vem de **variáveis de ambiente** (ou arquivo de config carregado no boot), nunca de código fixo por módulo.

Exemplo conceitual:

```text
AEGIS_MODULES=identity=http://identity:8080,financeiro=http://financeiro-api:9001,crm=http://crm-api:9002
```

## Fluxo browser → módulo

```text
SPA do módulo
    │  Authorization: Bearer <access_token>
    ▼
Aegis
    │  1. CORS / Request-Id
    │  2. Valida JWT (JWKS do Identity)
    │  3. Rejeita claims/headers de tenant enviados pelo cliente
    │  4. Proxy para o upstream do módulo (timeout)
    ▼
API do módulo
    │  Valida JWT novamente (JWKS)
    │  Aplica tenant_id e perms a partir das claims
    ▼
Resposta
```

## Fluxo back-end ↔ back-end

```text
financeiro-api  ──────HTTP direto──────►  crm-api
                     (rede interna)
                     sem passar pelo Aegis
```

O gateway **não** é ponto único de falha para integração entre módulos.

## Autenticação na borda

| Requisito | Comportamento |
|-----------|----------------|
| Token ausente ou inválido em rota protegida | `401` no gateway |
| Token válido | Encaminha ao módulo (módulo revalida) |
| JWKS | Consumido de URL do Identity (cache local; respeito a `kid` / rotação) |
| Emissão de token | Fora do Aegis (Identity) |

Claims relevantes no JWT de usuário (emitidas pelo Identity): `iss`, `aud`, `sub`, `tenant_id`, identidade, `roles` / `perms`, `iat`, `exp`.

## Isolamento por tenant

- A única fonte de `tenant_id` para requisições de usuário é o **JWT**.
- Valores de tenant no body, query string ou headers enviados pelo cliente devem ser **ignorados ou removidos** antes do proxy.
- Exceção de produto: token de serviço (Identity), em que o tenant pode ir explícito na chamada — regra definida pelo Identity, não inventada no gateway.

## Cross-Origin Resource Sharing (CORS)

- Política definida **somente no Aegis**.
- Módulos não precisam configurar CORS para o browser.
- Origens, métodos e headers permitidos vêm de configuração.

## Observabilidade e resiliência

| Aspecto | Expectativa |
|---------|-------------|
| `X-Request-Id` | Gerado ou propagado na entrada; encaminhado aos upstreams; ecoado na resposta |
| Timeout de proxy | Ordem de 3 segundos por chamada ao módulo |
| Módulo indisponível / timeout | `503` com indicação do módulo afetado (envelope estável) |
| Rate limit em `/public/**` | Por IP; `429` sem afetar rotas autenticadas |

## Responsabilidades dos integradores

### Casca

- Autenticar o usuário via Identity (através do gateway).
- Entregar sessão aos SPAs dos módulos (boot / contexto da mesma origem).
- Chamar apenas a origem pública do gateway.

### Identity

- Emitir e renovar tokens; publicar JWKS; expor `/auth/me` e demais APIs de identidade.

### Módulos

- Expor API sob o contrato acordado com o prefixo `/api/{codigo}/**`.
- Validar JWT localmente (JWKS); aplicar RBAC e filtro por `tenant_id`.
- Para falar com outro módulo: chamada direta na rede interna.
