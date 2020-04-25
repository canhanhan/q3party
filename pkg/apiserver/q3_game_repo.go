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
}

func (r *Q3GameRepository) Close() error {
	log.WithField("caller", "Q3GameRepository.Close").Trace("Started")
	defer log.WithField("caller", "Q3GameRepository.Close").Trace("Exited")

	r.cancel()

	return r.lister.Close()
}

func (r *Q3GameRepository) List() ([]*gamelister.Game, error) {
	log.WithField("caller", "Q3GameRepository.List").Trace("Started")
	defer log.WithField("caller", "Q3GameRepository.List").Trace("Exited")

	return r.lister.List()
}

func (r *Q3GameRepository) Refresh() error {
	log.WithField("caller", "Q3GameRepository.Refresh").Trace("Started")
	defer log.WithField("caller", "Q3GameRepository.Refresh").Trace("Exited")

	err := r.lister.Refresh()
	if err != nil {
		log.WithField("caller", "Q3GameRepository.Refresh").Error(err)
		return err
	}

	return nil
}

func (r *Q3GameRepository) refreshLoop() {
	log.WithField("caller", "Q3GameRepository.refresh").Trace("Started")
	defer log.WithField("caller", "Q3GameRepository.refresh").Trace("Exited")

	for {
		select {
		case <-r.ctx.Done():
			log.WithField("caller", "Q3GameRepository.refresh").Debug(r.ctx.Err())
			return
		default:
		}

		err := r.lister.Refresh()
		if err != nil {
			log.WithField("caller", "Q3GameRepository.refresh").Error(err)
			return
		}

		log.WithField("caller", "Q3GameRepository.refresh").Debug("Sleeping")
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
	}

	go r.refreshLoop()
	return &r, nil
}
