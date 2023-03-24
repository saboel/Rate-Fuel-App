package handlers

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgtype"
	"github.com/saboel/Rate-Fuel-App/clients"
	"github.com/saboel/Rate-Fuel-App/db"
)

//go:embed assets
var static embed.FS

type clientGetter interface {
	getClients(ctx context.Context, clientID pgtype.UUID ) (*clients.Client, error)
}

type Handlers struct {
	clientGetter clientGetter
	DB          db.Querier 
}

var indexHTMLTemplate = template.Must(template.ParseFS(static, "assets/Fuel_Quote_Form.html"))
var loginHTMLTemplate = template.Must(template.ParseFS(static, "assets/loginform.gohtml"))

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/fuel_quote":
		h.ServeFuelquote(w, r)
		return
	default:
		h.ServeLoginFrontPage(w, r)
	}
}

func (h *Handlers) getClient(ctx context.Context, uniqueID string) (*clients.Client, error) {
    // Retrieve client ID from the database using the unique identifier
    clientID, err := h.DB.GetClient_id(ctx, uniqueID)
    if err != nil {
        return nil, fmt.Errorf("failed to get client ID: %w", err)
    }

	//clientUUID := string(clientID.Bytes[:])


    // Retrieve client information using the client ID
    client, err := h.clientGetter.getClients(ctx, clientID)
    if err != nil {
        return nil, fmt.Errorf("failed to get client: %w", err)
    }

    return client, nil
}

func (h *Handlers) ServeFuelquote(w http.ResponseWriter, r *http.Request) {

		// Retrieve the client ID from the request context
		uniqueID, ok := r.Context().Value("unique_id").(string)


		
		
		if !ok {
			http.Error(w, "client ID not found", http.StatusInternalServerError)
			return
		}


	// Retrieve client information using the client ID
	client, err := h.getClient(r.Context(), uniqueID)

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


func init () {
	tpl = template.Must(template.ParseGlob("handlers/assets/*"))
}


func (h *Handlers) ServeLoginFrontPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		err := tpl.ExecuteTemplate(w, "loginform.gohtml",nil) 
		if err != nil {
			log.Fatalln(err)
		}
		

	case http.MethodPost: //new user, default method is if they have an account already 
		name:= r.FormValue("fullname") 
		address:= r.FormValue("address1")
		//optional_address:= r.FormValue("address2")
		city:= r.FormValue("city")
		state:= r.FormValue("state")
		clientID, err := h.DB.AddClient(r.Context(), db.AddClientParams{
			ClientName: name,
			AddressClient: address,
			ClientCity: city,
			ClientState: state,
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



