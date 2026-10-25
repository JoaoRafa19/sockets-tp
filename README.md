# sockets-tp
TPSockets

Aplicativo para compra de ingressos para eventos, com comunicação via **TCP** e mensagens em **JSON**.

## Como rodar

Subir o banco e rodar as migrations:

```bash
docker compose up -d db
docker compose run --rm migrate
```

Rodar o servidor:

```bash
port=:3000 go run ./cmd/server/main.go
```

## Protocolo

- Conexão TCP na porta `3000`.
- Cada requisição é um objeto JSON com `method` e `data`.
- Uma mesma conexão pode enviar várias requisições em sequência.
- Cada resposta é um JSON terminado com `\n`.

**Requisição**
```json
{ "method": "<nome-do-metodo>", "data": { } }
```

**Resposta de sucesso**
```json
{ "status": "ok", "data": { } }
```

**Resposta de erro**
```json
{ "status": "error", "error": "mensagem do erro" }
```

## Methods

> Os IDs são UUIDs gerados pelo banco. Nos exemplos abaixo, `evt-1` e `res-1` representam o `id` devolvido pelo `create-event` e pelo `reserve-ticket`.

### create-event
Cria um evento com os tipos de ingresso e a quantidade de cada um.

```json
{
    "method": "create-event",
    "data": {
        "name": "evento",
        "tickets": [
            { "type": "pista", "quantity": 1000 },
            { "type": "camarote", "quantity": 100 }
        ]
    }
}
```

Resposta:
```json
{
    "status": "ok",
    "data": {
        "id": "evt-1",
        "name": "evento",
        "tickets": [
            { "type": "pista", "quantity": 1000 },
            { "type": "camarote", "quantity": 100 }
        ]
    }
}
```

### reserve-ticket
Reserva ingressos de um tipo para um email. A quantidade é descontada do estoque do evento.

```json
{
    "method": "reserve-ticket",
    "data": {
        "event_id": "evt-1",
        "email": "joao@email.com",
        "type": "pista",
        "quantity": 2
    }
}
```

Resposta:
```json
{
    "status": "ok",
    "data": {
        "id": "res-1",
        "event_id": "evt-1",
        "email": "joao@email.com",
        "type": "pista",
        "quantity": 2,
        "status": "active"
    }
}
```

### get-reservation
Consulta uma reserva. O email precisa ser o mesmo usado na reserva.

```json
{
    "method": "get-reservation",
    "data": {
        "reservation_id": "res-1",
        "email": "joao@email.com"
    }
}
```

Resposta:
```json
{
    "status": "ok",
    "data": {
        "id": "res-1",
        "event_id": "evt-1",
        "email": "joao@email.com",
        "type": "pista",
        "quantity": 2,
        "status": "active"
    }
}
```

### cancel-reservation
Cancela uma reserva e devolve os ingressos ao estoque do evento. O email precisa ser o mesmo usado na reserva.

```json
{
    "method": "cancel-reservation",
    "data": {
        "reservation_id": "res-1",
        "email": "joao@email.com"
    }
}
```

Resposta:
```json
{
    "status": "ok",
    "data": {
        "id": "res-1",
        "event_id": "evt-1",
        "email": "joao@email.com",
        "type": "pista",
        "quantity": 2,
        "status": "cancelled"
    }
}
```

## Erros

| Erro | Quando acontece |
|---|---|
| `invalid payload` | JSON malformado ou com campos de tipo errado (a conexão é encerrada) |
| `invalid method` | `method` não existe |
| `event not found` | `event_id` não existe |
| `invalid ticket type` | `type` não existe no evento |
| `not enough tickets` | Quantidade pedida maior que o estoque |
| `reservation not found` | `reservation_id` não existe ou o email não confere |
| `reservation already cancelled` | Reserva já foi cancelada |

## Testando pelo terminal (PowerShell)

Cole esta função uma vez no terminal:

```powershell
function tcp($j){$c=[Net.Sockets.TcpClient]::new('localhost',3000);$s=$c.GetStream();$s.ReadTimeout=5000;$b=[Text.Encoding]::UTF8.GetBytes($j);$s.Write($b,0,$b.Length);([IO.StreamReader]::new($s)).ReadLine();$c.Close()}
```

Depois é só enviar as requisições:

```powershell
tcp '{"method":"create-event","data":{"name":"evento","tickets":[{"type":"pista","quantity":1000},{"type":"camarote","quantity":100}]}}'

tcp '{"method":"reserve-ticket","data":{"event_id":"evt-1","email":"joao@email.com","type":"pista","quantity":2}}'

tcp '{"method":"get-reservation","data":{"reservation_id":"res-1","email":"joao@email.com"}}'

tcp '{"method":"cancel-reservation","data":{"reservation_id":"res-1","email":"joao@email.com"}}'
```

No Linux/WSL, com `ncat`:

```bash
echo '{"method":"create-event","data":{"name":"evento","tickets":[{"type":"pista","quantity":1000}]}}' | ncat localhost 3000
```