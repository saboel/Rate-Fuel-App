package main 
import(
	"net/http"
	"github.com/saboel/Rate-Fuel-App/handlers"
	//"github.com/saboel/Rate-Fuel-App/clients"
	"github.com/saboel/Rate-Fuel-App/db"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4"
	"database/sql"
	"os"
	"log"
	"embed"
)



// start of our server 
func main() {
	// For example: POSTGRES_URL="postgres://postgres:mysecretpassword@localhost:5432/postgres"
	database, err := sql.Open("pgx", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("oops, db connection failed", err)
	}
	err = validateSchema(database)
	if err != nil {
		log.Fatal("oops, db migration failed", err)
	}

	
	h := &handlers.Handlers{
		DB: db.New(database),
	}

	http.Handle("/", h)
	


	log.Println("Serving on localhost:8080")

	err = http.ListenAndServe(":8080", nil) 
	if err != nil {
		log.Fatal("oops, server failed to start", err)
	}



}


//go:embed db/migrations/*.sql
var fs embed.FS

// Migrate migrates the Postgres schema to the current version.

func validateSchema(db *sql.DB) (retErr error) {
	sourceInstance, err := iofs.New(fs, "db/migrations")
	if err != nil {
		return err
	}
	defer func() {
		err := sourceInstance.Close()
		if retErr == nil {
			retErr = err
		}
	}()
	driverInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", driverInstance)
	if err != nil {
		return err
	}
	err = m.Up() // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}