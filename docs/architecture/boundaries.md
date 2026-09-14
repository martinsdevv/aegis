# Limites de escopo

Define fronteiras explícitas para evitar que o Aegis absorva responsabilidades da plataforma ou do Identity.

## Pertence ao Aegis

- Escuta HTTP pública (ou de borda) única.
- Roteamento por prefixo `/api/{modulo}` (e `/public` quando aplicável).
- Validação de JWT na entrada (consumidor de JWKS).
- CORS, timeouts de proxy, respostas `401` / `429` / `503` de borda.
- Propagação de `X-Request-Id`.
- Rate limiting de superfície pública por IP.
- Sanitização de headers que tentem forjar tenant/usuário.

## Pertence ao Identity (módulo)

- Login, refresh, logout, recuperação de senha, 2FA.
- Emissão e rotação de chaves (JWKS **publicado pelo Identity**).
- Cadastro de usuários, perfis, permissões, tenants.
- `GET /auth/me` e tokens de serviço.
- Política de senha, auditoria de acesso, sessões.

O Aegis **não** implementa Identity; apenas encaminha `/api/identity/**` e usa o JWKS para validar.

## Pertence à casca (plataforma)

- SPA mãe, menu, registro visual de módulos.
- Fluxo de UX de sessão entre SPAs (boot, renovação visível ao usuário).
- Tema, layout, health visual de módulos no menu.
- Decisões de produto sobre quais módulos aparecem para qual perfil.

A casca **usa** o Aegis como origem HTTP; não vive dentro do repositório do gateway.

## Pertence a cada módulo de negócio

- Regras de domínio e persistência (schema próprio).
- Validação JWT na API (segunda camada).
- RBAC fino (`perms`) e filtro por `tenant_id`.
- Chamadas a outros módulos na rede interna.
- Front SPA do módulo.

## Não misturar com o modelo legado de API key

O modelo atual de `api_keys`, quota mensal por chave e `/proxy` com upstream por chave é um **produto diferente** (gateway de consumidores/metering).

Para o encaixe na plataforma acadêmica/produto multi-módulo:

- Não usar API key como autenticação do browser.
- Não colocar cadastro de módulos de negócio na tabela de API keys.
- Tratar o legado como experimental ou feature opcional isolada — não como caminho padrão da casca.

## Resumo

| Pergunta | Resposta |
|----------|----------|
| Onde o usuário faz login? | Identity (via casca) |
| Onde o token é emitido? | Identity |
| Onde o token é checado na borda? | Aegis |
| Onde o token é checado de novo? | API de cada módulo |
| Onde está o menu? | Casca |
| Onde está a regra de cobrança/CRM/etc.? | Módulo dono |
| Chamada financeiro → CRM? | Direta, sem Aegis |
