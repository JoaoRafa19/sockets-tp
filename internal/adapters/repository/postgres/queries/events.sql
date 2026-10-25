-- name: CreateEvent :one
INSERT INTO events (name)
VALUES ($1)
RETURNING id, name, created_at;

-- name: CreateEventTicket :exec
INSERT INTO event_tickets (event_id, type, quantity)
VALUES ($1, $2, $3);

-- name: GetEventTicket :one
SELECT event_id, type, quantity
FROM event_tickets
WHERE event_id = $1 AND type = $2;

-- name: ListEventTickets :many
SELECT type, quantity
FROM event_tickets
WHERE event_id = $1
ORDER BY type;

-- Desconta do estoque só se houver ingressos suficientes.
-- Retorna o número de linhas afetadas: 0 = não reservou.
-- name: DecrementTickets :execrows
UPDATE event_tickets
SET quantity = quantity - sqlc.arg(amount)::int
WHERE event_id = sqlc.arg(event_id)
  AND type = sqlc.arg(type)
  AND quantity >= sqlc.arg(amount)::int;

-- name: IncrementTickets :exec
UPDATE event_tickets
SET quantity = quantity + sqlc.arg(amount)::int
WHERE event_id = sqlc.arg(event_id)
  AND type = sqlc.arg(type);


-- name: EventExists :one
SELECT EXISTS (SELECT 1 FROM events WHERE id = $1);