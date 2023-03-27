package clients
import (
	"github.com/jackc/pgtype"
)

type Client struct {
	ID            pgtype.UUID `json:"id" db:"client_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	City string `json:"city"`
	State string `json:"state"`
}




