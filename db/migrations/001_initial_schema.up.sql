CREATE TABLE IF NOT EXISTS clients (
    client_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    address_client TEXT NOT NULL UNIQUE,
    client_name TEXT NOT NULL,
    client_city TEXT NOT NULL, 
    client_state TEXT NOT NULL 
);

CREATE TABLE IF NOT EXISTS fuel_history (
fuel_history_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
client_id UUID NOT NULL REFERENCES clients(client_id),
gallons FLOAT NOT NULL,
price_per_gallon FLOAT NOT NULL,
total_price FLOAT NOT NULL,
delivery_date DATE NOT NULL,
timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);