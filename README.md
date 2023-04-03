# Rate-Fuel-App
predict the rate of the fuel based on certain criteria:

Client location (in-state or out-of-state) 

Client history (existing customer with previous purchase or new) 

Gallons requested 

Company profit margin(%) 

# How to start server
Currently the application expects a POSTGRES_URL environment variable. We are hosting the database locally
using docker. Here are the steps to follow:

Install docker if you don't have it installed:
https://docs.docker.com/get-docker/

I am also assuming you have golang properly installed on your machine, if not:
https://go.dev/doc/install
(Make sure you have it correctly installed by issuing:
go version)

Note: if you have a windows machine make sure you get the docker desktop version. 

Once docker is properly installed on your machine issue the following command:

```
docker run --rm -p 5433:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:14
```

This command binds your local port to the docker port, we issue a password to our database by doing "POSTGRES_PASSWORD=VALUE", and we select the postgres database -d flag. 

We then set our ```POSTGRES_URL``` environment variable to:

Windows 
```
set POSTGRES_URL=postgres://postgres:mysecretpassword@localhost:5433/postgres
```


Before starting the server issue this command:
```
go mod tidy
```

Then make sure you are root or admin terminal and start the server:

```
go run main.go
```

We are using sqlc to generate sql queries in go, to install sqlc:
https://docs.sqlc.dev/en/latest/overview/install.html

Then Open http://localhost:8080 with your browser to see the result.

Login route can be accessed on http://localhost:8080/login/

This endpoint can be edited in handlers/handler.go.


To check if the database is working you can issue this command: 

```
psql -h localhost -p 5433 -U postgres
```
The password is the password you set when running the
docker command: 
mysecretpassword

You can also check by going to the docker app and 
looking under inspect 

To check the clients table (make sure display mode is on):

```
postgres=# \x
TABLE CLIENTS;
```
Make sure there is no space at the end of the ";" 
If you currently have no registered clients then 
the table would be empty displaying:

(0 rows)

Currently we only have two tables:

CLIENTS and fuel_history

You can view the schema in db/migrations/001_initial_schema.up.sql 









NOTE: MAIN BRANCH IS BEHIND CHECK SABIN BRANCH
