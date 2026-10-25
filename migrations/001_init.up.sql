CREATE TABLE events (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE event_tickets (
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    type     TEXT NOT NULL,
    quantity INT  NOT NULL CHECK (quantity >= 0),
    PRIMARY KEY (event_id, type)
);

CREATE TABLE reservations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID NOT NULL,
    type       TEXT NOT NULL,
    email      TEXT NOT NULL,
    quantity   INT  NOT NULL CHECK (quantity > 0),
    status     TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (event_id, type) REFERENCES event_tickets(event_id, type)
);