package storage

import "context"

type Interface interface {
	GetCakes(context.Context) ([]Cake, error)
	AddCake(context.Context, Cake) error
}

// TODO: create separate entity
type Cake struct {
	ID 			int
	Description string
	Price		int
	Weight		int
}