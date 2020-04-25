package gamelister

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/finarfin/q3party/pkg/masterclient"
	"github.com/finarfin/q3party/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type GameLister struct {
	protocol string
	ctx      context.Context
	cancel   context.CancelFunc
	mc       *masterclient.MasterClient
	serverCh chan []string
	games    []*Game
	counter  int32
	mcBusy   bool
	mu       *sync.Mutex
}

func NewGameLister(ctx context.Context, proto string, masters []string) (*GameLister, error) {
	ctx, cancel := context.WithCancel(ctx)
	serverCh := make(chan []string)

	mc, err := masterclient.NewMasterClient(ctx, masters, serverCh)
	if err != nil {
		return nil, err
	}

	g := GameLister{
		protocol: proto,
		ctx:      ctx,
		cancel:   cancel,
		mc:       mc,
		serverCh: serverCh,
		games:    []*Game{},
		counter:  0,
		mcBusy:   true,
		mu:       &sync.Mutex{},
	}

	go g.gameLoop()

	err = g.Refresh()
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func (g *GameLister) Close() error {
	log.WithField("caller", "GameLister.Close").Trace("Started")
	defer log.WithField("caller", "GameLister.Close").Trace("Exited")

	g.cancel()
	return g.mc.Close()
}

func (g *GameLister) IsBusy() bool {
	log.WithField("caller", "GameLister.IsBusy").Trace("Started")
	defer log.WithField("caller", "GameLister.IsBusy").Trace("Exited")

	return g.mcBusy == true || atomic.LoadInt32(&g.counter) != 0
}

func (g *GameLister) List() ([]*Game, error) {
	log.WithField("caller", "GameLister.List").Trace("Started")
	defer log.WithField("caller", "GameLister.List").Trace("Exited")

	return g.games, nil
}

func (g *GameLister) Refresh() error {
	log.WithField("caller", "GameLister.Refresh").Trace("Started")
	defer log.WithField("caller", "GameLister.Refresh").Trace("Exited")

	g.mcBusy = true
	return g.mc.Refresh(&protocol.GetServersRequest{
		Protocol:     g.protocol,
		IncludeEmpty: true,
		IncludeFull:  true,
	})
}

func (g *GameLister) gameLoop() {
	log.WithField("caller", "GameLister.gameLoop").Trace("Started")
	defer log.WithField("caller", "GameLister.gameLoop").Trace("Exited")

	gameCh := make(chan *Game)
	for {
		select {
		case <-g.ctx.Done():
			log.WithField("caller", "GameLister.gameLoop").Debug(g.ctx.Err())
			return
		case res := <-g.serverCh:
			atomic.AddInt32(&g.counter, int32(len(res)))
			g.mcBusy = false
			for _, v := range res {
				go g.getServerInfo(v, gameCh)
			}
		case game := <-gameCh:
			atomic.AddInt32(&g.counter, -1)
			if game == nil {
				continue
			}

			g.games = append(g.games, game)
		}
	}
}

func (g *GameLister) getServerInfo(address string, gameCh chan<- *Game) {
	fields := log.Fields{"caller": "GameLister.getServerInfo", "addr": address}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	for i := 0; i < 3; i++ {
		game := g.tryGetServerInfo(address)
		if game != nil {
			log.WithFields(fields).Trace("Returning game")
			gameCh <- game
			return
		}
	}

	log.WithFields(fields).Trace("Returning empty")
	gameCh <- nil
}

func (g *GameLister) tryGetServerInfo(address string) *Game {
	fields := log.Fields{"caller": "GameLister.tryGetServerInfo", "addr": address}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	cli, err := protocol.NewQ3Client(g.ctx, address)
	if err != nil {
		log.WithFields(fields).Error("Error connecting to server %v", err)
		return nil
	}
	defer cli.Close()

	readCh := make(chan protocol.Q3Message)
	go cli.Listen(readCh)

	err = cli.Send(nil, &protocol.GetChallengeRequest{})
	if err != nil {
		log.WithFields(fields).Errorf("GetChallengeRequest failed: %v", err)
		return nil
	}

	var info *protocol.GetInfoResponse
	var status *protocol.GetStatusResponse
	var challenge string
	var start time.Time
	var latency time.Duration

	for i := 0; i < 5; i++ {
		select {
		case res := <-readCh:
			switch msg := res.Msg.(type) {
			case *protocol.GetInfoResponse:
				log.WithFields(fields).Trace("Received getInfoResponse")
				latency = time.Now().Sub(start)
				info = msg
			case *protocol.GetStatusResponse:
				log.WithFields(fields).Trace("Received getStatusResponse")
				status = msg
			case *protocol.GetChallengeResponse:
				log.WithFields(fields).Trace("Received getChallengeResponse")
				challenge = msg.Challenge
				start = time.Now()
				err = cli.Send(nil, &protocol.GetInfoRequest{Challenge: challenge})
				if err != nil {
					log.WithFields(fields).Errorf("GetInfoRequest failed: %v", err)
					break
				}
				err = cli.Send(nil, &protocol.GetStatusRequest{Challenge: challenge})
				if err != nil {
					log.WithFields(fields).Errorf("GetStatusRequest failed: %v", err)
					break
				}
			default:
				log.WithFields(fields).Errorf("Unknown message from server: %v", msg)
				break
			}
		case <-time.After(1 * time.Second):
			log.WithFields(fields).Debug("Timeout waiting info from server")
			break
		}

		if info != nil && status != nil {
			break
		}
	}

	if info != nil && status != nil {
		return &Game{
			Server: address,
			Info:   info.Data,
			Status: status.Data,
			Ping:   latency.Milliseconds(),
		}
	}

	return nil
}
