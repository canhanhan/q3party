package apiserver

import (
	"errors"
	"fmt"
	"time"
)

type ListService struct {
	repository ListRepository
}

func (s *ListService) Get(id string) (*List, error) {
	return s.repository.Get(id)
}

func (s *ListService) Create() (*List, error) {
	id := fmt.Sprintf("%X", time.Now().Unix())
	list := &List{
		ID:      id,
		Servers: []string{},
	}

	err := s.repository.Save(list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *ListService) AddToList(id string, server string) error {
	list, err := s.repository.Get(id)
	if err != nil {
		return nil
	}

	for _, v := range list.Servers {
		if v == server {
			return errors.New("Server is already in the list")
		}
	}

	list.Servers = append(list.Servers, server)
	return s.repository.Save(list)
}

func (s *ListService) RemoveFromList(id string, server string) error {
	list, err := s.repository.Get(id)
	if err != nil {
		return nil
	}

	pos := -1
	for i, v := range list.Servers {
		if v == server {
			pos = i
			break
		}
	}

	if pos == -1 {
		return errors.New("Server was not found in list")
	}

	list.Servers = append(list.Servers[0:pos], list.Servers[pos+1:]...)
	return s.repository.Save(list)
}

func NewListService(repository ListRepository) (*ListService, error) {
	return &ListService{
		repository: repository,
	}, nil
}
