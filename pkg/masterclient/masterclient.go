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
	cancel          context.CancelFunc
	masterConn      map[string]*protocol.Q3Conn
}

func (c *MasterClient) Close() error {
	c.cancel()

	for k := range c.masterConn {
		err := c.masterConn[k].Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MasterClient) Servers(req *protocol.GetServersRequest) (*protocol.GetServersResponse, error) {
	err := c.refreshMasterServers()
	if err != nil {
		return nil, err
	}

	return &protocol.GetServersResponse{
		Servers: c.servers,
	}, nil
}

func (c *MasterClient) readLoop(conn *protocol.Q3Conn) error {
	log.Trace("Starting read loop")
	defer log.Trace("Exited read loop")

	readCh := make(chan protocol.Q3Message)
	go conn.Listen(readCh)

	for {
		select {
		case <-c.ctx.Done():
			log.Errorf("Context done: %v", c.ctx.Err())
			return c.ctx.Err()
		default:
		}

		log.Trace("Waiting for a message")
		msg := <-readCh
		log.Trace("Received message")

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
	ctx, cancel := context.WithCancel(ctx)
	c := MasterClient{
		UpstreamServers: masters,
		ctx:             ctx,
		cancel:          cancel,
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

	err := c.refreshMasterServers()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
