package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/finarfin/q3party/pkg/gamelister"
)

type MockGameRepository struct {
	data []*gamelister.Game
}

func (s *MockGameRepository) List() ([]*gamelister.Game, error) {
	return s.data, nil
}

func NewMockGameRepository(file string) (*MockGameRepository, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var data []map[string]string
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	res := make([]*gamelister.Game, len(data))
	for i, v := range data {
		res[i] = &gamelister.Game{
			Server: v["server"],
			Info:   v,
		}
	}

	return &MockGameRepository{
		data: res,
	}, nil
}
