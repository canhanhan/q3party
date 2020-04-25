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
	Addr         *net.UDPAddr
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	BufferSize   int
	ctx          context.Context
	cancel       context.CancelFunc
	conn         *net.UDPConn
	cancelled    bool
}

func (s *Q3Conn) Close() error {
	fields := log.Fields{"caller": "Q3Conn.Close", "addr": s.Addr.String()}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	s.cancelled = true
	s.cancel()
	return s.conn.Close()
}

func (s *Q3Conn) Send(addr *net.UDPAddr, msg interface{}) error {
	fields := log.Fields{"caller": "Q3Conn.Send", "addr": s.Addr.String()}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	log.WithFields(fields).Tracef("Sending %T", msg)
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
	fields := log.Fields{"caller": "Q3Conn.Listen", "addr": s.Addr.String()}
	log.WithFields(fields).Trace("Started")
	defer log.WithFields(fields).Trace("Exited")

	for {
		select {
		case <-s.ctx.Done():
			log.WithFields(fields).Debug(s.ctx.Err())
			return
		default:
		}

		log.WithFields(fields).Trace("Waiting for data")
		buffer := make([]byte, s.BufferSize)
		s.conn.SetReadDeadline(time.Now().Add(s.ReadTimeout))
		length, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			if s.cancelled {
				return
			}

			if neterr, ok := err.(net.Error); ok {
				switch {
				case neterr.Timeout():
					log.WithFields(fields).Trace(err)
				case neterr.Temporary():
					log.WithFields(fields).Warning(err)
				default:
					log.WithFields(fields).Error(err)
					return
				}
			} else {
				log.WithFields(fields).Error(err)
				return
			}

			continue
		}

		data := buffer[0:length]
		msg, err := Unmarshal(data)
		if err != nil {
			log.WithFields(fields).Error(err)
			continue
		}

		log.WithFields(fields).Tracef("Received %T", msg)
		readCh <- Q3Message{Addr: addr, Msg: msg}
	}
}

func NewQ3Server(ctx context.Context, bind string) (*Q3Conn, error) {
	log.WithField("caller", "Q3Conn.NewQ3server").Trace("Started")
	defer log.WithField("caller", "Q3Conn.NewQ3Server").Trace("Exited")

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
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		BufferSize:   65507,
		ctx:          ctx,
		cancel:       cancel,
		conn:         conn,
	}, nil
}

func NewQ3Client(ctx context.Context, bind string) (*Q3Conn, error) {
	log.WithField("caller", "Q3Conn.NewQ3Client").Trace("Started")
	defer log.WithField("caller", "Q3Conn.NewQ3Client").Trace("Exited")

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
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 1 * time.Second,
		BufferSize:   65507,
		ctx:          ctx,
		cancel:       cancel,
		conn:         conn,
	}, nil
}
