package api

import "cakes-database-app/pkg/storage"

type API struct {
	db *storage.Interface
}