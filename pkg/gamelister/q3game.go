package gamelister

import (
	"context"
	"sync"
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
}

func NewGameLister(ctx context.Context, protocol string, masters []string) (*GameLister, error) {
	ctx, cancel := context.WithCancel(ctx)

	mc, err := masterclient.NewMasterClient(ctx, masters)
	if err != nil {
		return nil, err
	}

	return &GameLister{
		protocol: protocol,
		ctx:      ctx,
		cancel:   cancel,
		mc:       mc,
	}, nil
}

func (g *GameLister) Close() error {
	g.cancel()
	return g.mc.Close()
}

func (g *GameLister) retrieveMasters() {
	readCh := make(chan *protocol.GetServersResponse)
	for {
		
		select {
		case <-g.ctx.Done():
			log.Error(g.ctx.Err())
			return
		case res := <-readCh
			
		default:
		}
	}

}

func (g *GameLister) List() ([]*Game, error) {
	req := &protocol.GetServersRequest{
		Protocol:     g.protocol,
		IncludeEmpty: true,
		IncludeFull:  true,
	}

	var wg sync.WaitGroup
	wg.Add(len(res.Servers))

	gameCh := make(chan *Game, len(res.Servers))
	for _, v := range res.Servers {
		go g.getServerInfo(v, gameCh, &wg)
	}

	wg.Wait()
	close(gameCh)

	servers := []*Game{}
	for s := range gameCh {
		if s != nil {
			servers = append(servers, s)
		}
	}

	return servers, nil
}

func (g *GameLister) getServerInfo(address string, gameCh chan<- *Game, wg *sync.WaitGroup) {
	defer wg.Done()

	cli, err := protocol.NewQ3Client(g.ctx, address)
	if err != nil {
		log.Error("Error connecting to server %s: %v", address, err)
		wg.Done()
		return
	}
	defer cli.Close()

	readCh := make(chan protocol.Q3Message)
	go cli.Listen(readCh)

	cli.Send(nil, protocol.GetInfoRequest{})
	cli.Send(nil, protocol.GetChallengeRequest{})

	var info *protocol.GetInfoResponse
	var status *protocol.GetStatusResponse
	var challenge string

	for i := 0; i < 3; i++ {
		select {
		case res := <-readCh:
			switch msg := res.Msg.(type) {
			case *protocol.GetInfoResponse:
				info = msg
			case *protocol.GetStatusResponse:
				status = msg
			case *protocol.GetChallengeResponse:
				challenge = msg.Challenge
				err = cli.Send(nil, protocol.GetStatusRequest{Challenge: challenge})
				if err != nil {
					log.Errorf("Unknown message from server %s: %v", address, err)
					return
				}
			default:
				log.Errorf("Unknown message from server %s: %v", address, msg)
				return
			}
		case <-time.After(5 * time.Second):
			log.Errorf("Timeout waiting info from server: %s", address)
			return
		}

		if info != nil && status != nil {
			break
		}
	}

	if info != nil && status != nil {
		gameCh <- &Game{
			Server: address,
			Info:   info.Data,
			Status: status.Data,
		}
	}
}
