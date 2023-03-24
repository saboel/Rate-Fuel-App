package clients
import (
	"github.com/jackc/pgtype"
)

type Client struct {
	ID            pgtype.UUID `json:"id" db:"client_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	optional_address string 
	city string 
	state string 
}




