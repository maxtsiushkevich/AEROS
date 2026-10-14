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
	'Luggage',
	'Cargo', 
	'Mail', 
	'Equipment', 
	'Dangerous'
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
    id UUID DEFAULT gen_random_uuid() NOT NULL,

    manifest_id UUID NOT NULL,
    cargo_type cargo_type_enum NOT NULL,
    weight_kg INTEGER NOT NULL DEFAULT 0,
    passenger_id UUID,
    description TEXT NOT NULL,


    CONSTRAINT cargo_items_pkey
        PRIMARY KEY (id),

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

INSERT INTO cargo_manifests (created_at, updated_at, flight_id, max_weight_kg, current_weight_kg)
SELECT CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, f.id, m.max_weight_kg, 0
FROM (VALUES
    ('SU1001', 20000),
    ('BA2045', 30000),
    ('LH3201',  3000),
    ('AF4567', 15000),
    ('EK5123', 20000),
    ('KL2847',  3500),
    ('AA9876', 15000),
    ('TK6789',  3000),
    ('DL3456', 14000),
    ('AC1234', 18000)
) AS m(flight_number, max_weight_kg)
JOIN flights f ON f.flight_number = m.flight_number;

INSERT INTO cargo_items (manifest_id, cargo_type, weight_kg, passenger_id, description)
SELECT
    cm.id,
    i.cargo_type::cargo_type_enum,
    i.weight_kg,
    CASE WHEN i.cargo_type = 'Luggage' THEN gen_random_uuid() END,
    i.description
FROM (VALUES
    -- SU1001
('SU1001', 'Luggage',   23, 'Suitcase, blue, checked baggage'),
('SU1001', 'Luggage',   18, 'Travel bag, black'),
('SU1001', 'Cargo',   1200, 'Pallet with equipment spare parts'),
('SU1001', 'Mail',     150, 'Mail shipments, 12 bags'),
-- BA2045
('BA2045', 'Luggage',   25, 'Suitcase, red'),
('BA2045', 'Luggage',   20, 'Large backpack, checked'),
('BA2045', 'Cargo',   4500, 'Container with electronics'),
('BA2045', 'Equipment', 800, 'Spare landing gear for maintenance'),
-- LH3201
('LH3201', 'Luggage',   22, 'Suitcase, gray'),
('LH3201', 'Luggage',   15, 'Sports bag'),
('LH3201', 'Mail',      80, 'Business correspondence'),
-- AF4567
('AF4567', 'Luggage',   24, 'Suitcase, green'),
('AF4567', 'Cargo',   3200, 'Pharmaceutical products (temperature-controlled mode)'),
('AF4567', 'Dangerous', 120, 'Lithium batteries, Class 9'),
-- EK5123
('EK5123', 'Luggage',   30, 'Two pieces of luggage, family'),
('EK5123', 'Cargo',   6000, 'Textiles, 40 boxes'),
('EK5123', 'Equipment', 500, 'Musical equipment'),
-- KL2847
('KL2847', 'Luggage',   19, 'Suitcase, black'),
('KL2847', 'Cargo',    600, 'Cut flowers, 30 boxes'),
('KL2847', 'Mail',      90, 'Parcels'),
-- AA9876
('AA9876', 'Luggage',   21, 'Suitcase, orange'),
('AA9876', 'Cargo',   2800, 'Auto parts'),
('AA9876', 'Dangerous', 200, 'Aerosols and paints, Class 2'),
-- TK6789
('TK6789', 'Luggage',   17, 'Carry-on suitcase (checked as luggage)'),
('TK6789', 'Cargo',    450, 'Product samples for exhibition'),
-- DL3456 
('DL3456', 'Luggage',   26, 'Suitcase, large'),
('DL3456', 'Cargo',   3500, 'Medical equipment'),
('DL3456', 'Mail',     200, 'Express mail'),
-- AC1234
('AC1234', 'Luggage',   28, 'Suitcase, brown'),
('AC1234', 'Cargo',   5200, 'Industrial equipment'),
('AC1234', 'Dangerous', 300, 'Dry ice for sample transport, UN1845')
) AS i(flight_number, cargo_type, weight_kg, description)
JOIN flights f          ON f.flight_number = i.flight_number
JOIN cargo_manifests cm ON cm.flight_id = f.id;

UPDATE cargo_manifests cm
SET current_weight_kg = s.total,
    updated_at = CURRENT_TIMESTAMP
FROM (
    SELECT manifest_id, SUM(weight_kg)::INTEGER AS total
    FROM cargo_items
    GROUP BY manifest_id
) s
WHERE cm.id = s.manifest_id;

-- +goose Down

DROP TABLE IF EXISTS cargo_items;
DROP TABLE IF EXISTS cargo_manifests;
DROP TABLE IF EXISTS flights;

DROP TYPE IF EXISTS cargo_type_enum;
DROP TYPE IF EXISTS status_enum;
