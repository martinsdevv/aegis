# Visão geral e posicionamento

## Contexto

O produto é composto por vários times. Cada time entrega um **módulo** (SPA + API). A **plataforma (casca)** é a entidade mãe: login na UI, menu, navegação entre módulos. O **Identity** é um módulo transversal de autenticação e autorização. O **Aegis** é o **API gateway de borda** — não é a plataforma e não é o Identity.

## Papéis

```text
                         ┌─────────────────────┐
                         │  Casca (plataforma) │
                         │  SPA mãe / menu     │
                         └──────────┬──────────┘
                                    │
Browser ─── uma origem HTTP ───────►│
                                    ▼
                         ┌─────────────────────┐
                         │       Aegis         │
                         │  gateway de borda   │
                         └──────────┬──────────┘
              ┌─────────────────────┼─────────────────────┐
              ▼                     ▼                     ▼
       /api/identity/**     /api/{modulo}/**        /public/**
              │                     │
              ▼                     ▼
         Identity              APIs dos módulos
         (emite JWT,           (validam JWT via JWKS)
          JWKS, /me)
```

| Componente | Responsabilidade |
|------------|------------------|
| **Casca** | Experiência do usuário, menu, orquestração de SPAs |
| **Identity** | Login, tokens, JWKS, usuários, tenant, permissões |
| **Aegis** | Única entrada HTTP: roteamento, validação na borda, CORS, limites |
| **Módulos** | Domínio de negócio; consomem Identity; não emitem token |

## Princípios

1. **Uma URL pública.** O navegador não conhece hosts internos dos módulos.
2. **Auth em um lugar só (Identity).** Nenhum módulo tem tela de senha nem emite JWT de usuário.
3. **Duas camadas na borda e no módulo.** O gateway rejeita token inválido; o módulo valida de novo (JWKS).
4. **Back-end ↔ back-end sem Aegis.** Chamadas entre containers usam o nome do serviço na rede interna.
5. **Tenant vem do token.** Corpo, query ou header enviado pelo cliente com `tenantId` são ignorados (exceto fluxos de token de serviço definidos pelo Identity).

## O que o Aegis é

Um **reverse proxy seguro e configurável** na entrada do sistema:

- Encaminha `/api/{modulo}/**` ao upstream do módulo (URL via ambiente).
- Valida JWT na entrada (JWKS publicado pelo Identity).
- Centraliza CORS e política de rotas públicas.
- Propaga identificadores de requisição e responde de forma previsível quando um módulo falha.

## O que o Aegis não é

- Não hospeda login, cadastro de usuários ou emissão de tokens.
- Não é a casca (menu, registro visual de módulos, tema).
- Não é barramento obrigatório para tráfego interno entre módulos.
- Não substitui a validação JWT dentro de cada API de módulo.
