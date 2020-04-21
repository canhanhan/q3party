package apiserver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type MockListRepository struct {
	data map[string]*List
}

func (r *MockListRepository) Get(id string) (*List, error) {
	list, ok := r.data[id]
	if !ok {
		return nil, errors.New("List was not found")
	}

	return list, nil
}

func (r *MockListRepository) Save(list *List) error {
	r.data[list.ID] = list
	return nil
}

func NewMockListRepository(file string) (*MockListRepository, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var data map[string]*List
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return &MockListRepository{
		data: data,
	}, nil
}
