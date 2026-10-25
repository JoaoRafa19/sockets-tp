-- name: CreateReservation :one
INSERT INTO reservations (event_id, type, email, quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetReservation :one
SELECT *
FROM reservations
WHERE id = $1 AND email = $2;

-- Só cancela reservas ativas. Sem linha retornada = não existe ou já cancelada.
-- name: CancelReservation :one
UPDATE reservations
SET status = 'cancelled'
WHERE id = $1 AND email = $2 AND status = 'active'
RETURNING *;