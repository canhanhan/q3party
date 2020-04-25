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
	ch              chan<- []string
}

func (c *MasterClient) Close() error {
	log.WithField("caller", "MasterClient.Close").Trace("Started")
	defer log.WithField("caller", "MasterClient.Close").Trace("Exited")

	c.cancel()

	for k := range c.masterConn {
		err := c.masterConn[k].Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MasterClient) Refresh(req *protocol.GetServersRequest) error {
	log.WithField("caller", "MasterClient.Refresh").Trace("Started")
	defer log.WithField("caller", "MasterClient.Refresh").Trace("Exited")

	for _, m := range c.UpstreamServers {
		err := c.masterConn[m].Send(nil, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *MasterClient) readLoop(conn *protocol.Q3Conn) {
	fields := log.Fields{"caller": "MasterClient.readLoop", "addr": conn.Addr.String()}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	readCh := make(chan protocol.Q3Message)
	go conn.Listen(readCh)

	for {
		select {
		case <-c.ctx.Done():
			log.WithFields(fields).Errorf("Context done: %v", c.ctx.Err())
			return
		default:
		}

		log.WithFields(fields).Trace("Waiting for a message")
		msg := <-readCh
		log.WithFields(fields).Trace("Received message")

		switch v := msg.Msg.(type) {
		case *protocol.GetServersResponse:
			log.WithFields(fields).Debug("Received server response")
			c.ch <- v.Servers
		default:
			err := fmt.Errorf("Unknown message type: %s", v)
			log.WithFields(fields).Error(err)
			continue
		}
	}
}

func NewMasterClient(ctx context.Context, masters []string, ch chan<- []string) (*MasterClient, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := MasterClient{
		UpstreamServers: masters,
		ctx:             ctx,
		cancel:          cancel,
		masterConn:      map[string]*protocol.Q3Conn{},
		servers:         []string{},
		ch:              ch,
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
