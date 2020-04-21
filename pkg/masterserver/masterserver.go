package masterserver

import (
	"context"
	"fmt"
	"net"

	"github.com/finarfin/q3party/pkg/protocol"
	log "github.com/sirupsen/logrus"
)

type ServerReader interface {
	Servers(*protocol.GetServersRequest) (*protocol.GetServersResponse, error)
}

type MasterServer struct {
	reader ServerReader
	conn   *protocol.Q3Conn
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *MasterServer) Close() error {
	s.cancel()
	return s.conn.Close()
}

func (s *MasterServer) listen() error {
	readCh := make(chan protocol.Q3Message)
	go s.conn.Listen(readCh)

	for {
		select {
		case <-s.ctx.Done():
			log.Error(s.ctx.Err())
			return s.ctx.Err()
		default:
		}

		msg := <-readCh

		switch v := msg.Msg.(type) {
		case *protocol.GetServersRequest:
			s.sendGameServersResponse(msg.Addr, v)
		default:
			err := fmt.Errorf("Unknown message type: %s", v)
			log.Error(err)
			continue
		}
	}
}

func (s *MasterServer) sendGameServersResponse(addr *net.UDPAddr, req *protocol.GetServersRequest) error {
	res, err := s.reader.Servers(req)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.conn.Send(addr, res)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func NewMasterServer(ctx context.Context, bind string, reader ServerReader) (*MasterServer, error) {
	ctx, cancel := context.WithCancel(ctx)
	conn, err := protocol.NewQ3Server(ctx, bind)
	if err != nil {
		return nil, err
	}

	s := &MasterServer{
		conn:   conn,
		reader: reader,
		ctx:    ctx,
		cancel: cancel,
	}

	go s.listen()
	return s, nil
}
