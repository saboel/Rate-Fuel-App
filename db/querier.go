// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"

	"github.com/jackc/pgtype"
)

type Querier interface {
	AddClient(ctx context.Context, arg AddClientParams) (Client, error)
	GetClient_id(ctx context.Context, addressClient string) (pgtype.UUID, error)
	GetFuelQuotes(ctx context.Context, clientID pgtype.UUID) ([]FuelHistory, error)
	ListClients(ctx context.Context) ([]Client, error)
}

var _ Querier = (*Queries)(nil)