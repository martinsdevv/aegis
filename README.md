# Aegis 
Aegis é um API Gateway que atua como proxy reverso, oferecendo health checks, e futuramente observabilidade, autenticação, rate-limiting e quotas.

---

## Status de Desenvolvimento

- [x] Upstream mock server
- [x] Health check
- [x] Reverse proxy com rewrite de path
- [x] Testes de integração do proxy
- [ ] Middlewares (request-id, logging, recover)
- [ ] API Keys
- [ ] Rate limiting

## Como Rodar

- Garanta que o Go está instalado na sua máquina

```bash
go version
```

- Da raiz do projeto, rode o servidor upstream mock

```bash
go run ./cmd/upstream-mock/main.go
```

- Da raiz do projeto, rode o gateway

```bash
go run ./cmd/gateway/main.go
```

- Portas: 
- **gateway**: `:8000`
- **upstream-mock**: `:9000`

## Teste

Caso esteja no windows, recomendo utilizar o Invoke-RestMethod pra testar, ou um software de teste como Postman ou Insomnia. Caso esteja em sistemas linux/unix like o curl funcionará perfeitamente.

### Upstream (:9000)

Endpoints:
- /ping -> retorna {"pong": true}
- /healthz -> retorna {"ok": true}
- /echo -> retorna o mesmo body inserido na requisição (POST)

```bash
curl -i http://localhost:9000/ping
curl -i http://localhost:9000/healthz
curl -i -X POST http://localhost:9000/echo -H "Content-Type: application/json" -d '{"name":"Joaozin","age":21}'
```

---

### Proxy / Gateway (:8000)

Endpoints:
- /proxy/ping -> chama o /ping da upstream -> retorna {"pong": true}
- /proxy/echo -> chama o /echo da upstream -> retorna o mesmo body inserido na requisição (POST)
- /proxy/healthz -> chama o /healthz da upstream -> retorna {"ok": true}
- /healthz -> retorna {"ok": true} (status do gateway)

```bash
curl -i http://localhost:8000/proxy/ping
curl -i -X POST http://localhost:8000/proxy/echo -H "Content-Type: application/json" -d '{"name":"Joaozin","age":21}'
curl -i http://localhost:8000/proxy/healthz
curl -i http://localhost:8000/healthz
```
