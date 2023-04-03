package handlers

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	//"github.com/docker/distribution/uuid"
	"github.com/gofrs/uuid"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgtype"
	"github.com/saboel/Rate-Fuel-App/clients"
	"github.com/saboel/Rate-Fuel-App/db"
)

//go:embed assets
var static embed.FS

type clientGetter interface {
	getClients(ctx context.Context, clientID pgtype.UUID) (*clients.Client, error)
}

type Handlers struct {
	clientGetter clientGetter
	DB           db.Querier
}

var indexHTMLTemplate = template.Must(template.ParseFS(static, "assets/Fuel_Quote_Form.html"))
var loginHTMLTemplate = template.Must(template.ParseFS(static, "assets/loginform.gohtml"))
var login_HTMLTemplate = template.Must(template.ParseFS(static, "assets/login.gohtml"))

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/login":
		h.Servelogin(w, r)
	case "/fuel_quote":
		h.ServeFuelquote(w, r)
		return
	default:
		h.ServeLoginFrontPage(w, r)
	}
}



type MyContext struct {
	ClientID string
}

// Define a new context key of type *MyContext.
type contextKey string

var myContextKey contextKey = "my-context"

func (h *Handlers) Servelogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := login_HTMLTemplate.Execute(w, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to render html template: %v", err), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		//check submitted login creds
		name := r.FormValue("fullname")
		address := r.FormValue("address")
		client, err := h.DB.GetClientByNameAndAddress(r.Context(), db.GetClientByNameAndAddressParams{
			ClientName:    name,
			AddressClient: address,
		})

		if err != nil {
			log.Printf("Failed to retrieve client information: %v", err)

			// Handle the error (e.g. show an error message to the user)
			http.Error(w, "Failed to retrieve client information", http.StatusInternalServerError)
			return
		}

		if client.ClientID.Status != pgtype.Present {
			// If the client does not exist in the database, display an error message
			http.Error(w, "invalid login credentials", http.StatusBadRequest)
			return
		}

		//client exists redirect to the fuel quote form
		// Get the byte slice representation of the UUID
		idBytes := client.ClientID.Bytes
		// Convert the byte slice to a uuid.UUID value
		id, err := uuid.FromBytes(idBytes[:])
		if err != nil {
			// Handle the error
		}



		ctx := context.WithValue(r.Context(), myContextKey, &MyContext{ClientID: id.String()})
		log.Printf("client ID: %v", id.String())
		http.Redirect(w, r.WithContext(ctx), fmt.Sprintf("/fuel_quote?client_id=%s", id.String()), http.StatusFound)

	}
}

func (h *Handlers) getClient(ctx context.Context, uniqueID pgtype.UUID) (*clients.Client, error) {
	// Retrieve client ID from the database using the unique identifier
	myCtx, ok := ctx.Value(myContextKey).(*MyContext)

	if !ok {
        return nil, fmt.Errorf("invalid context type")
    }
	log.Printf("myCtx: %+v", myCtx) // Check if myCtx is correct
	

	//clientUUID := string(clientID.Bytes[:])

	// Retrieve client information using the client ID
	clientdb, err := h.DB.GetClientInfo(ctx, uniqueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	client := &clients.Client{
		ID: clientdb.ClientID,
		Name: clientdb.ClientName,
		Address: clientdb.AddressClient,
		City: clientdb.ClientCity,
		State: clientdb.ClientState,
	}

	return client, nil
}

func (h *Handlers) ServeFuelquote(w http.ResponseWriter, r *http.Request) {

	// Retrieve the client ID from the request context
	//clientID := r.Context().Value("client_id")
	//fmt.Printf("clientID type: %T\n", clientID)

	
	uniqueID := r.URL.Query().Get("client_id")
	log.Printf("uniqueID: %v", uniqueID)



	if uniqueID == "" {
		http.Error(w, "client ID not found", http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(r.Context(), myContextKey, &MyContext{ClientID: uniqueID})

	myCtx, ok := ctx.Value(myContextKey).(*MyContext)
    if !ok {
        http.Error(w, "invalid context type", http.StatusInternalServerError)
        return
    }
    log.Printf("myCtx: %+v", myCtx)
	var clientID pgtype.UUID
	err := clientID.Set(uniqueID)

	if err != nil {
        http.Error(w, fmt.Sprintf("invalid client ID: %v", err), http.StatusBadRequest)
        return
    }


	// Retrieve client information using the client ID
	client, err := h.getClient(ctx, clientID)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get client: %v", err), http.StatusInternalServerError)
		return
	}
	if client == nil {
		// Client is not logged in, redirect to login page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}


	quotes, err := h.DB.GetFuelQuotes(r.Context(), client.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get fuel quotes: %v", err), http.StatusInternalServerError)
		return
	}

	// Render fuel quotes template with retrieved quotes
	err = indexHTMLTemplate.Execute(w, quotes)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to render template: %v", err), http.StatusInternalServerError)
		return
	}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("handlers/assets/*"))
}

func (h *Handlers) ServeLoginFrontPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		err := tpl.ExecuteTemplate(w, "loginform.gohtml", nil)
		if err != nil {
			log.Fatalln(err)
		}

	case http.MethodPost: //new user, default method is if they have an account already
		name := r.FormValue("fullname")
		address := r.FormValue("address1")
		//optional_address:= r.FormValue("address2")
		city := r.FormValue("city")
		state := r.FormValue("state")
		clientID, err := h.DB.AddClient(r.Context(), db.AddClientParams{
			ClientName:    name,
			AddressClient: address,
			ClientCity:    city,
			ClientState:   state,
		})

		if err != nil {
			var e *pgconn.PgError
			if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
				http.Error(w, "that client is already in database", http.StatusBadRequest)
				return
			}
			http.Error(w, fmt.Sprintf("failed to add client : %v", err), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "client_id", clientID)
		http.Redirect(w, r.WithContext(ctx), "/fuel_quote", http.StatusFound)

	}

}
