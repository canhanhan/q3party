package protocol

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Q3Message struct {
	Addr *net.UDPAddr
	Msg  interface{}
}

type Q3Conn struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BufferSize   int
	ctx          context.Context
	cancel       context.CancelFunc
	conn         *net.UDPConn
}

func (s *Q3Conn) Close() error {
	log.Trace("Started Q3Conn.Close")
	defer log.Trace("Exited Q3Conn.Close")

	s.cancel()
	return s.conn.Close()
}

func (s *Q3Conn) Send(addr *net.UDPAddr, msg interface{}) error {
	log.Trace("Started Q3Conn.Send")
	defer log.Trace("Exited Q3Conn.Send")

	data, err := Marshal(msg)
	if err != nil {
		return err
	}

	for _, chunk := range data {
		_, err = s.conn.Write(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Q3Conn) Listen(readCh chan<- Q3Message) {
	log.Trace("Started Q3Conn.Listen")
	defer log.Trace("Exited Q3Conn.Listen")

	for {
		select {
		case <-s.ctx.Done():
			log.Error(s.ctx.Err())
			return
		default:
		}

		log.Trace("Waiting for data")
		buffer := make([]byte, s.BufferSize)
		s.conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		length, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			if neterr, ok := err.(net.Error); ok {
				switch {
				case neterr.Timeout():
					log.Trace(err)
				case neterr.Temporary():
					log.Warning(err)
				default:
					log.Error(err)
					return
				}
			} else {
				log.Error(err)
				return
			}

			continue
		}

		data := buffer[0:length]
		msg, err := Unmarshal(data)
		if err != nil {
			log.Error(err)
			continue
		}

		readCh <- Q3Message{Addr: addr, Msg: msg}
	}
}

func NewQ3Server(ctx context.Context, bind string) (*Q3Conn, error) {
	log.Trace("Started Q3Conn.NewQ3Server")
	defer log.Trace("Exited Q3Conn.NewQ3Server")

	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Q3Conn{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		BufferSize:   65507,
		ctx:          ctx,
		cancel:       cancel,
		conn:         conn,
	}, nil
}

func NewQ3Client(ctx context.Context, bind string) (*Q3Conn, error) {
	log.Trace("Started Q3Conn.NewQ3Client")
	defer log.Trace("Exited Q3Conn.NewQ3Client")

	addr, err := net.ResolveUDPAddr("udp", bind)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	return &Q3Conn{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		BufferSize:   65507,
		ctx:          ctx,
		cancel:       cancel,
		conn:         conn,
	}, nil
}
