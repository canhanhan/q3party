package apiserver

import (
	"context"
	"time"

	"github.com/finarfin/q3party/pkg/gamelister"
	log "github.com/sirupsen/logrus"
)

type Q3GameRepository struct {
	ctx    context.Context
	cancel context.CancelFunc
	lister *gamelister.GameLister
	games  []*gamelister.Game
}

func (r *Q3GameRepository) Close() error {
	r.cancel()

	return r.lister.Close()
}

func (r *Q3GameRepository) List() ([]*gamelister.Game, error) {
	return r.games, nil
}

func (r *Q3GameRepository) refresh() {
	for {
		select {
		case <-r.ctx.Done():
			log.Error(r.ctx.Err())
			return
		default:
		}

		g, err := r.lister.List()
		if err != nil {
			log.Error(err)
			return
		}

		r.games = g
		time.Sleep(30 * time.Minute)
	}
}

func NewQ3GameRepository(ctx context.Context, protocol string, masters []string) (*Q3GameRepository, error) {
	gl, err := gamelister.NewGameLister(ctx, protocol, masters)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	r := Q3GameRepository{
		ctx:    ctx,
		cancel: cancel,
		lister: gl,
		games:  []*gamelister.Game{},
	}

	go r.refresh()
	return &r, nil
}
