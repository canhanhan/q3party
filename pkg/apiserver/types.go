package apiserver

import "github.com/finarfin/q3party/pkg/gamelister"

type GameRepository interface {
	List() ([]*gamelister.Game, error)
	Refresh() error
}

type List struct {
	ID      string   `json:"id"`
	Servers []string `json:"servers"`
}

type ListRepository interface {
	Get(id string) (*List, error)
	Save(list *List) error
}
