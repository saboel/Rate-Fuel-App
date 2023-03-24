-- name: ListClients :many
SELECT * FROM clients;

-- name: AddClient :one
INSERT INTO clients(address_client, client_name, client_city, client_state) VALUES ($1, $2, $3, $4) RETURNING *;


 -- name: GetFuelQuotes :many
 SELECT * FROM fuel_history WHERE client_id = $1;

-- name: GetClient_id :one 
 SELECT client_id FROM clients WHERE address_client = $1;
