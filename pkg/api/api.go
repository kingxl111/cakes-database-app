package api

import "cakes-database-app/pkg/storage"

// dependency inversion principle
type API struct {
	db storage.Interface
}

func New(db storage.Interface) *API {
	return &API{db: db}
}
