package apiserver

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/finarfin/q3party/pkg/gamelister"
	log "github.com/sirupsen/logrus"
)

type MockGameRepository struct {
	ctx  context.Context
	file string
	data []*gamelister.Game
}

func (s *MockGameRepository) List() ([]*gamelister.Game, error) {
	return s.data, nil
}

func (s *MockGameRepository) Refresh() error {
	f, err := os.Open(s.file)
	if err != nil {
		return err
	}

	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var data []*gamelister.Game
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	s.data = data
	return nil
}

func (r *MockGameRepository) refreshLoop() {
	log.WithField("caller", "MockGameRepository.refreshLoop").Trace("Started")
	defer log.WithField("caller", "MockGameRepository.refreshLoop").Trace("Exited")

	for {
		select {
		case <-r.ctx.Done():
			log.WithField("caller", "MockGameRepository.refreshLoop").Debug(r.ctx.Err())
			return
		default:
		}

		err := r.Refresh()
		if err != nil {
			log.WithField("caller", "MockGameRepository.refreshLoop").Error(err)
			return
		}

		log.WithField("caller", "Q3GameRepository.refreshLoop").Debug("Sleeping")
		time.Sleep(30 * time.Minute)
	}
}

func NewMockGameRepository(ctx context.Context, file string) (*MockGameRepository, error) {
	r := MockGameRepository{
		ctx:  ctx,
		file: file,
		data: []*gamelister.Game{},
	}

	go r.refreshLoop()

	return &r, nil
}
