-- Insert 20 different flights with various parameters
INSERT INTO flights (flight_number, origin, destination, date, status, aircraft) VALUES
-- International flights (various statuses)
('SU1001', 'SVO', 'JFK', CURRENT_TIMESTAMP + INTERVAL '2 hours', 'Scheduled', 'Boeing 777-300ER'),
('BA2045', 'LHR', 'CDG', CURRENT_TIMESTAMP - INTERVAL '3 hours', 'Arrived', 'Airbus A380'),
('LH3201', 'FRA', 'MUC', CURRENT_TIMESTAMP + INTERVAL '30 minutes', 'Boarding', 'Airbus A320'),
('AF4567', 'CDG', 'LAX', CURRENT_TIMESTAMP - INTERVAL '8 hours', 'Departed', 'Boeing 787-9'),
('EK5123', 'DXB', 'LHR', CURRENT_TIMESTAMP + INTERVAL '1 hour', 'CheckIn', 'Boeing 777-300ER'),

-- Regional flights
('SU2101', 'VKO', 'LED', CURRENT_TIMESTAMP - INTERVAL '2 days', 'Arrived', 'Airbus A320'),
('S76589', 'DME', 'KZN', CURRENT_TIMESTAMP + INTERVAL '4 hours', 'Scheduled', 'Boeing 737-800'),
('U64234', 'SVO', 'EKB', CURRENT_TIMESTAMP - INTERVAL '1 hour', 'Arrived', 'Airbus A321'),
('FV5678', 'VKO', 'AER', CURRENT_TIMESTAMP + INTERVAL '45 minutes', 'CheckIn', 'Airbus A319'),
('N84562', 'LED', 'SVO', CURRENT_TIMESTAMP + INTERVAL '6 hours', 'Scheduled', 'Boeing 737-700'),

-- Delayed and cancelled flights
('AA9876', 'JFK', 'LAX', CURRENT_TIMESTAMP + INTERVAL '2 hours 30 minutes', 'Delayed', 'Boeing 787-10'),
('UA5432', 'ORD', 'SFO', CURRENT_TIMESTAMP - INTERVAL '4 hours', 'Cancelled', 'Boeing 737-900'),
('DL3456', 'ATL', 'MIA', CURRENT_TIMESTAMP + INTERVAL '3 hours', 'Delayed', 'Airbus A350'),

-- Redirected flights
('AC1234', 'YYZ', 'LAX', CURRENT_TIMESTAMP - INTERVAL '5 hours', 'Redirected', 'Boeing 777-200LR'),

-- More variety
('KL2847', 'AMS', 'FCO', CURRENT_TIMESTAMP + INTERVAL '12 hours', 'Scheduled', 'Boeing 737-8K2'),
('IB8765', 'MAD', 'BCN', CURRENT_TIMESTAMP - INTERVAL '30 minutes', 'Arrived', 'Airbus A319'),
('OS4321', 'VIE', 'PRG', CURRENT_TIMESTAMP + INTERVAL '7 hours', 'Scheduled', 'Airbus A320'),
('TK6789', 'IST', 'ATH', CURRENT_TIMESTAMP + INTERVAL '5 hours 15 minutes', 'Boarding', 'Boeing 737-8F2'),
('LX9999', 'ZRH', 'MXP', CURRENT_TIMESTAMP - INTERVAL '1 hour 30 minutes', 'Arrived', 'Airbus A220');

-- Verify insertion
SELECT COUNT(*) as total_flights FROM flights;
SELECT * FROM flights ORDER BY date;