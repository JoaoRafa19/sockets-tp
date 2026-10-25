#!/usr/bin/env python3
"""Testa todos os métodos do servidor TCP.

Uso: python3 scripts/test_api.py [host] [porta]
"""
import json
import socket
import sys
import threading
import uuid

HOST = sys.argv[1] if len(sys.argv) > 1 else "localhost"
PORT = int(sys.argv[2]) if len(sys.argv) > 2 else 3000

passed = failed = 0


class Client:
    """Uma conexão TCP persistente: várias requisições, uma resposta por linha."""

    def __init__(self):
        self.sock = socket.create_connection((HOST, PORT), timeout=5)
        self.file = self.sock.makefile("r", encoding="utf-8")

    def raw(self, text):
        self.sock.sendall(text.encode())
        line = self.file.readline()
        return json.loads(line) if line else None

    def call(self, method, data=None):
        msg = {"method": method}
        if data is not None:
            msg["data"] = data
        return self.raw(json.dumps(msg))

    def close(self):
        self.sock.close()


def check(name, cond, resp=None):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✓ {name}")
    else:
        failed += 1
        print(f"  ✗ {name}\n      resposta: {resp}")


def is_ok(r):
    return r is not None and r.get("status") == "ok"


def is_err(r, msg):
    return r is not None and r.get("status") == "error" and msg in r.get("error", "")


def main():
    c = Client()

    print("\ncreate-event")
    r = c.call("create-event", {"name": "Show de teste", "tickets": [
        {"type": "pista", "quantity": 1000},
        {"type": "camarote", "quantity": 2},
    ]})
    check("cria evento", is_ok(r) and r["data"].get("id"), r)
    event_id = r["data"]["id"] if is_ok(r) else str(uuid.uuid4())

    r = c.call("create-event", {"name": "  ", "tickets": [{"type": "a", "quantity": 1}]})
    check("recusa nome vazio", is_err(r, "invalid event name"), r)
    r = c.call("create-event", {"name": "x", "tickets": []})
    check("recusa evento sem ingressos", is_err(r, "at least one ticket"), r)
    r = c.call("create-event", {"name": "x", "tickets": [{"type": "a", "quantity": 0}]})
    check("recusa quantidade zero", is_err(r, "quantity must be greater"), r)
    r = c.call("create-event", {"name": "x", "tickets": [{"type": "a", "quantity": 1}, {"type": "a", "quantity": 2}]})
    check("recusa tipo de ingresso repetido", is_err(r, "duplicate ticket type"), r)

    print("\nreserve-ticket")
    r = c.call("reserve-ticket", {"event_id": event_id, "email": "Joao@Email.com", "type": "camarote", "quantity": 2})
    check("reserva 2 camarotes", is_ok(r) and r["data"].get("status") == "active", r)
    check("email salvo em minúsculo", is_ok(r) and r["data"].get("email") == "joao@email.com", r)
    res_id = r["data"]["id"] if is_ok(r) else str(uuid.uuid4())

    r = c.call("reserve-ticket", {"event_id": event_id, "email": "a@b.com", "type": "camarote", "quantity": 1})
    check("sem estoque -> not enough tickets", is_err(r, "not enough tickets"), r)
    r = c.call("reserve-ticket", {"event_id": event_id, "email": "a@b.com", "type": "vip", "quantity": 1})
    check("tipo inexistente -> invalid ticket type", is_err(r, "invalid ticket type"), r)
    r = c.call("reserve-ticket", {"event_id": str(uuid.uuid4()), "email": "a@b.com", "type": "pista", "quantity": 1})
    check("evento inexistente -> event not found", is_err(r, "event not found"), r)
    r = c.call("reserve-ticket", {"event_id": "nao-e-uuid", "email": "a@b.com", "type": "pista", "quantity": 1})
    check("id inválido -> event not found", is_err(r, "event not found"), r)
    r = c.call("reserve-ticket", {"event_id": event_id, "email": "nao-e-email", "type": "pista", "quantity": 1})
    check("email inválido -> invalid email", is_err(r, "invalid email"), r)
    r = c.call("reserve-ticket", {"event_id": event_id, "email": "a@b.com", "type": "pista", "quantity": 0})
    check("quantidade zero -> recusada", is_err(r, "quantity must be greater"), r)

    print("\nget-reservation")
    r = c.call("get-reservation", {"reservation_id": res_id, "email": "joao@email.com"})
    check("consulta a reserva", is_ok(r) and r["data"].get("id") == res_id, r)
    r = c.call("get-reservation", {"reservation_id": res_id, "email": "outro@email.com"})
    check("email errado -> reservation not found", is_err(r, "reservation not found"), r)
    r = c.call("get-reservation", {"reservation_id": str(uuid.uuid4()), "email": "joao@email.com"})
    check("id inexistente -> reservation not found", is_err(r, "reservation not found"), r)

    print("\ncancel-reservation")
    r = c.call("cancel-reservation", {"reservation_id": res_id, "email": "joao@email.com"})
    check("cancela a reserva", is_ok(r) and r["data"].get("status") == "cancelled", r)
    r = c.call("cancel-reservation", {"reservation_id": res_id, "email": "joao@email.com"})
    check("cancelar de novo -> already cancelled", is_err(r, "already cancelled"), r)
    r = c.call("get-reservation", {"reservation_id": res_id, "email": "joao@email.com"})
    check("consulta mostra cancelled", is_ok(r) and r["data"].get("status") == "cancelled", r)
    r = c.call("reserve-ticket", {"event_id": event_id, "email": "a@b.com", "type": "camarote", "quantity": 2})
    check("estoque devolvido: reserva os 2 camarotes de novo", is_ok(r), r)

    print("\nprotocolo")
    r = c.call("metodo-que-nao-existe", {})
    check("método inexistente -> invalid method", is_err(r, "invalid method"), r)
    r = c.raw('{"mehtod":"create-event","data":{}}')
    check("typo em method -> invalid method", is_err(r, "invalid method"), r)
    r = c.call("create-event", {"name": "x", "tickets": [{"type": "a", "quantity": 1}], "extra": 1})
    check("campo desconhecido -> invalid payload", is_err(r, "invalid payload"), r)
    r = c.call("create-event")
    check("sem data -> invalid payload", is_err(r, "invalid payload"), r)
    r = c.call("get-reservation", {"reservation_id": res_id, "email": "joao@email.com"})
    check("conexão continua aberta depois dos erros", is_ok(r), r)
    c.close()

    b = Client()
    r = b.raw('{"method":"create-event","data":{"name":"quebrado",}}')
    check("JSON quebrado -> invalid payload", is_err(r, "invalid payload"), r)
    check("JSON quebrado -> servidor fecha a conexão", b.file.readline() == "", "conexão continuou aberta")
    b.close()

    print("\nconcorrência")
    c = Client()
    r = c.call("create-event", {"name": "Concorrência", "tickets": [{"type": "pista", "quantity": 10}]})
    c.close()
    if not is_ok(r):
        check("cria evento para o teste de concorrência", False, r)
        return
    conc_id = r["data"]["id"]

    results = []
    lock = threading.Lock()

    def buy(i):
        cl = Client()
        resp = cl.call("reserve-ticket", {"event_id": conc_id, "email": f"user{i}@test.com", "type": "pista", "quantity": 1})
        cl.close()
        with lock:
            results.append(resp)

    threads = [threading.Thread(target=buy, args=(i,)) for i in range(30)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    oks = sum(1 for x in results if is_ok(x))
    sold_out = sum(1 for x in results if is_err(x, "not enough tickets"))
    check(f"30 compras simultâneas, 10 ingressos: {oks} ok, {sold_out} sem estoque",
          oks == 10 and sold_out == 20, results)


if __name__ == "__main__":
    try:
        main()
    except ConnectionRefusedError:
        print(f"Não conectou em {HOST}:{PORT}. O servidor está rodando?")
        sys.exit(1)
    print(f"\n{passed} passaram, {failed} falharam")
    sys.exit(1 if failed else 0)
