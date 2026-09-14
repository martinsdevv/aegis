# Documentação do Aegis

Índice da documentação de arquitetura e encaixe do Aegis no ecossistema.

| Documento | Conteúdo |
|-----------|----------|
| [Visão geral e posicionamento](architecture/overview.md) | Papel do Aegis, casca, Identity e módulos |
| [Contrato do gateway](architecture/gateway-contract.md) | Rotas, JWT, o que passa / não passa |
| [Estado atual (`main`)](architecture/current-state.md) | O que o código oferece hoje |
| [Lacunas e roadmap](roadmap/gaps.md) | O que falta para o Aegis se encaixar |
| [Limites de escopo](architecture/boundaries.md) | O que não pertence ao Aegis |

**Premissas do projeto**

- Browser → API de módulo passa pelo **gateway** (`/api/{modulo}/**`).
- Back-end ↔ back-end **não** passa pelo gateway.
- Identity é módulo (emite JWT / JWKS); casca é a entidade mãe.
- Fronts são **SPAs separados** — **sem iframe**.
