-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE status_enum AS ENUM (
	'Scheduled',
	'CheckIn',
	'Boarding',
	'Delayed',
	'Departed',
	'Arrived',
	'Cancelled',
	'Redirected'
);

CREATE TYPE cargo_type_enum AS ENUM (
	'LUGGAGE',
	'CARGO', 
	'MAIL', 
	'EQUIPMENT', 
	'DANGEROUS'
);

CREATE TABLE flights (
	id UUID DEFAULT gen_random_uuid() NOT NULL,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	deleted_at TIMESTAMPTZ,
	flight_number VARCHAR(8) NOT NULL,
	origin VARCHAR(3) NOT NULL,
	destination VARCHAR(3) NOT NULL,
	date TIMESTAMPTZ NOT NULL,
	status status_enum DEFAULT 'Scheduled',
	aircraft TEXT NOT NULL,
	CONSTRAINT flights_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_flights_deleted_at ON flights (deleted_at);

CREATE TABLE cargo_manifests (
    id UUID DEFAULT gen_random_uuid() NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    flight_id UUID NOT NULL,
    max_weight_kg INTEGER NOT NULL,
    current_weight_kg INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT cargo_manifests_pkey
        PRIMARY KEY (id),

    CONSTRAINT cargo_manifests_flight_id_fkey
        FOREIGN KEY (flight_id)
        REFERENCES flights(id)
        ON DELETE CASCADE,

    CONSTRAINT cargo_manifests_flight_id_unique
        UNIQUE (flight_id),

    CONSTRAINT cargo_manifests_weight_check
        CHECK (current_weight_kg >= 0),

    CONSTRAINT cargo_manifests_max_weight_check
        CHECK (max_weight_kg >= 0),

    CONSTRAINT cargo_manifests_current_weight_limit_check
        CHECK (current_weight_kg <= max_weight_kg)
);

CREATE INDEX idx_cargo_manifests_deleted_at
    ON cargo_manifests (deleted_at);

CREATE INDEX idx_cargo_manifests_flight_id
    ON cargo_manifests (flight_id);


CREATE TABLE cargo_items (
    item_id UUID DEFAULT gen_random_uuid() NOT NULL,

    manifest_id UUID NOT NULL,
    cargo_type cargo_type_enum NOT NULL,
    weight_kg INTEGER NOT NULL DEFAULT 0,
    passenger_id UUID,
    description TEXT NOT NULL,
    packed_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT cargo_items_pkey
        PRIMARY KEY (item_id),

    CONSTRAINT cargo_items_manifest_id_fkey
        FOREIGN KEY (manifest_id)
        REFERENCES cargo_manifests(id)
        ON DELETE CASCADE,

    CONSTRAINT cargo_items_weight_check
        CHECK (weight_kg >= 0)
);

CREATE INDEX idx_cargo_items_manifest_id
    ON cargo_items (manifest_id);

INSERT INTO flights (flight_number, origin, destination, date, status, aircraft) VALUES
('SU1001', 'SVO', 'JFK', CURRENT_TIMESTAMP + INTERVAL '2 hours', 'Scheduled', 'Boeing 777-300ER'),
('BA2045', 'LHR', 'CDG', CURRENT_TIMESTAMP - INTERVAL '3 hours', 'Arrived', 'Airbus A380'),
('LH3201', 'FRA', 'MUC', CURRENT_TIMESTAMP + INTERVAL '30 minutes', 'Boarding', 'Airbus A320'),
('AF4567', 'CDG', 'LAX', CURRENT_TIMESTAMP - INTERVAL '8 hours', 'Departed', 'Boeing 787-9'),
('EK5123', 'DXB', 'LHR', CURRENT_TIMESTAMP + INTERVAL '1 hour', 'CheckIn', 'Boeing 777-300ER'),
('SU2101', 'VKO', 'LED', CURRENT_TIMESTAMP - INTERVAL '2 days', 'Arrived', 'Airbus A320'),
('S76589', 'DME', 'KZN', CURRENT_TIMESTAMP + INTERVAL '4 hours', 'Scheduled', 'Boeing 737-800'),
('U64234', 'SVO', 'EKB', CURRENT_TIMESTAMP - INTERVAL '1 hour', 'Arrived', 'Airbus A321'),
('FV5678', 'VKO', 'AER', CURRENT_TIMESTAMP + INTERVAL '45 minutes', 'CheckIn', 'Airbus A319'),
('N84562', 'LED', 'SVO', CURRENT_TIMESTAMP + INTERVAL '6 hours', 'Scheduled', 'Boeing 737-700'),
('AA9876', 'JFK', 'LAX', CURRENT_TIMESTAMP + INTERVAL '2 hours 30 minutes', 'Delayed', 'Boeing 787-10'),
('UA5432', 'ORD', 'SFO', CURRENT_TIMESTAMP - INTERVAL '4 hours', 'Cancelled', 'Boeing 737-900'),
('DL3456', 'ATL', 'MIA', CURRENT_TIMESTAMP + INTERVAL '3 hours', 'Delayed', 'Airbus A350'),
('AC1234', 'YYZ', 'LAX', CURRENT_TIMESTAMP - INTERVAL '5 hours', 'Redirected', 'Boeing 777-200LR'),
('KL2847', 'AMS', 'FCO', CURRENT_TIMESTAMP + INTERVAL '12 hours', 'Scheduled', 'Boeing 737-8K2'),
('IB8765', 'MAD', 'BCN', CURRENT_TIMESTAMP - INTERVAL '30 minutes', 'Arrived', 'Airbus A319'),
('OS4321', 'VIE', 'PRG', CURRENT_TIMESTAMP + INTERVAL '7 hours', 'Scheduled', 'Airbus A320'),
('TK6789', 'IST', 'ATH', CURRENT_TIMESTAMP + INTERVAL '5 hours 15 minutes', 'Boarding', 'Boeing 737-8F2'),
('LX9999', 'ZRH', 'MXP', CURRENT_TIMESTAMP - INTERVAL '1 hour 30 minutes', 'Arrived', 'Airbus A220');

-- +goose Down

DROP TABLE IF EXISTS cargo_items;
DROP TABLE IF EXISTS cargo_manifests;
DROP TABLE IF EXISTS flights;

DROP TYPE IF EXISTS cargo_type_enum;
DROP TYPE IF EXISTS status_enum;
