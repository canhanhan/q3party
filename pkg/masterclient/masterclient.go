package masterclient

import (
	"context"
	"fmt"

	"github.com/finarfin/q3party/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type MasterClient struct {
	UpstreamServers []string
	servers         []string
	ctx             context.Context
	masterConn      map[string]*protocol.Q3Conn
}

func (c *MasterClient) Close() error {
	for k := range c.masterConn {
		err := c.masterConn[k].Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MasterClient) Servers(req *protocol.GetServersRequest) (*protocol.GetServersResponse, error) {
	return &protocol.GetServersResponse{
		Servers: c.servers,
	}, nil
}

func (c *MasterClient) readLoop(conn *protocol.Q3Conn) error {
	readCh := make(chan protocol.Q3Message)
	for {
		select {
		case <-c.ctx.Done():
			log.Error(c.ctx.Err())
			return c.ctx.Err()
		default:
		}

		msg := <-readCh
		switch v := msg.Msg.(type) {
		case *protocol.GetServersResponse:
			c.updateServerList(v)
		default:
			err := fmt.Errorf("Unknown message type: %s", v)
			log.Error(err)
			continue
		}
	}

	return nil
}

func (c *MasterClient) updateServerList(res *protocol.GetServersResponse) error {
	servers := map[string]bool{}
	for _, s := range c.servers {
		servers[s] = false
	}
	for _, s := range res.Servers {
		servers[s] = false
	}

	sl := make([]string, len(servers))
	i := 0
	for k := range servers {
		sl[i] = k
		i++
	}
	c.servers = sl
	return nil
}

func (c *MasterClient) refreshMasterServers() error {
	req := protocol.GetServersRequest{Protocol: "68"}
	for _, m := range c.UpstreamServers {
		err := c.masterConn[m].Send(nil, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMasterClient(ctx context.Context, masters []string) (*MasterClient, error) {
	c := MasterClient{
		UpstreamServers: masters,
		ctx:             ctx,
		masterConn:      map[string]*protocol.Q3Conn{},
		servers:         []string{},
	}

	for _, m := range masters {
		conn, err := protocol.NewQ3Client(ctx, m)
		if err != nil {
			return nil, err
		}

		c.masterConn[m] = conn
		go c.readLoop(conn)
	}

	return &c, nil
}
